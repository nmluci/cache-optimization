package dto

type PublicUserPayload struct {
	Fullname   string  `json:"fullname"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Priv       uint64  `json:"access_lv"`
	SessionKey string  `json:"-"`
	MasterKey  *string `json:"-"`
}

type PublicUserLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
