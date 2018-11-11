package errors

type BadRequest struct {
	anyError
}

func (e BadRequest) IsBadRequest() {
}
