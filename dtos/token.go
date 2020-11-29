package dtos

// Token
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpireAt    int64  `json:"expire_at"`
}

// UserInToken
type UserInToken struct {
	UserID int    `json:"user_id"`
	Phone  string `json:"phone"`
	Utype  int    `json:"utype"`
}

// SmsToken
type SmsToken struct {
	Smsid    string `json:"smsid"`
	ExpireAt int64  `json:"expire_at"`
}
