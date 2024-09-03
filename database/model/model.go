package model

import "time"

var ALL = []interface{}{
	Cats{},
}

type Cats struct {
	UUID      string `gorm:"primaryKey"`
	Caption   string
	Image     string
	CreatedAt time.Time
}
