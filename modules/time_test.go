package modules

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestPrint(t *testing.T) {
	t1 := time.Now()
	log.Println("hello")
	t2 := time.Now()
	fmt.Println("hello")
	t3 := time.Now()
	fmt.Println("t2-t1", t2.Sub(t1), "t3-t2", t3.Sub(t2), "t3-t1", t3.Sub(t1))
	// output t2-t1 22.9682ms t3-t2 546.5µs t3-t1 23.5147ms
}

func TestSleep(t *testing.T) {
	t1 := time.Now()
	<-time.After(time.Microsecond)
	t2 := time.Now()
	fmt.Println("chan", t2.Sub(t1))
	//output: chan 1.0511ms
	t3 := time.Now()
	time.Sleep(time.Nanosecond)
	t4 := time.Now()
	time.Sleep(time.Microsecond)
	t5 := time.Now()
	time.Sleep(time.Millisecond)
	t6 := time.Now()
	time.Sleep(time.Second)
	t7 := time.Now()
	fmt.Println("nano", t4.Sub(t3), "micro", t5.Sub(t4), "milli", t6.Sub(t5), "second", t7.Sub(t6), "all", t7.Sub(t3))
	// output: nano 2.1539ms micro 1.0383ms milli 1.7096ms second 1.0002482s all 1.00515s
	t8 := time.Now()
	t9 := time.Now()
	fmt.Println("now cost", t9.Sub(t8))
	// now cost 0s
}

func TestLoop(t *testing.T) {
	t1 := time.Now()
	for i := 0; i < 10000; i++ {
	}
	t2 := time.Now()
	for i := 0; i < 100000; i++ {
	}
	t3 := time.Now()
	for i := 0; i < 1000000; i++ {
	}
	t4 := time.Now()
	for i := 0; i < 10000000; i++ {
	}
	t5 := time.Now()
	fmt.Println("t2-t1", t2.Sub(t1), "t3-t2", t3.Sub(t2), "t4-t3", t4.Sub(t3), "t5-t4", t5.Sub(t4))
	// output: t2-t1 0s t3-t2 0s t4-t3 513.3µs t5-t4 5.9848ms
	// output: t2-t1 0s t3-t2 1.0018ms t4-t3 0s t5-t4 3.9832ms
}

func TestFileTime(t *testing.T) {

}
