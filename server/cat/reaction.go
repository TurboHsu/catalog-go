package cat

import (
	"catalog-go/config"
	"catalog-go/database"
	"catalog-go/database/model"
	"catalog-go/database/query"
	"log"

	"github.com/gin-gonic/gin"
)

func getValidReactions(c *gin.Context) {
	c.JSON(200, gin.H{"data": config.CONFIG.Database.AllowedReactions})
}

func addReactionHandler(c *gin.Context) {
	fingerprint, ok := c.GetQuery("fingerprint")
	if !ok || fingerprint == "" {
		c.JSON(400, gin.H{"error": "fingerprint not provided"})
		return
	}
	if len(fingerprint) > 32 {
		c.JSON(400, gin.H{"error": "fingerprint too long"})
		return
	}
	catUUID, ok := c.GetQuery("cat")
	if !ok {
		c.JSON(400, gin.H{"error": "cat not provided"})
		return
	}
	reaction, ok := c.GetQuery("reaction")
	if !ok {
		c.JSON(400, gin.H{"error": "reaction not provided"})
		return
	}
	if !contains(config.CONFIG.Database.AllowedReactions, reaction) {
		c.JSON(400, gin.H{"error": "reaction not allowed"})
		return
	}

	q := query.Use(database.DB)
	result, err := q.WithContext(c.Request.Context()).Reactions.Where(q.Reactions.CatUUID.Eq(catUUID),
		q.Reactions.Emoji.Eq(reaction),
		q.Reactions.Client.Eq(fingerprint)).Count()

	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		log.Printf("[E] %v", err)
		return
	}

	if result > 0 {
		c.JSON(400, gin.H{"error": "already reacted"})
		return
	}

	err = q.WithContext(c.Request.Context()).Reactions.Create(&model.Reactions{
		CatUUID: catUUID,
		Emoji:   reaction,
		Client:  fingerprint,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		log.Printf("[E] %v", err)
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

func removeReactionHandler(c *gin.Context) {
	fingerprint := c.Query("fingerprint")
	catUUID := c.Query("cat")
	reaction := c.Query("reaction")
	if !contains(config.CONFIG.Database.AllowedReactions, reaction) {
		c.JSON(400, gin.H{"error": "reaction not allowed"})
		return
	}
	q := query.Use(database.DB)

	_, err := q.WithContext(c.Request.Context()).Reactions.Where(q.Reactions.CatUUID.Eq(catUUID),
		q.Reactions.Emoji.Eq(reaction),
		q.Reactions.Client.Eq(fingerprint)).Delete()

	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		log.Printf("[E] %v", err)
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
