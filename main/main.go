package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/latifrons/gows"
	"github.com/sirupsen/logrus"
	"time"
)

func handleFunc(client *gows.Client, bytes []byte) {
	fmt.Println(client.Conn.RemoteAddr().String())
	fmt.Println(string(bytes))
}

func main() {
	r := gin.New()
	r.Use(gows.GinLogger(logrus.StandardLogger(), time.RFC3339, true))

	hub := gows.NewHub()
	go hub.Run(false)

	r.GET("/ws", func(c *gin.Context) {
		gows.ServeWs(hub, c.Writer, c.Request, handleFunc)
	})

	go func() {
		for {
			t := time.Now().String()
			// send message.
			hub.Broadcast([]byte(t))
			time.Sleep(time.Second)
		}
	}()

	r.Run("0.0.0.0:8000")
}
