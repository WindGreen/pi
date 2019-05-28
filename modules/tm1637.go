package modules

import (
	"errors"
	"time"

	"github.com/WindGreen/pi"
)

var TM1637_MAP = map[rune]byte{
	0:   0x3f,
	1:   0x06,
	2:   0x5b,
	3:   0x4f,
	4:   0x66,
	5:   0x6d,
	6:   0x7d,
	7:   0x07,
	8:   0x7f,
	9:   0x6f,
	'0': 0x3f,
	'1': 0x06,
	'2': 0x5b,
	'3': 0x4f,
	'4': 0x66,
	'5': 0x6d,
	'6': 0x7d,
	'7': 0x07,
	'8': 0x7f,
	'9': 0x6f,
	'.': 0x00,
	'A': 0x77,
	'b': 0x7c,
	'C': 0x39,
	'c': 0x58,
	'd': 0x5e,
	'E': 0x79,
	'F': 0x71,
	'g': 0x6f,
	'H': 0x76,
	'h': 0x7c,
	'I': 0x06,
	'J': 0x0d,
	'L': 0x38,
	'l': 0x06,
	'O': 0x3f,
	'o': 0x5c,
	'P': 0x73,
	'p': 0x73,
	'q': 0x67,
	'r': 0x50,
	'S': 0x6d,
	's': 0x6d,
	'U': 0x3e,
	'u': 0x1c,
	'V': 0x3e,
	'v': 0x1c,
}

const (
	TM1637_ADDR_AUTO  = 0x40
	TM1637_ADDR_FIXED = 0x44
	TM1637_START_ADDR = 0xC0
)

type TM1637 struct {
	colon      byte //colon number
	clckPin    *pi.Pin
	dataPin    *pi.Pin
	brightness int
}

func OpenTM1637(clkPin, dataPin, brightness int) (tm *TM1637, err error) {
	tm = new(TM1637)
	if brightness < 0 || brightness > 7 {
		return nil, errors.New("brightness out of define")
	}
	tm.brightness = brightness
	tm.clckPin, err = pi.OpenPin(clkPin)
	if err != nil {
		return nil, err
	}
	tm.dataPin, err = pi.OpenPin(dataPin)
	if err != nil {
		return nil, err
	}
	return tm, nil
}

func (tm *TM1637) Show(f [4]rune, colon bool) error {
	for i, r := range f {
		tm.WriteRune(i, r)
	}
	if colon {
		tm.Colon(true)
	} else {
		tm.Colon(false)
	}
	return nil
}

func (tm *TM1637) Clear() error {
	tm.Show([4]rune{'.', '.', '.', '.'}, false)
	return nil
}

func (tm *TM1637) WriteRune(pos int, r rune) error {
	if pos < 0 || pos > 3 {
		return errors.New("tm1637 index out of range")
	}
	tm.WriteByte(pos, TM1637_MAP[r])
	return nil
}

func (tm *TM1637) WriteByte(pos int, b byte) error {
	if pos < 0 || pos > 3 {
		return errors.New("tm1637 index out of range")
	}
	tm.start()
	tm.writeByte(TM1637_ADDR_FIXED)
	tm.br()
	tm.writeByte(TM1637_START_ADDR | byte(pos))
	tm.writeByte(b)
	tm.br()
	tm.writeByte(0x88 | byte(tm.brightness))
	tm.stop()
	if pos == 1 {
		tm.colon = b
	}
	return nil
}

func (tm *TM1637) Write(d []byte) (int, error) {

	return 4, nil
}

func (tm *TM1637) Colon(b bool) error {
	if b {
		tm.colon |= 0x80
	} else {
		tm.colon &= 0x7f
	}
	tm.WriteByte(1, tm.colon)
	return nil
}

func (tm *TM1637) Close() error {
	tm.clckPin.Close()
	tm.dataPin.Close()
	return nil
}

func (tm *TM1637) flush() error {
	return nil
}

func (tm *TM1637) start() error {
	tm.clckPin.Write(pi.HIGH)
	tm.dataPin.Write(pi.HIGH)
	tm.dataPin.Write(pi.LOW)
	tm.clckPin.Write(pi.LOW)
	return nil
}

func (tm *TM1637) stop() error {
	tm.clckPin.Write(pi.LOW)
	tm.dataPin.Write(pi.LOW)
	tm.clckPin.Write(pi.HIGH)
	tm.dataPin.Write(pi.HIGH)
	return nil
}

func (tm *TM1637) br() error {
	tm.stop()
	tm.start()
	return nil
}

func (tm *TM1637) writeByte(d byte) error {
	for i := 0; i < 8; i++ {
		tm.clckPin.Write(pi.LOW)
		if d&0x01 == 1 {
			tm.dataPin.Write(pi.HIGH)
		} else {
			tm.dataPin.Write(pi.LOW)
		}
		d >>= 1
		tm.clckPin.Write(pi.HIGH)
	}

	//ack
	tm.clckPin.Write(pi.LOW)
	tm.dataPin.Write(pi.HIGH)
	tm.clckPin.Write(pi.HIGH)
	tm.dataPin.Set(pi.IN)

	for {
		v, err := tm.dataPin.Read()
		if err != nil {
			panic(err)
		}
		if v == pi.LOW {
			// tm.clckPin.Write(pi.HIGH)
			break
		} else {
			// tm.dataPin.SetDirection(pi.OUT)
			// tm.dataPin.Write(pi.LOW)
			// tm.dataPin.SetDirection(pi.IN)
			time.Sleep(time.Millisecond)
		}
	}
	tm.dataPin.Set(pi.OUT)
	return nil
}
