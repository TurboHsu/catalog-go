package cat

import (
	"catalog-go/config"
	"catalog-go/database"
	"catalog-go/database/model"
	"catalog-go/database/query"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func addReactionHandler(c *gin.Context) {
	fingerprint, ok := c.GetQuery("fingerprint")
	if !ok || fingerprint == "" {
		c.JSON(400, gin.H{"error": "fingerprint not provided"})
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
	cat, err := q.WithContext(c.Request.Context()).Cats.Preload(q.Cats.Reactions).Where(q.Cats.UUID.Eq(catUUID)).First()
	if err != nil {
		c.JSON(400, gin.H{"error": "cat not found"})
		return
	}
	
	
	if reactionExist(cat, reaction) {
		reaction, err := q.WithContext(c.Request.Context()).Reactions.Where(q.Reactions.CatUUID.Eq(catUUID),
			q.Reactions.Emoji.Eq(reaction)).First()
		if cli, err := reaction.GetClients(); err != nil || contains(cli, fingerprint) {
			c.JSON(400, gin.H{"error": "already reacted"})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{"error": "reaction not found"})
			return
		}
		reaction.AppendClient(fingerprint)
		q.WithContext(c.Request.Context()).Reactions.Save(reaction)
	} else {
		clientJson, _ := json.Marshal([]string{fingerprint})
		err = q.WithContext(c.Request.Context()).Reactions.Create(&model.Reactions{
			CatUUID: catUUID,
			Emoji:   reaction,
			Clients: string(clientJson),
		})
		if err != nil {
			c.JSON(500, gin.H{"error": "reaction not created"})
			return
		}
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
	cat, err := q.WithContext(c.Request.Context()).Cats.Preload(q.Cats.Reactions).Where(q.Cats.UUID.Eq(catUUID)).First()
	if err != nil {
		c.JSON(400, gin.H{"error": "cat not found"})
		return
	}
	if !reactionExist(cat, reaction) {
		c.JSON(400, gin.H{"error": "reaction not found"})
		return
	}
	reactionObj, err := q.WithContext(c.Request.Context()).Reactions.Where(q.Reactions.CatUUID.Eq(catUUID),
		q.Reactions.Emoji.Eq(reaction)).First()
	if err != nil {
		c.JSON(400, gin.H{"error": "reaction not found"})
		return
	}

	if cli, err := reactionObj.GetClients(); err != nil || !contains(cli, fingerprint) {
		c.JSON(400, gin.H{"error": "no such record"})
		return
	}

	err = reactionObj.RemoveClient(fingerprint)
	if err != nil {
		c.JSON(500, gin.H{"error": "reaction remove error"})
		return
	}
	
	err = q.WithContext(c.Request.Context()).Reactions.Save(reactionObj)
	if err != nil {
		c.JSON(500, gin.H{"error": "reaction save error"})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func reactionExist(c *model.Cats, reaction string) bool {
	for _, r := range c.Reactions {
		if r.Emoji == reaction {
			return true
		}
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
