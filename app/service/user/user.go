package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(20); not null"`
	Password string `json:"password" gorm:"size: 255; not null;"`
}

func (u User) GetUid() string {
	return u.Name
}
