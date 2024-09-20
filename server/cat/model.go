package cat

import (
	"catalog-go/database/model"
	"time"

	"github.com/goccy/go-json"
)

type CatResponse struct {
	UUID      string      `json:"uuid"`
	Caption   string      `json:"caption"`
	Image     string      `json:"image"`
	Thumbnail string      `json:"thumbnail"`
	CreatedAt time.Time   `json:"created_at"`
	Reactions []Reactions `json:"reactions"`
}

type Reactions struct {
	Emoji     string `json:"emoji"`
	Count     int    `json:"count"`
	IsReacted bool   `json:"is_reacted"`
}

func (r *CatResponse) FromCats(cat *model.Cats, fingerprint string) error {
	r.UUID = cat.UUID
	r.Caption = cat.Caption
	r.Image = cat.Image
	r.Thumbnail = cat.Thumbnail
	r.CreatedAt = cat.CreatedAt

	for _, reaction := range cat.Reactions {
		var clients []string
		if len(reaction.Clients) > 2 { // Bigger than []
			err := json.Unmarshal([]byte(reaction.Clients), &clients)
			if err != nil {
				return err
			}
			
			r.Reactions = append(r.Reactions, Reactions{
				Emoji:     reaction.Emoji,
				Count:     len(clients),
				IsReacted: contains(clients, fingerprint),
			})
		}
	}
	return nil
}
