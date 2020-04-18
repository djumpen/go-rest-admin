package apperrors

type Validation struct {
	namespace string
	message   string
}

func NewValidation(ns, msg string) *Validation {
	return &Validation{
		namespace: ns,
		message:   msg,
	}
}

func (e *Validation) Error() string {
	return e.message
}

func (e *Validation) Namespace() string {
	return e.namespace
}
