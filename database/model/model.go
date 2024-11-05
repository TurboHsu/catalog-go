package model

import (
	"time"
)

var ALL = []interface{}{
	Cats{},
	Reactions{},
}

type Cats struct {
	UUID      string      `gorm:"primaryKey" json:"uuid"`
	Caption   string      `json:"caption"`
	Image     string      `json:"image"`
	Thumbnail string      `json:"thumbnail"`
	CreatedAt time.Time   `json:"created_at"`
	Reactions []Reactions `gorm:"foreignKey:CatUUID;references:UUID" json:"reactions"`
}

type Reactions struct {
	ID      uint   `gorm:"primarykey"`
	Emoji   string `gorm:"type:varchar(255)"`
	Client  string `gorm:"type:varchar(255)"`
	CatUUID string `gorm:"type:varchar(255)"`
}
