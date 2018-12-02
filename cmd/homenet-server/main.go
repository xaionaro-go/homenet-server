package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/methods"
	"github.com/xaionaro-go/homenet-server/middleware"
	"github.com/xaionaro-go/homenet-server/storage"
)

func fatalIfError(err error) {
	if err != nil {
		logrus.Fatalf("%v", err)
	}
}

func main() {
	storage.Init(`/home/homenetsrv/storage`)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.PUT("/:net", methods.RegisterNet)
	netAuthed := router.Group("/:net")
	netAuthed.Use(middleware.GetNetwork())
	netAuthed.GET("", methods.GetNet)
	netAuthed.GET("/peers", methods.GetPeers)
	netAuthed.PUT("/peers/:id", methods.RegisterPeer)
	netAuthed.DELETE("/peers/:id", methods.UnregisterPeer)
	netAuthed.PUT("/negotiationMessage/:peer_id_to/:peer_id_from", methods.SetNegotiationMessage)
	netAuthed.GET("/negotiationMessage/:peer_id_to", methods.GetNegotiationMessages)
	netAuthed.GET("/negotiationMessage/:peer_id_to/:peer_id_from", methods.GetNegotiationMessage)
	fatalIfError(router.Run())
}
