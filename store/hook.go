package store

import (
	"catalog-go/config"
	"catalog-go/store/qiniu"
)

func PutFileHook(path string) error {
	if config.CONFIG.Store.Qiniu.Enable {
		err := qiniu.PutFile(path)
		if err != nil {
			return err
		}
	}

	return nil
}
