package dto

type PublicUserPayload struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CustID   string `json:"cust_id"`
}

type PublicUserLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
