package main

import (
	"log"
	"time"

	"github.com/WindGreen/pi/modules"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	str        = kingpin.Flag("string", "string content").Short('s').Default("HELLO").String()
	animation  = kingpin.Flag("animation", "animation of string:flow,wipe,flash,normal").Short('a').Default("normal").String()
	duration   = kingpin.Flag("duration", "duration of diget change").Short('d').Default("1000").Int()
	brightness = kingpin.Flag("brightness", "brightness of display:0-7").Short('b').Default("7").Int()
	clkPin     = kingpin.Flag("clk", "clock pin no.").Default("2").Int()
	dioPin     = kingpin.Flag("dlo", "data io pin no.").Default("3").Int()
	loop       = kingpin.Flag("loop", "loop").Short('l').Default("true").Bool()
	app        = kingpin.Arg("app", "sub command:time,string,light,clear").Required().String()
)

var (
	tm *modules.TM1637
)

func main() {
	kingpin.Parse()
	var err error
	tm, err = modules.OpenTM1637(*clkPin, *dioPin, *brightness)
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	switch *app {
	case "time":
		displayTime()
	case "string":
		displayString()
	case "clear":
		log.Println("clear")
		tm.Show([4]rune{'.', '.', '.', '.'}, false)
	case "light":
		log.Println("light")
		tm.Show([4]rune{8, 8, 8, 8}, true)
	}
}

func displayString() {
	log.Println("dispay string:", *str)

	frames := make([][4]rune, 0)
	switch *animation {
	case "normal":
		frame := [4]rune{'.', '.', '.', '.'}
		index := 0
		for i, r := range *str {
			frame[index] = r
			if index == 3 || i == len(*str)-1 {
				frames = append(frames, frame)
				index = 0
				frame = [4]rune{'.', '.', '.', '.'}
			} else {
				index++
			}
		}
	case "flow":
		frame := [4]rune{'.', '.', '.', '.'}
		for i := 0; i < len(*str)*2; i++ {
			for j := 0; j < 3; j++ {
				frame[j] = frame[j+1]
			}
			if i%2 == 0 {
				frame[3] = rune((*str)[i/2])
			} else {
				frame[3] = '.'
			}
			frames = append(frames, frame)
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				frame[j] = frame[j+1]
			}
			frame[3] = '.'
			frames = append(frames, frame)
		}
	}
	for i := 0; ; i++ {
		if i >= len(frames) {
			if *loop {
				i = 0
			} else {
				break
			}
		}
		tm.Show(frames[i], false)
		time.Sleep(time.Duration(*duration) * time.Millisecond)
	}
}

func displayTime() {
	log.Println("dislpay time")

	colon := true
	for {
		now := time.Now()
		h := now.Hour()
		m := now.Minute()
		a := h / 10
		b := h % 10
		c := m / 10
		d := m % 10
		tm.Show([4]rune{rune(a), rune(b), rune(c), rune(d)}, colon)
		colon = !colon
		time.Sleep(time.Second)
	}
}
