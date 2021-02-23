package main

import (
	"entry_task/usermanagementsys/src/conf"
	"entry_task/usermanagementsys/src/rpcclient"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/astaxie/beego/core/logs"
	"github.com/gin-gonic/gin"
)

var config conf.HTTPConf

func init() {
	// parse the confFile into HTTPConf struct
	var confFile string
	flag.StringVar(&confFile, "c", "../conf/httpserver.yaml", "config file")
	flag.Parse()
	err := conf.ConfParser(confFile, &config)
	if err != nil {
		logs.Critical("Parser config failed, err:", err.Error())
		os.Exit(-1)
	} else {
		logs.Critical("Parser config succ")
	}

	// init log
	logConfig := fmt.Sprintf(`{"filename":"%s","level":%s,"maxlines":0,"maxsize":0,"daily":true,"maxdays":%s}`,
		config.Log.Logfile, config.Log.Loglevel, config.Log.Maxdays)
	logs.SetLogger(logs.AdapterFile, logConfig)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.Async()

	// init userclient (pool) with the parsed config struct
	err = rpcclient.InitPool(config.Rpcserver.Addr, config.Pool.Initsize, config.Pool.Capacity, time.Duration(config.Pool.Maxidle)*time.Second)
	if err != nil {
		logs.Critical("InitPool failed, err:", err.Error())
		os.Exit(-2)
	} else {
		logs.Critical("Init userclient Pool succ")
	}
}

// cleanup global objects
func finalize() {
	rpcclient.DestoryPool()
}

// start a gin instance and register handlers
func main() {
	defer finalize()

	gin.SetMode(gin.DebugMode)
	// gin.DefaultWriter = ioutil.Discard

	r := gin.Default()
	// home page
	r.Any("/welcome", webRoot)
	r.POST("/login", loginHandler)
	r.POST("/logout", logoutHandler)
	r.GET("/getuserinfo", getUserinfoHandler)
	r.POST("/editnickname", editNicknameHandler)
	r.POST("/uploadpic", uploadHeadurlHandler)

	r.Static("/static/", "./static/")
	r.Static("/upload/images/", "./upload/images/")
	fmt.Println("about to start")
	// listen to http server port : 8080
	r.Run(fmt.Sprintf(":%d", config.Server.Port))
	fmt.Println("started")
}

func webRoot(context *gin.Context) {
	context.String(http.StatusOK, "gin works fine")
}
