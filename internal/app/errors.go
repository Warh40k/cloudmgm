package app

// NotFoundError not found
type NotFoundError struct {
	Message string `json:"message"`
	Err     error
}

func (e NotFoundError) Error() string {
	return e.Message
}

func (e NotFoundError) Unwrap() error {
	return e.Err
}

// BadRequestError bad request
type BadRequestError struct {
	Message string `json:"message"`
	Err     error
}

func (e BadRequestError) Error() string {
	return e.Message
}

func (e BadRequestError) Unwrap() error {
	return e.Err
}

// InternalError internal error
type InternalError struct {
	Message string `json:"message"`
	Err     error
}

func (e InternalError) Error() string {
	return e.Message
}

func (e InternalError) Unwrap() error {
	return e.Err
}
