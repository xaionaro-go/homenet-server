package errors

type anyError struct {
	message string
}

func (e anyError) Error() string {
	return e.message
}
