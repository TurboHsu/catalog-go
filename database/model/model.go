package model

import "time"

var ALL = []interface{}{
    Cats{},
}

type Cats struct {
    UUID      string    `gorm:"primaryKey" json:"uuid"`
    Caption   string    `json:"caption"`
    Image     string    `json:"image"`
    Thumbnail string    `json:"thumbnail"`
    CreatedAt time.Time `json:"created_at"`
}