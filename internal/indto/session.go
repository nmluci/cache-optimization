package indto

type UserSession struct {
	UserID   uint64 `json:"user_id"`
	Email    uint64 `json:"email"`
	Fullname uint64 `json:"fullname"`
	LoginAt  uint64 `json:"login_at"`
	Priv     uint64 `json:"priv"`
}
