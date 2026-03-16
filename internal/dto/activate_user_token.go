package dto

type ActivateUserToken struct {
	Token string `json:"token" binding:"required"`
}
