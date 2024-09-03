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

	"github.com/google/uuid"
)

func Place(file File, caption string, ctx context.Context) (err error) {
	uuid := uuid.NewString()
	// Save file
	filename := config.CONFIG.Store.StorePath + "/" + uuid + "." + file.Type
	if _, err := os.Stat(config.CONFIG.Store.StorePath); os.IsNotExist(err) {
		os.Mkdir(config.CONFIG.Store.StorePath, 0755)
	}
	err = os.WriteFile(filename, file.Buffer, 0644)
	if err != nil {
		return err
	}

	cat := model.Cats{
		UUID:   uuid,
		Caption: caption,
		Image: uuid + "." + file.Type,
		CreatedAt: time.Now(),
	}
	c := query.Use(database.DB).Cats
	err = c.WithContext(ctx).Create(&cat)

	go func(filename string)  {
		store.PutFileHook(filename)
	}(filename)
	
	return
}

type File struct {
	Buffer []byte
	Type  string // Extension
}