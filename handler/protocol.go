package handler

type BaseInfo struct {
	Account     string `json:"account"`
	IsActive    bool `json:"isactive"`
	IsSuperUser bool `json:"issuperuser"`
	UserType    int `json:"usertype"`
}