package dtos

// CLaw 服务条款
type CLaw struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Agreement string `gorm:"type:text" json:"agreement"`
}
