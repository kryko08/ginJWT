package utils

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}

func (r *RequestError) Status() int {
	return r.StatusCode
}
