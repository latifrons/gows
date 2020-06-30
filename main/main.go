package main

import (
	"github.com/gin-gonic/gin"
	"github.com/latifrons/gows"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	r := gin.New()
	r.Use(gows.GinLogger(logrus.StandardLogger(), time.RFC3339, true))

	hub := wshandler.NewHub()
	go hub.Run(false)

}
