package cat

import "github.com/gin-gonic/gin"

func ConfigureRoute(r *gin.Engine) {
	r.GET("/api/cat/get_all", getAllHandler)
	r.GET("/api/cat/get", getHandler)
	r.GET("/api/cat/get_by_id", getByIdHandler)
}
