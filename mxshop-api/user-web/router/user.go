package router //路由
import(
	"github.com/gin-gonic/gin"
	
)
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")

	{
		UserRouter.Get("list", api.GetUserList)
	}
}
