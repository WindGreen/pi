package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"time"
	"tickpay/notify/email"
	"fmt"
	"strings"
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
		if temp >= alarmNum {
			log.Printf("Temp is too high:%d\n", temp)
			e := email.Zoho{
				To:      []string{"admin@tickpay.org"},
				Subject: "pi's temperature is too high",
				Content: fmt.Sprintf(`pi's temperature is <b style="color:red">%d</b>,large than %d`, temp, alarmNum),
			}
			e.Send()
		} else {
			log.Println(temp)
		}
		time.Sleep(time.Duration(duration) * time.Second)
	}
}