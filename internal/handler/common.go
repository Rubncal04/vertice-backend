package handler

// ErrorResponse representa una respuesta de error estándar
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
