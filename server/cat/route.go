package cat

import "github.com/gin-gonic/gin"

func ConfigureRoute(r *gin.Engine) {
	r.GET("/api/cat/get", getHandler)
	r.GET("/api/cat/get_by_id", getByIdHandler)
	r.GET("/api/cat/add_reaction", addReactionHandler)
	r.GET("/api/cat/remove_reaction", removeReactionHandler)
}
