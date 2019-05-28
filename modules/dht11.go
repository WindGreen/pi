package modules

import (
	"errors"
	"time"

	"github.com/WindGreen/pi"
)

type DHT11 struct {
	ioPin *pi.Pin
}

func OpenDHT11(ioPin int) (*DHT11, error) {
	d := new(DHT11)
	pin, err := pi.OpenPin(ioPin)
	if err != nil {
		return nil, err
	}
	d.ioPin = pin
	return d, nil
}

func (d *DHT11) Read() (hum, tmp float64, err error) {
	d.ioPin.Set(pi.OUT)
	d.ioPin.Write(pi.HIGH)
	d.ioPin.Write(pi.LOW)
	time.Sleep(time.Millisecond * 20) // >= 18 ms
	d.ioPin.Write(pi.HIGH)            // cost about 15~20 us, Println cost about 150us
	d.ioPin.Set(pi.IN)                // cost about 18 us
	// time.Sleep(time.Microsecond * 20) // 20-40 us
	for {
		v, _ := d.ioPin.Read() // cost about 32 us
		if v != pi.LOW {
			break
		}
	}
	for {
		v, _ := d.ioPin.Read()
		if v != pi.HIGH {
			break
		}
	}

	data := [5]byte{}
	index := 0
	for i := 0; i < 40; i++ {
		for {
			v, _ := d.ioPin.Read()
			if v != pi.LOW {
				break
			}
		}
		st := time.Now()
		for {
			v, _ := d.ioPin.Read()
			if v != pi.HIGH {
				break
			}
		}
		dr := time.Now().Sub(st)
		if dr <= 50*time.Microsecond { //28us + 32us(compensate)
			data[index] <<= 1
		} else {
			data[index] = data[index]<<1 | 0x01
		}
		if i%8 == 7 {
			index++
		}
	}
	a := int(data[0]) + int(data[1]) + int(data[2]) + int(data[3])
	b := int(data[4])
	if a != b {
		return 0, 0, errors.New("check failed")
	}
	h := int(data[0]) // int(data[0])<<52 | int(data[1])<<44 //preserved
	t := int(data[2]) // int(data[0])<<52 | int(data[1])<<44 //preserved
	hum = float64(h)
	tmp = float64(t)
	return hum, tmp, nil
}

func (d *DHT11) Close() error {
	d.ioPin.Close()
	return nil
}
