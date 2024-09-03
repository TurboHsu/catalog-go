package cat

import (
	"catalog-go/database"
	"catalog-go/database/query"

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
