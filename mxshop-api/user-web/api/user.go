package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/form"
	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/global"
	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/global/reponse"
	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/proto"
)

// 处理嵌套结构体的名称
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for filed, err := range fileds {
		field := ""
		rsp[filed[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

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

// 处理验证器错误
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
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

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
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

func PassWordLogin(c *gin.Context) {
	//表单验证
	passwordLoginForm := form.PassWordLoginForm{}
	if err := c.ShouldBindJSON(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
		//如何返回错误信息
		//ch——07？？
	}
	if err := c.ShouldBindJSON(&passwordLoginForm); err != nil {

	}
	//拨号连接grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure()) //指定不加密   insecure——不上锁的
	if err != nil {
		zap.S().Errorw("[GetUserList]连接【用户服务失败】",
			"msg", err.Error())
	}
	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	//登录的逻辑
	if rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		} else {
			//只是查找到了用户而已，并没有检查密码
			if passRsp, passErr := userSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
				PassWord:          passwordLoginForm.PassWord,
				EncryptedPassword: rsp.PassWord,
			}); passErr != nil {
				c.JSON(http.StatusInternalServerError, map[string]string{
					"password": "登录失败",
				})
			} else {
				if passRsp.Success {
					c.JSON(http.StatusOK, map[string]string{
						"msg": "登录成功",
					})
				} else {
					c.JSON(http.StatusBadRequest, map[string]string{
						"msg": "登录失败",
					})
				}
			}

		}
	}
}
