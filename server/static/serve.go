package static

import (
	"catalog-go/config"

	"github.com/gin-gonic/gin"
)

func ConfigureRoute(r *gin.Engine) {
	r.Static("/static", config.CONFIG.Store.StorePath)
}
