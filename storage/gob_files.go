package storage

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/xaionaro-go/atomicmap"

	"github.com/xaionaro-go/homenet-server/iface"
)

var savingQueue = make(chan iface.Object)
var storageDir string

func init() {
	gob.Register(atomicmap.New())
}

func Init(newStorageDir string) {
	if storageDir != "" {
		logrus.Panicf(`The storage is already initialized`)
	}
	if newStorageDir == "" {
		logrus.Errorf(`newStorageDir = ""`)
		return
	}
	storageDir = newStorageDir

	for typeName, storage := range storages.ToSTDMap() {
		restore(storage.(atomicmap.Map), typeName.(string))
	}

	go savingQueueHandler()
	go periodicallySaveEverything()
}

func GetSavingQueue() chan iface.Object {
	return savingQueue
}

func savingQueueHandler() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		<-ticker.C
		hasAnythingToSave := false
		{
			for {
				select {
				case <-savingQueue:
					fmt.Println("checkQueue: found")
					hasAnythingToSave = true
					continue
				default:
				}
				break
			}
		}

		if !hasAnythingToSave {
			continue
		}

		saveEverything() // TODO: remove this from here
		/*objs := map[string]atomicmap.Map{}
		for obj, ok := <-savingQueue; ok {
			objTypeName := fmt.Sprintf(`%T`, obj)
			if objs[objTypeName] == nil {
				objs[objTypeName] = atomicmap.New()
			}
			atomicmap.Set(obj.IGetID(), obj)
		}

		for objTypeName, objMap := range objs {
			f, err := newHistoryFile(objTypeName)
			if err != nil {
				logrus.
			}
			err = gob.Marshal(objMap.ToSTDMap(), )
			f.Write()
			for id, obj := range objMap.ToSTDMap() {
				f.
			}
		}*/
	}
}

var storages = atomicmap.New()

func Get(sample iface.Object) atomicmap.Map {
	sampleTypeName := fmt.Sprintf("%T", sample)
	if len(sampleTypeName) > 1 && sampleTypeName[0] == '*' {
		sampleTypeName = sampleTypeName[1:]
	}
	storage, err := storages.Get(sampleTypeName)
	switch err {
	case nil:
	case atomicmap.NotFound:
		gob.Register(sample)
		storage = atomicmap.New()
		if err != nil {
			logrus.Errorf("Cannot restore the data for %T: %v", sample, err)
		}
		storages.Set(sampleTypeName, storage)
	default:
		logrus.Errorf("Unexpected error from storages.Get(): %v", err)
		storage = atomicmap.New()
		storages.Set(sampleTypeName, storage)
	}
	return storage.(atomicmap.Map)
}

func restore(storage atomicmap.Map, typeName string) error {
	storagePath := filepath.Join(storageDir, `current`, typeName+`.gob`)
	decodeFile, err := os.Open(storagePath)
	if err != nil {
		if err.(*os.PathError).Err == syscall.ENOENT {
			return nil
		}
		return err
	}
	defer decodeFile.Close()
	decoder := gob.NewDecoder(decodeFile)

	stdMap := map[atomicmap.Key]interface{}{}
	err = decoder.Decode(&stdMap)
	if err != nil {
		return err
	}

	storage.FromSTDMap(stdMap)
	return nil
}

func saveEverything() {
	now := time.Now()
	ts := now.UnixNano()
	storePath := filepath.Join(storageDir, fmt.Sprintf(`full-%d`, ts))
	err := os.MkdirAll(storePath, os.FileMode(0700))
	if err != nil {
		logrus.Errorf("Cannot create a directory to save the storage data: %v", err)
		return
	}

	f, err := os.Create(filepath.Join(storePath, `README.md`))
	if err == nil {
		f.Write([]byte("```\n" + `curl -s https://golang.org/src/encoding/gob/dump.go?m=text > /tmp/gobdump.$$.go && go run /tmp/gobdump.$$.go models.network.gob; rm -f /tmp/gobdump.$$.go` + "\n```\n"))
		f.Close()
	}

	for typeName, objMap := range storages.ToSTDMap() {
		fPath := filepath.Join(storePath, typeName.(string)+`.gob`)
		f, err := os.Create(fPath)
		if err != nil {
			logrus.Errorf("Cannot open file to save the storage data of %s: %v", typeName, err)
			continue
		}
		defer f.Close()

		gobEncoder := gob.NewEncoder(f)
		err = gobEncoder.Encode(objMap.(atomicmap.Map).ToSTDMap())
		if err != nil {
			logrus.Errorf("Cannot save the storage data of %s: %v", typeName, err)
			continue
		}
	}

	symlinkPath := filepath.Join(storageDir, `current`)

	oldDir, err := os.Readlink(symlinkPath)
	if err != nil && err.(*os.PathError).Err != syscall.ENOENT {
		logrus.Errorf(`Unable to read a symlink "%s": %v`, symlinkPath, err)
	}

	err = os.Remove(symlinkPath)
	if err != nil && err.(*os.PathError).Err != syscall.ENOENT {
		logrus.Errorf(`Unable to unlink the symlink to the actual data: %v (%v)`, err)
	}

	err = os.Symlink(storePath, symlinkPath)
	if err != nil {
		logrus.Errorf(`Unable to create the symlink to the actual data: %v`, err)
		return
	}

	if oldDir != "" {
		err = os.RemoveAll(oldDir)
		if err != nil {
			logrus.Errorf(`Unable to delete the old data directory "%v": %v`, oldDir, err)
			return
		}
	}
}

func periodicallySaveEverything() {
	ticker := time.NewTicker(24 * time.Hour)
	for {
		<-ticker.C
		saveEverything()
	}
}
