package requests

type ChangePasswordPayload struct {
	Email       string `json:"email" example:"john.doe@mail.com"`
	OldPassword string `json:"old_password" example:"password1"`
	NewPassword string `json:"new_password" example:"password2"`
}
