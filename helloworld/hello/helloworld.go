package hello

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

// ---------------------------------
const LIMIT = 30

func testDefer() {
	println(DeferFunc1(1))
	println(DeferFunc2(1))
	println(DeferFunc3(1))
}

func DeferFunc1(i int) (t int) {
	t = i
	defer func() {
		t += 3
	}()
	return t
}

func DeferFunc2(i int) int {
	t := i // t的左右与仅为内部函数, 不是整个 DeferFunc2, 不会影响外界的返回值
	defer func() {
		t += 3
	}()
	return t
}

func DeferFunc3(i int) (t int) {
	defer func() {
		t += i
	}()
	return 2
}

const (
	x  = iota // 0
	y         // 1
	z  = "zz" // zz
	kk        // zz
	p  = iota // 4
)

// func() int 是函数 fibonacci() 的返回值, 返回值是一个 返回一个int类型值的函数
func fibonacci() func() int {
	back1, back2 := 0, 1
	return func() int {
		// 重新赋值
		back1, back2 = back2, back1+back2
		return back1
	}
}

func demoFibonacci() {
	f := fibonacci() // 返回一个闭包函数
	var array [LIMIT]int
	for i := 0; i < LIMIT; i++ {
		array[i] = f()
	}
	fmt.Println(array)
}

// 求最大公约数
func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

// 计算斐波那契数列  // 返回第n个菲波那契数
func fib(n int) int {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		x, y = y, x+y
	}
	return x
}

func testPrint() {
	fmt.Println(x, y, z, kk, p)
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5}
	s1 = append(s1, s2...) // 如果是 append(s1, s2) 报错. 不能直接合并
	fmt.Println(s1)
}

func testHelloWorld() {
	fmt.Println("hello world")
	fmt.Println("hello world!", runtime.Version())
	demoFibonacci()
	fmt.Printf("第 %d 个斐波那契数: %d \n", 8, fib(8))
	fmt.Printf("%d 和 %d 的最大公约数是: %d \n", 136, 36, gcd(136, 88))
}

func testTimeSwitch() {
	// 获取本地 location
	// 注意这里的时间两边不能有空格, 要去空格
	toBeCharge := "2017-11-08 14:07:31"                             //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                                            //转化为时间戳 类型是int64
	fmt.Println(theTime)                                            //打印输出theTime 2015-01-01 15:15:00 +0800 CST
	fmt.Println(theTime.Unix())
	fmt.Println("sr  ", sr) //打印输出时间戳 1420041600

	//时间戳转日期
	dataTimeStr := time.Unix(sr, 0).Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
	fmt.Println(dataTimeStr)

	// 获取时间戳
	t := time.Now()
	fmt.Println(t)
	fmt.Println(t.UTC().Format(time.UnixDate))
	fmt.Println(t.Unix())
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	fmt.Println(timestamp)
	timestamp = timestamp[:10]
	fmt.Println(timestamp)
}

// 周期任务机制
func testPeriodTask() {
	ticker := time.NewTicker(time.Millisecond * 100)
	ch := make(chan int, 10)
	go func() {
		i := 0
		for value := range ticker.C {
			fmt.Printf("ticked at %v\n", time.Now())
			fmt.Println("value =", value)
			ch <- value.Second()
			i++
			if i > 10 {
				close(ch)
				break
			}
		}
	}()
	for value := range ch {
		fmt.Println("value ---", value)
	}
}

func TestHelloWorld() {
	testPrint()
	testDefer()
	testHelloWorld()
	testTimeSwitch()
	testPeriodTask()
}
