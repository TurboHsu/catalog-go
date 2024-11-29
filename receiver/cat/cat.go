package cat

import (
	"catalog-go/config"
	"catalog-go/database"
	"catalog-go/database/model"
	"catalog-go/database/query"
	"catalog-go/store"
	"context"
	"os"
	"time"

	u "github.com/google/uuid"
)

func Place(image File, thumbnail File, caption string, ctx context.Context) (uuid string, err error) {
	uuid = u.NewString()
	// Save file
	filename := config.CONFIG.Store.StorePath + "/" + uuid + "." + image.Type
	if _, err := os.Stat(config.CONFIG.Store.StorePath); os.IsNotExist(err) {
		os.Mkdir(config.CONFIG.Store.StorePath, 0755)
	}
	err = os.WriteFile(filename, image.Buffer, 0644)
	if err != nil {
		return "", err
	}

	thumbnailFilename := config.CONFIG.Store.StorePath + "/" + uuid + "." + "thumbnail" + "." + thumbnail.Type
	err = os.WriteFile(thumbnailFilename, thumbnail.Buffer, 0644)
	if err != nil {
		return "", err
	}

	cat := model.Cats{
		UUID:      uuid,
		Caption:   caption,
		Image:     uuid + "." + image.Type,
		Thumbnail: uuid + "." + "thumbnail" + "." + thumbnail.Type,
		CreatedAt: time.Now(),
	}
	c := query.Use(database.DB).Cats
	err = c.WithContext(ctx).Create(&cat)

	go func(filename string, thumbnailFilename string) {
		store.PutFileHook(filename)
		store.PutFileHook(thumbnailFilename)
	}(filename, thumbnailFilename)

	return
}

func Remove(uuid string, ctx context.Context) (err error) {
	q := query.Use(database.DB)
	cat, err := q.WithContext(ctx).Cats.Preload(q.Cats.Reactions).Where(q.Cats.UUID.Eq(uuid)).First()
	if err != nil {
		return err
	}

	for _, r := range cat.Reactions {
		reaction, err := q.WithContext(ctx).Reactions.Where(q.Reactions.ID.Eq(r.ID)).First()
		if err != nil {
			return err
		}
		_, err = q.WithContext(ctx).Reactions.Delete(reaction)
		if err != nil {
			return err
		}
	}

	// Delete file
	filename := config.CONFIG.Store.StorePath + "/" + cat.Image
	thumbnailFilename := config.CONFIG.Store.StorePath + "/" + cat.Thumbnail
	os.Remove(filename)
	os.Remove(thumbnailFilename)
	go func(filename string, thumbnailFilename string) {
		store.RemoveFileHook(filename)
		store.RemoveFileHook(thumbnailFilename)
	}(filename, thumbnailFilename)

	_, err = q.WithContext(ctx).Cats.Delete(cat)
	return err
}

type File struct {
	Buffer []byte
	Type   string // Extension
}
