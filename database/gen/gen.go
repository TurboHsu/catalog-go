package main

import (
	"catalog-go/database/model"

	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../query",
	})

	for _, model := range model.ALL {
		g.ApplyBasic(model)
	}

	g.Execute()
}
