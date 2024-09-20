package model

import (
	"time"

	"gorm.io/gorm"
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
	gorm.Model
	Emoji   string `gorm:"type:varchar(255)"`
	Clients string `gorm:"type:text" json:"clients"`
	CatUUID string `gorm:"type:varchar(255)"`
}
