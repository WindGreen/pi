package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/WindGreen/pi"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	rgb = kingpin.Flag("rgb", "rgb").Default("false").Bool()
	r   = kingpin.Flag("red", "read").Short('r').Bool()
	g   = kingpin.Flag("green", "green").Short('g').Bool()
	b   = kingpin.Flag("blue", "blue").Short('b').Bool()
	pin = kingpin.Flag("pin", "pin No.").Default("13,19,26").Short('p').String()

	cmd = kingpin.Arg("action", "on|flash").Default("on").String()
)

func main() {
	kingpin.Parse()

	if *rgb {
		pins := strings.Split(*pin, ",")
		if len(pins) != 3 {
			log.Fatal("pins like 13,19,26")
			return
		}
		p1, e1 := strconv.Atoi(pins[0])
		p2, e2 := strconv.Atoi(pins[1])
		p3, e3 := strconv.Atoi(pins[2])
		if e1 != nil || e2 != nil || e3 != nil {
			log.Fatal("pin number parse failed")
		}
		pin1, err := pi.OpenPin(p1)
		if err != nil {
			panic(err)
		}
		defer pin1.Close()
		pin2, err := pi.OpenPin(p2)
		if err != nil {
			panic(err)
		}
		defer pin2.Close()
		pin3, err := pi.OpenPin(p3)
		if err != nil {
			panic(err)
		}
		defer pin3.Close()

		pin1.Set(pi.OUT)
		pin2.Set(pi.OUT)
		pin3.Set(pi.OUT)

		var ct int32
		if *cmd == "on" {
			ct = 1
		} else if *cmd == "flash" {
			ct = 1 << 30
		}
		for i := 0; i < int(ct); i++ {
			if *r {
				pin1.Write(pi.HIGH)
			}
			if *g {
				pin2.Write(pi.HIGH)
			}
			if *b {
				pin3.Write(pi.HIGH)
			}
			time.Sleep(time.Second)
			if ct > 1 {
				pin1.Write(pi.LOW)
				pin2.Write(pi.LOW)
				pin3.Write(pi.LOW)
				time.Sleep(time.Second)
			}
		}
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
