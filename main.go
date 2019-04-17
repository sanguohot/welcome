package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanguohot/welcome/etc"
	"github.com/sanguohot/welcome/pkg/common/log"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

var (
	dir = "frontend"
	host = "0.0.0.0"
	port = 8443
)

func noRouteHandler(c *gin.Context) {
	p := filepath.Join(dir, c.Request.RequestURI)
	info, err := os.Stat(p)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	if info.IsDir() {
		log.Logger.Error(fmt.Sprintf("dir %s not supported", p))
	}
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	if strings.HasSuffix(p, ".json") {
		c.Data(http.StatusOK, "application/json", data)
	} else {
		c.Data(http.StatusOK, http.DetectContentType(data), data)
	}
}

func main() {
	//http.ListenAndServe(":4200", nil)
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())
	// 默认设置logger，但启用logger会导致吞吐量大幅度降低
	if os.Getenv("GIN_LOG") != "off" {
		r.Use(gin.Logger())
	}
	r.MaxMultipartMemory = 10 << 20 // 10 MB
	r.Static("/", filepath.Join(etc.GetServerDir(), dir))
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", host, port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Logger.Fatal(err.Error())
		}
	}()
	log.Sugar.Infof("[http] listening => %s, serv => %s", server.Addr, dir)
	// apiserver发生错误后延时五秒钟，优雅关闭
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Logger.Fatal(err.Error())
	}
	log.Sugar.Infof("stop server => %s, serv => %s", server.Addr, dir)

	//http.Handle("/data/*", &dataHandler{})
	//http.Handle("/", http.FileServer(http.Dir(dir)))
	//
	//log.Sugar.Infof("正在监听4200, dir => %s", dir)
	//http.ListenAndServe(":4200", nil)
}