package model

type Users struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Priv     uint64 `json:"priv"`
}
