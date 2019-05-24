package modules

import (
	"errors"

	"github.com/WindGreen/pi"
)

var TM1637_ = map[rune]byte{
	0:   0x0,
	'a': 0x01,
}

const (
	TM1637_ADDR_AUTO  = 0X40
	TM1637_ADDR_FIXED = 0X44
	TM1637_START_ADDR = 0XC0
)

type TM1637 struct {
	data       [4]byte
	clkPin     *pi.Pin
	dataPin    *pi.Pin
	brightness float64
}

func OpenTM1637(clkPin, dataPin int) (tm *TM1637, err error) {
	tm = new(TM1637)
	tm.clkPin, err = pi.OpenPin(clkPin, pi.OUT)
	if err != nil {
		return nil, err
	}
	tm.dataPin, err = pi.OpenPin(dataPin, pi.OUT)
	if err != nil {
		return nil, err
	}
	return tm, nil
}

func (tm *TM1637) WriteRune(pos int, r rune) error {
	if pos < 0 || pos > 3 {
		return errors.New("tm1637 index out of range")
	}
	tm.start()
	tm.
	return nil
}

func (tm *TM1637) Write(d []byte) (int, error) {

	return 4, nil
}

func (tm *TM1637) Close() error {

	return nil
}

func (tm *TM1637) Clear() error {

	return nil
}

func (tm *TM1637) flush() error {

}

func (tm *TM1637) start() error {
	tm.clkPin.Write(pi.HIGH)
	tm.dataPin.Write(pi.HIGH)
	tm.dataPin.Write(pi.LOW)
	tm.clkPin.Write(pi.LOW)
	return nil
}

func (tm *TM1637) stop() error {
	tm.clkPin.Write(pi.LOW)
	tm.dataPin.Write(pi.LOW)
	tm.clkPin.Write(pi.HIGH)
	tm.dataPin.Write(pi.HIGH)
	return nil
}

func (tm *TM1637)writeByte(d byte)error{
	for i:=0;i<8;i++{
		tm.clkPin.Write(pi.LOW)
		if d&0x01{
			tm.dataPin.Write(pi.HIGH)
		}else{
			tm.data.Pin.Write(pi.LOW)
		}
		tm.clkPin.Write(pi.HIGH)
	}
	tm.clkPin.Write(pi.LOW)
	tm.dataPin.Write(pi.HIGH)
	tm.clk
	tm.dataPin.SetDirection(pi.IN)

	tm.dataPin.SetDirection(pi.OUT)
	return nil
}