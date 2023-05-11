package main

import (
	"fmt"
	"os"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	common.DefaultReadConfig(os.Args)
	server.ConfigureLogging()

	server.RawGet("/ping", func(ctx *gin.Context) {
		logrus.Infof("Received ping from '%v'", ctx.RemoteIP())
		ctx.Data(200, "text/plain", []byte(fmt.Sprintf("pong at %s\n", time.Now())))
	})

	server.BootstrapServer()
}
