package model

var ALL = []interface{}{
	Cats{},
}

type Cats struct {
	UUID string `gorm:"primaryKey"`
	Caption string
	Image string
}