package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	file     string
	alarmNum int
	duration int
)

func main() {
	flag.StringVar(&file, "f", "/sys/class/thermal/thermal_zone0/temp", "temperature file of raspbian")
	flag.IntVar(&alarmNum, "a", 80, "alert temp")
	flag.IntVar(&duration, "d", 5, "duration of seconds")
	flag.Parse()
	for {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println(err)
			continue
		}
		temp, _ := strconv.Atoi(strings.TrimSpace(string(data)))
		temp /= 1000
		log.Printf("Temperature:%d\n", temp)
		time.Sleep(time.Duration(duration) * time.Second)
	}
}
