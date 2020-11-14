package main

import (
	"bytes"
	"github.com/WindGreen/pi/modules"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	brightness = kingpin.Flag("brightness", "brightness of display:0-7").Short('b').Default("3").Int()
	clkPin     = kingpin.Flag("clk", "clock pin no.").Default("2").Int()
	dioPin     = kingpin.Flag("dio", "data io pin no.").Default("3").Int()
	action     = kingpin.Flag("action", "if execute actions when battery change").Short('a').Default("false").Bool()
	minBattery = kingpin.Flag("min", "min battery percent").Default("10").Float64()
	maxBattery = kingpin.Flag("max", "max battery percent").Default("20").Float64()
	minCmd     = kingpin.Flag("mincmd", "action to execute").Default(`echo "min battery triggered"`).String()
	maxCmd     = kingpin.Flag("maxcmd", "action to execute").Default(`echo "max battery triggered"`).String()
	once       = kingpin.Flag("once", "action execute time").Default("true").Bool()
)

var (
	tm             *modules.TM1637
	minCmdExecuted int
	maxCmdExecuted int
)

func main() {
	kingpin.Parse()
	var err error
	tm, err = modules.OpenTM1637(*clkPin, *dioPin, *brightness)
	if err != nil {
		panic(err)
	}
	defer tm.Close()
	go func() {
		for {
			percent, _, timeLeft, battLeft, status := load()

			switch status {
			case "ONLINE":
				tm.Show([4]rune{'L', 'I', 0x52 | 0x80, 'E'}, false)
				time.Sleep(time.Second * 2)
			case "ONBATT":
				tm.Show([4]rune{'b', 'A', 0x31 | 0x80, 0x31 | 0x80}, false)
				time.Sleep(time.Second * 2)
			case "COMMLOST":
				tm.Show([4]rune{'L', 'O', 'S', 0x31 | 0x80}, false)
				time.Sleep(time.Second * 1)
				tm.Show([4]rune{'-', '-', '-', '-'}, false)
				time.Sleep(time.Second * 1)
				continue
			}

			result := [4]rune{'L', 0, 0, 0}
			p := int(math.Round(percent))
			for i := 3; i > 0; i-- {
				result[i] = rune(p % 10)
				p /= 10
			}
			if result[1] == 0 {
				result[1] = '.'
				if result[2] == 0 {
					result[2] = '.'
				}
			}
			tm.Show(result, false)
			time.Sleep(time.Second * 2)

			watt := int(math.Round(percent * 0.01 * 390))
			result = [4]rune{'E', 0, 0, 0}
			for i := 3; i > 0; i-- {
				result[i] = rune(watt % 10)
				watt /= 10
			}
			if result[1] == 0 {
				result[1] = '.'
			}
			tm.Show(result, false)
			time.Sleep(time.Second * 2)

			l := int(math.Round(battLeft))
			result = [4]rune{'b', 0, 0, 0}
			for i := 3; i > 0; i-- {
				result[i] = rune(l % 10)
				l /= 10
			}
			if result[1] == 0 {
				result[1] = '.'
			}
			tm.Show(result, false)
			time.Sleep(time.Second * 2)

			l = int(math.Round(timeLeft))
			result = [4]rune{0x31 | 0x80, '.', 0, 0}
			for i := 3; i > 1; i-- {
				result[i] = rune(l % 10)
				l /= 10
			}
			tm.Show(result, false)
			time.Sleep(time.Second * 2)

			//result = [4]rune{0, 0, 0, 'V'}
			//v := int(math.Round(volt))
			//for i := 2; i >= 0; i-- {
			//	result[i] = rune(v % 10)
			//	v /= 10
			//}
			//tm.Show(result, false)
			//time.Sleep(time.Second * 2)

			// schedule
			if *action {
				var cmd *exec.Cmd
				var action string
				if battLeft <= *minBattery {
					action = *minCmd
					if *once {
						if minCmdExecuted > 0 {
							continue
						}
						maxCmdExecuted = 0 //重置max计数
					}
					minCmdExecuted++
				} else if battLeft >= *maxBattery {
					action = *maxCmd
					if *once {
						if maxCmdExecuted > 0 {
							continue
						}
						minCmdExecuted = 0 //重置min计数
					}
					maxCmdExecuted++
				}
				cmd = exec.Command("sh", "-c", action)
				if out, err := cmd.Output(); err != nil {
					log.Printf("Action[%s]:%s\n", action, err)
				} else {
					log.Printf("Action[%s]:%s", action, string(out))
				}
			}
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func load() (loadPercent, volt, timeLeft, battLeft float64, state string) {
	info := status()

	// load
	load := info["LOADPCT"]
	str := strings.Split(load, " ")
	loadPercent, _ = strconv.ParseFloat(str[0], 64)

	// volt
	v := info["LINEV"]
	str = strings.Split(v, " ")
	volt, _ = strconv.ParseFloat(str[0], 64)

	// left
	l := info["TIMELEFT"]
	str = strings.Split(l, " ")
	timeLeft, _ = strconv.ParseFloat(str[0], 64)

	//status
	state = strings.TrimSpace(info["STATUS"])

	//charge percent
	b := info["BCHARGE"]
	str = strings.Split(b, " ")
	battLeft, _ = strconv.ParseFloat(str[0], 64)

	return
}

func status() map[string]string {
	cmd := exec.Command("/sbin/apcaccess", "status")
	out := bytes.Buffer{}
	cmd.Stdout = &out
	cmd.Run()
	str := out.String()

	lines := strings.Split(str, "\n")
	pairs := map[string]string{}
	for _, line := range lines {
		i := strings.Index(line, ":")
		if i == -1 || i == len(line)-1 {
			continue
		}
		k := strings.TrimSpace(line[:i])
		v := strings.TrimSpace(line[i+1:])
		pairs[k] = v
	}
	return pairs
}
