package cat

import (
	"catalog-go/database"
	"catalog-go/database/query"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getAllHandler(c *gin.Context) {
	q := query.Use(database.DB)

	// Get all cat
	cats, err := q.WithContext(c.Request.Context()).Cats.Find()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": cats,
	})
}

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

	q := query.Use(database.DB)
	cats, err := q.WithContext(c.Request.Context()).Cats.
		Order(q.Cats.CreatedAt.Desc()).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"data": cats,
	})
}
