package initialize

import(
	"go.uber.org/zap"
)
func InitLogger(){
	logger,_:=zap.NewDevelopment()//日志格式
	zap.ReplaceGlobals(logger)

}