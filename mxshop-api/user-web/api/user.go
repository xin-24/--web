package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/global"
	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/global/reponse"
	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/proto"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() { //拿到状态码
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"mss": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code,
				})
			}
			return

		}
	}
}

func GetUserList(ctx *gin.Context) {

	//拨号连接grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure()) //指定不加密   insecure——不上锁的
	if err != nil {
		zap.S().Errorw("[GetUserList]连接【用户服务失败】",
			"msg", err.Error())
	}
	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)

	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    0,
		PSize: 0,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//zap.S().Debug("获取用户列表")
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		// data := make(map[string]interface{})

		user := reponse.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			// BirthDay:time.Time(time.Unix(int64(value.BirthDay),0)).Format("2025-2-25"),//法二
			// BirthDay:time.Time(time.Unix(int64(value.BirthDay),0)),//法一
			BirthDay: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)), //法三
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		// data["id"] = value.Id
		// data["name"] = value.NickName
		// data["birthday"] = value.BirthDay
		// data["gender"] = value.Gender
		// data["mobile"] = value.Mobile

		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}
