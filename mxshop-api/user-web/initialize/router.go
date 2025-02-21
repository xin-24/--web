package initialize
//1-4/1-5
import (
	"github.com/gin-gonic/gin"
	"github.com/xin-24/go/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default

	Apigroup := router.Group("/v1") //v1版本号
	router.InitUserRouter(Apigroup)
	return Router
}
