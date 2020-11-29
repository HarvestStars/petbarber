package dtos

type User struct {
	UserID   uint   `json:"user_id"`
	Phone    string `json:"phone"`
	UserType int    `json:"user_type"`
}

type UserSigninReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
	Smsid string `json:"smsid"`
}

type UserSigninRep struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}

type UserPetHouseProfileRep struct {
	User  TuPethouse `json:"user"`
	Token Token      `json:"token"`
}

type UserGroomerProfileRep struct {
	User  TuGroomer `json:"user"`
	Token Token     `json:"token"`
}
