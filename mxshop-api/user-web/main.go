package main

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin/binding"
	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/global"
	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/initialize"
	myvalidator "github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/validator"

	"go.uber.org/zap"
)

func main() {
	// logger,_:=zap.NewProduction()//json格式
	//1.初始化logger
	initialize.InitLogger()
	//2.初始化配置文件
	initialize.InitConfig()
	//3.初始化routers
	Router := initialize.Routers()
	//4.初始化翻译
	if err:= initialize.InitTrans("zh");err!=nil{
		panic(err)
	}
	//注册验证器
	if v,ok:=binding.Validator.Engine().(*validator.Validate);ok{
		_=v.RegisterValidation("mobile",myvalidator.ValidateMobile)
		_=v.RegisterTranslation("mobile",global.Trans,func(ut ut.Translator)error{
			return ut.Add("mobile","{0}非法的手机号！",true)
		},func(ut ut.Translator,fe validator.FieldError)string{
			t,_:=ut.T("mobile",fe.Field())
			return t
		})
	}
	/*
		1.S()可以获取全局的suger，可以让我们自己设置一个全局的logger
		2.日志分级别的，debug, info, warn, error, fetal
		3.S函数和L函数很有用，提供了一个全局的安全访问logger的途径
	*/

	// logger,_:=zap.NewProduction()
	// defer logger.Sync()
	// suger:=logger.Suager()
	//或者？！
	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)

	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}

}
