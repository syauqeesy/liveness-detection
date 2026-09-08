package common

type ApplicationException struct {
	HttpStatusCode int
	Message        string
}

func (e ApplicationException) Error() string {
	return e.Message
}

func CreateException(httpStatusCode int, message string) *ApplicationException {
	return &ApplicationException{
		HttpStatusCode: httpStatusCode,
		Message:        message,
	}
}
