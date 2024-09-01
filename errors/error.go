package errors

type ErrResponse struct {
	Message   string `json:"message"`
	Status    int
}