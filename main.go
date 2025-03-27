package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
)

var args struct {
	ListenAddr string `arg:"-l,--listen" help:"Address to listen on" default:"127.0.0.1:8085"`
	GpioPin    int    `arg:"-p,--pin" help:"GPIO pin to use" default:"19"`
}

var logger = logrus.StandardLogger()

func main() {
	arg.MustParse(&args)
	e := gin.Default()
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	e.Use(cors.New(cfg))
	e.POST("/api/v1/set", setLight)
	if err := e.Run(args.ListenAddr); err != nil {
		logger.Fatalf("error starting server: %v", err)
	}
}

type SetLightReq struct {
	On         bool `json:"on"`
	Brightness int  `json:"brightness"`
}

func setLight(context *gin.Context) {
	var req SetLightReq
	if err := context.BindJSON(&req); err != nil {
		context.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	var err error
	if !req.On {
		err = setBrightness(0)
	} else {
		if req.Brightness < 0 || req.Brightness > 32 {
			context.JSON(400, gin.H{"error": "invalid brightness"})
			return
		}
		err = setBrightness(req.Brightness)
	}
	if err != nil {
		logger.Errorf("error setting brightness: %v", err)
		context.JSON(500, gin.H{"error": "error setting brightness"})
		return
	}
	context.JSON(200, gin.H{"success": true})
}

func setBrightness(brightness int) error {
	if err := rpio.Open(); err != nil {
		return fmt.Errorf("error opening GPIO: %v", err)
	}
	defer rpio.Close()

	pin := rpio.Pin(args.GpioPin)
	pin.Mode(rpio.Pwm)
	pin.Freq(64000)
	pin.DutyCycle(uint32(brightness), 32)
	return nil
}
