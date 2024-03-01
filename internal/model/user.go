package model

type User struct {
	ID       string `json:"id,omitempty"`   //uuid пользователя
	Login    string `json:"login"`          //login
	Password string `json:"password"`       //login
	Hash     string `json:"hash,omitempty"` //hash for password
}
