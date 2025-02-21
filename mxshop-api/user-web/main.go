package main
import(
	"github.com/gin-gonic/gin"
)
func main(){
	router:=gin.Default()

	Apigroup:=router.Group("/v1")//v1版本号
	router2.InitUserRouter(Apigroup)
}