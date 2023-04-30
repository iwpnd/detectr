package errors

// ErrRequestError ...
type ErrRequestError struct {
	Status int    `json:"status"`
	Source string `json:"source,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"details,omitempty"`
}

// ErrorResponse ...
type ErrorResponse struct {
	Errors []*ErrRequestError `json:"errors"`
}

// NewRequestError ...
func NewRequestError(e ...*ErrRequestError) ErrorResponse {
	error := ErrorResponse{}

	error.Errors = append(error.Errors, e...)

	return error
}
