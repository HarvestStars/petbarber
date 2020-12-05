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

type UserSignupRep struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}

type PetHouseSigninRep struct {
	User     User       `json:"user"`
	PetHouse TuPethouse `json:"user_info"`
	Token    Token      `json:"token"`
}

type GroomerSigninRep struct {
	User    User      `json:"user"`
	Groomer TuGroomer `json:"user_info"`
	Token   Token     `json:"token"`
}

type UserPetHouseProfileRep struct {
	User  TuPethouse `json:"user_info"`
	Token Token      `json:"token"`
}

type UserGroomerProfileRep struct {
	User  TuGroomer `json:"user_info"`
	Token Token     `json:"token"`
}
