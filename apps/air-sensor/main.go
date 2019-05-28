package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/WindGreen/pi/modules"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	ioPin    = kingpin.Flag("pin", "io pin number").Short('p').Default("4").Int()
	duration = kingpin.Flag("duration", "duration").Short('d').Default("5").Int()
	output   = kingpin.Flag("output", "tm1637,2,3, 2-clk pin 3-dio pin").String()
	app      = kingpin.Arg("app", "ht").Required().String()
)

func main() {
	kingpin.Parse()
	dht, err := modules.OpenDHT11(*ioPin)
	if err != nil {
		panic(err)
	}
	log.Println("detect dht11 on", *ioPin)
	defer dht.Close()
	time.Sleep(time.Second)
	for {
		h, t, err := dht.Read()
		if err != nil {
			log.Println("error", err)
		} else {
			Print(h, t)
		}
		time.Sleep(time.Second * time.Duration(*duration))
	}
}

func Print(h, t float64) {
	log.Println("temperature", t, "humidity", h)
	if *output == "" {

	} else {
		args := strings.Split(*output, ",")
		if len(args) == 0 {
			panic("--output is not valid")
		}
		if args[0] == "tm1637" {
			clkPin, err := strconv.Atoi(args[1])
			if err != nil {
				panic(err)
			}
			dioPin, err := strconv.Atoi(args[2])
			if err != nil {
				panic(err)
			}
			tm, err := modules.OpenTM1637(clkPin, dioPin, 7)
			if err != nil {
				panic(err)
			}
			tm.Show([4]rune{rune(t / 10), rune(int(t) % 10), rune(h / 10), rune(int(h) % 10)}, false)
		}
	}
}
