package router //路由
import (
	"github.com/gin-gonic/gin"
	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/api"

	"go.uber.org/zap"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关url")
	{
		UserRouter.GET("list", api.GetUserList)
		UserRouter.POST("pwd_Login", api.PassWordLogin)
	}
}
