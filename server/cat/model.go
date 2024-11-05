package cat

import (
	"catalog-go/database/model"
	"time"
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
	r.Reactions = []Reactions{}

	if len(cat.Reactions) > 0 {
		reactionMap := make(map[string]*Reactions)
		for _, reaction := range cat.Reactions {

			if existingReaction, exists := reactionMap[reaction.Emoji]; exists {
				existingReaction.Count++
				if reaction.Client == fingerprint {
					existingReaction.IsReacted = true
				}
			} else {
				reactionMap[reaction.Emoji] = &Reactions{
					Emoji:     reaction.Emoji,
					Count:     1,
					IsReacted: reaction.Client == fingerprint,
				}
			}
		}

		for _, reaction := range reactionMap {
			r.Reactions = append(r.Reactions, *reaction)
		}
	}

	return nil
}
