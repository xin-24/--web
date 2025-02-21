//可能要删除不然会有问题

package main

import (
	
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction() //生产环境//json
	// logger,_:=zap.NewDevelopment()//开发环境//日志
	defer logger.Sync() // flushes buffer, if any
	url := "https//imooc.com"
	logger.Info("failed to fetch URL",
	zap.String("url",url),
	zap.Int("nums",3),
)
	// sugar := logger.Sugar()
	// sugar.Infow("failed to fetch URL",//输出时用的
	// 	// Structured context as loosely typed key-value pairs.
	// 	"url", url,
	// 	"attempt", 3,
	// 	"backoff", time.Second,
	// )
	// sugar.Infof("Failed to fetch URL: %s", url)
}
