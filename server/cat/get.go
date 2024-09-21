package cat

import (
	"catalog-go/database"
	"catalog-go/database/query"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getHandler(c *gin.Context) {
	page := 1
	pageSize := 10
	var err error

	if p, ok := c.GetQuery("page"); ok && p != "" {
		page, err = strconv.Atoi(p)
		if page < 1 {
			err = errors.New("invalid offset")
		}
		if err != nil {
			c.JSON(400, gin.H{
				"error": "bad request",
			})
			return
		}
	}
	if s, ok := c.GetQuery("page-size"); ok && s != "" {
		pageSize, err = strconv.Atoi(s)
		if pageSize >= 25 && pageSize < 1 {
			err = errors.New("invalid pageSize")
		}
		if err != nil {
			c.JSON(400, gin.H{
				"error": "bad request",
			})
			return
		}
	}
	fingerprint, ok := c.GetQuery("fingerprint")
	if !ok || fingerprint == "" {
		c.JSON(400, gin.H{
			"error": "bad fingerprint",
		})
		return
	}

	q := query.Use(database.DB)
	cats, err := q.WithContext(c.Request.Context()).Cats.
		Preload(q.Cats.Reactions).
		Order(q.Cats.CreatedAt.Desc()).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find()

	data := make([]CatResponse, len(cats))
	for i, cat := range cats {
		data[i].FromCats(cat, fingerprint)
	}

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"data": data,
	})
}

func getByIdHandler(c *gin.Context) {
	uuid, ok := c.GetQuery("uuid")
	if !ok {
		c.JSON(400, gin.H{
			"error": "bad request",
		})
		return
	}
	fingerprint, ok := c.GetQuery("fingerprint")
	if !ok || fingerprint == "" {
		c.JSON(400, gin.H{
			"error": "bad fingerprint",
		})
		return
	}

	q := query.Use(database.DB)
	cat, err := q.WithContext(c.Request.Context()).Cats.
		Preload(q.Cats.Reactions).
		Where(q.Cats.UUID.Eq(uuid)).
		First()

	if err != nil {
		c.JSON(404, gin.H{
			"error": "not found",
		})
		return
	}

	var data CatResponse
	err = data.FromCats(cat, fingerprint)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": data,
	})
}
