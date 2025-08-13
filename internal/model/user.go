package model

type RequestCreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseLoginUser struct {
	AccessToken string `json:"access_token"`
	ID          string `json:"id"`
	Username    string `json:"username"`
}
