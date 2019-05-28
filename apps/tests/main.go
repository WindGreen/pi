package main

/*
	time.Sleep(time.Microsecond) cost about 10)us but not stable(26-440us)
	C.usleep(1) cost about 72us but not stable(88-175us)
	usleep(1) cost about 65-75us stable
*/

/*
#include <stdio.h>
#include <sys/time.h>
#include <time.h>
#include <unistd.h>

void testbase(){
	int i;
	struct timeval t1,t2,r21;

	gettimeofday(&t1, NULL);
	gettimeofday(&t2, NULL);
	timersub(&t2, &t1, &r21);
	printf("gettime:%ld\n", r21.tv_usec);
	// 1
	// 0

	gettimeofday(&t1, NULL);
	for (i=0;i<1000000;i++){
	}
	gettimeofday(&t2, NULL);
	timersub(&t2, &t1, &r21);
	printf("loop:%ld\n", r21.tv_usec);
	// 0
	// 0
	// loop cost => 0

	gettimeofday(&t1, NULL);
	for (i=0;i<1000;i++){
		gettimeofday(&t2, NULL);
	}
	timersub(&t2, &t1, &r21);
	printf("sub1:%ld\n", r21.tv_usec);
	// sub1:247
	// sub1:249
	// sub1:247
	// => cost 0.25us

	gettimeofday(&t1, NULL);
	for (i=0;i<1000;i++){
		usleep(1);
	}
	gettimeofday(&t2, NULL);
	timersub(&t2, &t1, &r21);
	printf("sub2:%ld\n", r21.tv_usec);
	//sub2:65525
	//sub2:65180
	//suub2:65142
	//=> cost 65us
}

void testusleep(){
	struct timeval t1,t2,t3,t4,t5,r21,r32,r43,r54;
	gettimeofday(&t1, NULL);
	usleep(1);
	gettimeofday(&t2, NULL);
	usleep(10);
	gettimeofday(&t3, NULL);
	usleep(100);
	gettimeofday(&t4, NULL);
	usleep(1000);
	gettimeofday(&t5, NULL);

	timersub(&t2, &t1, &r21);
	timersub(&t3, &t2, &r32);
	timersub(&t4, &t3, &r43);
	timersub(&t5, &t4, &r54);

	printf("1:%ld, 10:%ld, 100:%ld, 1000:%ld\n", r21.tv_usec,r32.tv_usec,r43.tv_usec,r54.tv_usec);
	// 1:76, 10:309, 100:174, 1000:1072
	// 1:74, 10:81, 100:169, 1000:1071
	// => cost 70us
}

void testclock(){
	struct timespec t1,t2,t3,t4,t5;
	clock_gettime(CLOCK_REALTIME, &t1);
	usleep(1);
	clock_gettime(CLOCK_REALTIME, &t2);
	usleep(10);
	clock_gettime(CLOCK_REALTIME, &t3);
	usleep(100);
	clock_gettime(CLOCK_REALTIME, &t4);
	usleep(1000);
	clock_gettime(CLOCK_REALTIME, &t5);

	printf("1:%ld, 10:%ld, 100:%ld, 1000:%ld\n",t2.tv_nsec-t1.tv_nsec,t3.tv_nsec-t2.tv_nsec,t4.tv_nsec-t3.tv_nsec,t5.tv_nsec-t4.tv_nsec);
	// 1:74218, 10:79791, 100:169218, 1000:1071399
}

void testfile(){
	struct timeval t1,t2,r;
	FILE *f = NULL;
	char c='\0';
	char buf[8];
	char out[4]="out\0";
	f=fopen("test.txt","w+");
	if(f==NULL){
		return;
	}

	gettimeofday(&t1, NULL);
	fputc('0',f);
	gettimeofday(&t2, NULL);
	timersub(&t2, &t1, &r);
	printf("fputc:%ld\n",r.tv_usec);

	gettimeofday(&t1, NULL);
	fputs("in",f);
	gettimeofday(&t2, NULL);
	timersub(&t2, &t1, &r);
	printf("fputs:%ld\n",r.tv_usec);

	fclose(f);


	f=fopen("test.txt","r");

	gettimeofday(&t1, NULL);
	c=fgetc(f);
	gettimeofday(&t2, NULL);
	timersub(&t2, &t1, &r);
	printf("fgetc:%ld,value:%c\n",r.tv_usec,c);

	gettimeofday(&t1, NULL);
	fgets(buf,8,f);
	gettimeofday(&t2, NULL);
	timersub(&t2, &t1, &r);
	printf("fgets:%ld,value:%s\n",r.tv_usec,buf);

	fclose(f);


	f=fopen("test.txt","wb");

	gettimeofday(&t1, NULL);
	fwrite(out,sizeof(out),3,f);
	gettimeofday(&t2, NULL);
	timersub(&t2, &t1, &r);
	printf("fwrite:%ld,value:%s\n",r.tv_usec,out);

	fclose(f);


	f=fopen("test.txt","wb");

	gettimeofday(&t1, NULL);
	fread(buf,8,3,f);
	gettimeofday(&t2, NULL);
	timersub(&t2, &t1, &r);
	printf("fread:%ld,value:%s\n",r.tv_usec,buf);

	fclose(f);


	//fputc:25
	//fputs:4
	//fgetc:23,value:0
	//fgets:6,value:in
	//fwrite:12,value:out
	//fread:21,value:in

	//fputc:26
	//fputs:4
	//fgetc:20,value:0
	//fgets:7,value:in
	//fwrite:7,value:out
	//fread:18,value:in

	//fputc:52
	//fputs:5
	//fgetc:19,value:0
	//fgets:7,value:in
	//fwrite:9,value:out
	//fread:18,value:in
}
 */
import "C"

import (
	"fmt"
	"log"
	"time"

	//"github.com/WindGreen/pi"
)

func main() {
	TestBase()
	TestPrint()
	TestSleep()
	TestLoop()
	TestUsleep()
	TestCusleep()
	//TestPiSleep()
}

func TestBase() {
	var st, en time.Time

	st = time.Now()
	for i := 0; i < 1000; i++ {
	}
	en = time.Now()
	fmt.Println("loop 1000:", en.Sub(st))
	// 6.198µs
	// 6.146µs

	st = time.Now()
	for i := 0; i < 1000000; i++ {
	}
	en = time.Now()
	fmt.Println("loop 10000:", en.Sub(st))
	// 5.031234ms
	// 5.104463ms

	st = time.Now()
	for i := 0; i < 1000; i++ {
		en = time.Now()
	}
	fmt.Println("sub1", en.Sub(st))
	// 2.043638ms
	// 2.092232ms
	// => time.Now cost 2.1us

	st = time.Now()
	for i := 0; i < 1000; i++ {
		C.usleep(1)
	}
	en = time.Now()
	fmt.Println("sub2", en.Sub(st))
	// 72.231384ms
	// 71.292327ms
	// => C.usleep() cost 72us

	st = time.Now()
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Microsecond)
	}
	en = time.Now()
	fmt.Println("sub3", en.Sub(st))
	// 9.334863ms
	// 9.154885ms
	// 9.391638ms
	// => 9.3us
}

func TestPrint() {
	t1 := time.Now()
	log.Println("hello")
	t2 := time.Now()
	fmt.Println("hello")
	t3 := time.Now()
	fmt.Println("t2-t1", t2.Sub(t1), "t3-t2", t3.Sub(t2), "t3-t1", t3.Sub(t1))
	//windows output t2-t1 22.9682ms t3-t2 546.5µs t3-t1 23.5147ms
	//raspi output t2-t1 702.707µs t3-t2 358.072µs t3-t1 1.060779ms
}

func TestSleep() {
	t1 := time.Now()
	<-time.After(time.Microsecond)
	t2 := time.Now()
	fmt.Println("chan", t2.Sub(t1))
	//windows output: chan 1.0511ms
	//linux output: chan 227.447µs
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
	// raspi: nano 172.708µs micro 9.062µs milli 1.106769ms second 1.000152885s all 1.001441424s
	// 		nano 149.119µs micro 12.553µs milli 1.111183ms second 1.000168266s all 1.001441121s

	d10 := time.Microsecond * 10
	d100 := time.Microsecond * 100
	d1000 := time.Microsecond * 1000
	t3 = time.Now()
	time.Sleep(time.Microsecond)
	t4 = time.Now()
	time.Sleep(d10)
	t5 = time.Now()
	time.Sleep(d100)
	t6 = time.Now()
	time.Sleep(d1000)
	t7 = time.Now()
	fmt.Println("1us", t4.Sub(t3), "10us", t5.Sub(t4), "100us", t6.Sub(t5), "1000us", t7.Sub(t6), "all", t7.Sub(t3))
	// d*n
	// 1us 26.146µs 10us 125.833µs 100us 207.238µs 1000us 1.111141ms all 1.470358ms
	// 1us 68.229µs 10us 159.999µs 100us 209.27µs 1000us 1.112547ms all 1.550045ms
	// 1us 67.291µs 10us 127.135µs 100us 205.207µs 1000us 1.10911ms all 1.508743ms
	// dn
	// 1us 24.584µs 10us 104.744µs 100us 190.424µs 1000us 1.110254ms all 1.430006ms
	// 1us 442.986µs 10us 113.911µs 100us 212.821µs 1000us 1.109834ms all 1.879552ms
	// 1us 388.868µs 10us 184.381µs 100us 206.83µs 1000us 1.106186ms all 1.886265ms
	// 1us 29.792µs 10us 119.796µs 100us 607.312µs 1000us 1.286709ms all 2.043609ms
	// => time.Sleep() not stable

	t8 := time.Now()
	t9 := time.Now()
	fmt.Println("now cost", t9.Sub(t8))
	//now cost 1.25µs

}

func TestLoop() {
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
	//windows output: t2-t1 0s t3-t2 0s t4-t3 513.3µs t5-t4 5.9848ms
	//windows output: t2-t1 0s t3-t2 1.0018ms t4-t3 0s t5-t4 3.9832ms

	//raspi output: t2-t1 353.072µs t3-t2 543.228µs t4-t3 5.040458ms t5-t4 50.494008ms
	//raspi output: t2-t1 51.205µs t3-t2 501.105µs t4-t3 5.001735ms t5-t4 50.101938ms
}

func TestUsleep() {
	t1 := time.Now()
	C.usleep(1)
	t2 := time.Now()
	C.usleep(10)
	t3 := time.Now()
	C.usleep(100)
	t4 := time.Now()
	C.usleep(1000)
	t5 := time.Now()
	fmt.Println("1us:", t2.Sub(t1), "10us:", t3.Sub(t2), "100us", t4.Sub(t3), "1000us", t5.Sub(t4))
	//raspi output: 1us: 160.314µs 10us: 85.573µs 100us 175.679µs 1000us 1.075581ms
	//				1us: 88.076µs 10us: 88.076µs 100us 198.704µs 1000us 1.12285ms
	//				1us: 86.615µs 10us: 395.262µs 100us 218.542µs 1000us 1.078493ms
}

func TestCusleep() {
	C.testbase()
	C.testusleep()
	C.testclock()
	C.testfile()
}

//func TestPiSleep() {
//	t1 := time.Now()
//	pi.Sleep(1)
//	t2 := time.NOw()
//	pi.Sleep(10)
//	t3 := time.NOw()
//	pi.Sleep(100)
//	t4 := time.NOw()
//	pi.Sleep(1000)
//	t5 := time.NOw()
//	fmt.Println("1us:", t2.Sub(t1), "10us:", t3.Sub(t2), "100us", t4.Sub(t3), "1000us", t5.Sub(t4))
//	// 1us: 1.979µs 10us: 7.761µs 100us 67.604µs 1000us 667.55µs
//	// 1us: 2.031µs 10us: 7.812µs 100us 67.604µs 1000us 720.936µs
//	// 1us: 1.979µs 10us: 7.812µs 100us 69.271µs 1000us 740.155µs 
//}
