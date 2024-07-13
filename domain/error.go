package domain


type ICustomError interface {
	Error() string
}


type CustomError struct {
	StatusCode uint
	ErrMsg     string
}

func (e *CustomError) Error() string {
	return e.ErrMsg
}

func NewErr(StatusCode uint, ErrMsg string) error {
	return &CustomError{
		StatusCode: StatusCode,
		ErrMsg:     ErrMsg,
	}
}
