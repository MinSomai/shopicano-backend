package errors

import "encoding/json"

type preparedError struct {
	Err    ValidationError `json:"errors,omitempty"`
	Status int             `json:"-"`
	Code   string          `json:"code,omitempty"`
	Title  string          `json:"title,omitempty"`
}

func NewPreparedError() *preparedError {
	return &preparedError{
		Err: ValidationError{},
	}
}

func IsPreparedError(v interface{}) bool {
	_, ok := v.(*preparedError)
	return ok
}

func (pe *preparedError) Error() string {
	b, _ := json.Marshal(pe)
	return string(b)
}

func (pe preparedError) AddValidationError(key, value string) {
	pe.Err[key] = append(pe.Err[key], value)
}
