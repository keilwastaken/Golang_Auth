package Models

type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"-"`
}
