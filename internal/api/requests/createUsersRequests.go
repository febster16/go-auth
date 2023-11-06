package requests

type UserPayload struct {
	Email    string `json:"email" example:"john.doe@mail.com"`
	Password string `json:"password" example:"password"`
}
