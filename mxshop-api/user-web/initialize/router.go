package initialize

//1-4/1-5
import (
	"github.com/gin-gonic/gin"
	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	Apigroup := Router.Group("/u/v1") //v1版本号
	router.InitUserRouter(Apigroup)
	return Router
}
