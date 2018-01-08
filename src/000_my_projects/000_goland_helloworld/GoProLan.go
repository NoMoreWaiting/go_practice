package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	htmltemp "html/template"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/cmplx"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	texttemp "text/template"
	"time"
	"unicode"
	"unicode/utf8"
)

func testGoProLan() {
	//glog_demo()
	//goplio_ch1_echo()
	//goplio_ch1_dup1()
	//goplio_ch1_dup2()
	//goplio_ch1_lissajous(os.Stdout)
	//goplio_ch1_fetch()
	//goplio_ch1_fetchall()
	//goplio_ch1_server()
	//goplio_ch2_echo()
	//goplio_ch2_cf()
	//goplio_ch3_surface()
	//goplio_ch3_mandelbrot(os.Stdout)
	//goplio_ch3_basename()
	//goplio_ch3_comma()
	//goplio_ch3_printints()
	//testSameStrOutOrder()
	//interConvStrNum()
	//goplio_ch3_netflag()
	//goplio_ch4_sha256()
	//testReverse()
	//testAppendInt()
	//traverseMap()
	//goplio_ch4_dedup()
	//testConvKey()
	//goplio_ch4_charcount()
	//goplio_ch4_wordfreq()
	//testBinaryTree()
	//goplio_ch4_movie()
	//goplio_ch4_issues()
	//goplio_ch4_issuesreport()
	//goplio_ch4_issueshtml()
	//testAnonymousfunc()
	goplio_ch5_topsort()
	//testVariablePara()
	//testJoin()
}

// glog 测试
func glog_demo() {
	// 命令行参数 -log_dir="./" -v=2
	flag.Parse() // 1 glog 解析命令行参数

	glog.Info("This is a Info log") // 2
	glog.Warning("This is a Warning log")
	glog.Error("This is a Error log")

	glog.V(1).Infoln("level 1") // 3
	glog.V(2).Infoln("level 2")
	glog.V(2).Info("%s  %s\n", "test ", "real down")
	//glog.V(3).Infoln("level 3")

	glog.Flush() // 4 输入到日志
}

// 命令行变量
func goplio_ch1_echo() {
	// 变量的声明方式
	// s := ""  // 只能在函数中使用不能导出, 并且注意作用域的问题
	// var s string
	// var s = ""
	// var s string = ""

	var s, sep string
	for i := 0; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	glog.Info(s)

	s = ""
	for i, arg := range os.Args[0:] {
		sep = strconv.Itoa(i) + ": "
		s += sep + arg + " "
	}
	glog.Info(s)
	glog.Info(strings.Join(os.Args[0:], " "))
}

// 查找重复的行, 从文件或者标准输入统计
func goplio_ch1_dup1() {
	//var i, j, k int
	//fmt.Scanln(&i, &j, &k)
	//glog.Info(i, j, k)

	counts := make(map[string]int) // 内置函数 make 创建空 map, 哈希map
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}

	for line, n := range counts {
		if n > 1 {
			/*
				%d 十进制整数
				%x, %o, %b 十六进制，八进制，二进制整数。
				%f, %g, %e 浮点数： 3.141593 3.141592653589793 3.141593e+00
				%t 布尔：true或false
				%c 字符（rune） (Unicode码点)
				%s 字符串
				%q 带双引号的字符串"abc"或带单引号的字符'c'
				%v 变量的自然形式（natural format）
				%T 变量的类型
				%% 字面上的百分号标志（无操作数）
			*/
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int) {
	// 需要使用终端去执行, goland的终端不一定靠谱
	// 这里的获取内容无法执行, 不要总是纠结暂时无意义的细节, 可以问人, 也可以以后解决, 调整心态

	const testText = "1,2,3,4,"
	input := bufio.NewScanner(f)
	//input := bufio.NewScanner(strings.NewReader(testText))
	for input.Scan() {
		counts[input.Text()]++
		if counts[input.Text()] > 1 {
			fmt.Println(f.Name())
		}

		//line := input.Text()
		//counts[line] += 1
	}
	// NOTE: ignoring potential errors from input.Err()
}

// 从文件中统计重复的行
func goplio_ch1_dup2() {
	counts := make(map[string]int)
	for _, filename := range os.Args[1:] {
		data, err := ioutil.ReadFile(filename) // ReadFile 函数返回一个字节切片（byte slice), 需要转化为string使用
		if nil != err {
			fmt.Fprintf(os.Stderr, "dump2: %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			counts[line]++
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

// gif图片
func goplio_ch1_lissajous(out io.Writer) {
	// 调色板数组
	var palette = []color.Color{color.RGBA{0xaa, 0xaa, 0xaa, 0xff}, color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff}}

	const (
		whiteIndex = 0 // first color in palette
		blackIndex = 1 // next color in palette
		blueIndex  = 2 // third color in palette
	)
	const (
		cycles  = 3     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 120   // number of animation frames  帧数
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				blueIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

// 获取url并打印信息, curl
func goplio_ch1_fetch() {
	// args: ./fetch http://www.baidu.com
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
			fmt.Println("url: ", url)
		}

		resp, err := http.Get(url)
		defer resp.Body.Close()
		if nil != err {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if nil != err {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(2)
		}
		x, err := io.Copy(os.Stderr, resp.Body)
		if nil != err {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(2)
		}
		fmt.Printf("%s\n", b)
		fmt.Printf("%s\n", x)
		fmt.Println(resp.Status) // http 返回的状态码
	}
}

// 同时获取所有URL
func goplio_ch1_fetchall() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // 猜测 ch, 在有数据库的时候才会打印, 没有数据就会堵塞
		/*
			当一个goroutine尝试在一个channel上做receive或者send操作时，这个goroutine会阻塞在调用处，直到另一个goroutine往这个channel里写入、或者接收值，这样两个goroutine才会继续执行channel操作之后的逻辑
			主线程也是一个goroutine, 所以堵塞在这里
		*/
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

// 获取URL, 传递获取时间和字节数
func fetch(url string, ch chan<- string) {
	start := time.Now()
	var nbytes int64
	for i := 0; i < 2; i++ {
		resp, err := http.Get(url)
		defer resp.Body.Close()
		if nil != err {
			ch <- fmt.Sprint(err) // send to channel ch
			return
		}
		nbytes, err = io.Copy(ioutil.Discard, resp.Body)
		if nil != err {
			ch <- fmt.Sprintf("while reading %s: %v", url, err)
			return
		}
		// 请求一个网站(百度)2次, 第二次应该在第一次的缓存上加了某些内容
		secs := time.Since(start).Seconds()
		fmt.Printf("-------- %.2fs  %7d  %s\n", secs, nbytes, url)
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
	fmt.Printf("%.2fs: , %s\n", secs, url)
}

var mu sync.Mutex
var count int

// 简单web服务器程序
func goplio_ch1_server() {
	http.HandleFunc("/", handlerDash) // 路径结尾以 "/"结尾都会匹配, "/xx/" 都可以
	http.HandleFunc("/count", handlerCounter)
	http.HandleFunc("/svg", handlerSVG)
	http.HandleFunc("/mandelbrot", handlerMandelbrot)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}

// 打印http请求信息
func handlerDash(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); nil != err {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
	/*
		URL.Path = "/"
		GET / HTTP/1.1
		Header["Connection"] = ["keep-alive"]
		Header["Upgrade-Insecure-Requests"] = ["1"]
		Header["User-Agent"] = ["Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"]
		Header["Accept"] = ["text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*\/*;q=0.8"]
		Header["Accept-Encoding"] = ["gzip, deflate, sdch, br"]
		Header["Accept-Language"] = ["zh-CN,zh;q=0.8"]
		Host = "127.0.0.1:8000"
		RemoteAddr = "127.0.0.1:40065"
	*/
}

// 打印http请求次数
func handlerCounter(w http.ResponseWriter, r *http.Request) {
	goplio_ch1_lissajous(w)
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

var heads int
var tail int

// switch 用例1
func switchDemo(str string) {
	switch str {
	case "heads":
		heads++
	case "tail":
		tail++
	case "other":
		fmt.Println("other")
		fallthrough // 不会break, 进入下一层判断
	default:
		fmt.Println("landed on edge")
	}
}

// switch 用例2
func sigNum(x int) int {
	switch { // 相当于switch true, 根据case中的值的true或false来判断执行
	case x > 0:
		return +1
	default: // x == 0
		return 0
	case x < 0:
		return -1
	}
}

var n = flag.Bool("n", false, "omit trailing newline")
var sep = flag.String("s", " ", "separator")

func goplio_ch2_echo() {
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), *sep))
	if !*n {
		fmt.Println()
	}
}

type Celsius float64
type Fahrenheit float64

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

// String() 方法
func (c Celsius) String() string    { return fmt.Sprintf("%g°C", c) }
func (f Fahrenheit) String() string { return fmt.Sprintf("%g°F", f) }

// 转换函数
// CToF converts a Celsius temperature to Fahrenheit.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC converts a Fahrenheit temperature to Celsius.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

/*
import (
"fmt"
"os"
"strconv"
"gopl.io/ch2/tempconv"
)
*/
// 温度转换
func goplio_ch2_cf() {
	for _, arg := range os.Args[1:] {
		t, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cf: %v\n", err)
			os.Exit(1)
		}
		f := Fahrenheit(t)
		c := Celsius(t)
		fmt.Printf("%s = %s, %s = %s\n",
			f, FToC(f), c, CToF(c))
	}
}

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

// 生成svg曲面图形
func goplio_ch3_surface() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Println("</svg>")
}

// 将svg图形写入 io.Writer
func goplio_ch3_surface_http(out io.Writer) {
	var svgpicture string
	svgpicture += fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			svgpicture += fmt.Sprintf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	svgpicture += fmt.Sprintln("</svg>")
	out.Write([]byte(svgpicture))
}

// 计算svg图形的点值
func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z
	z := heightSvg(x, y)

	// Compute (x, y, z) isometrically onto 2-D SVG canvas(sx, sy) // 投影到二维
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

// 计算svg图形的纵向高度
func heightSvg(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

// http请求svg图形 http.ResponseWriter 应该是一个引用类型
func handlerSVG(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml") // 默认以png解析, svg图形需要加此规则, 可以根据前面的512个字节自动输出对应的头部
	goplio_ch3_surface_http(w)
}

// complex128 复数算法生成一个Mandelbrot
func goplio_ch3_mandelbrot(out io.Writer) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(out, img)
}

// 生成分形图色彩
func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15
	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.RGBA{contrast * n, 255 - contrast*n, (255 - contrast*n) / 2, 0xff}
		}
	}
	return color.Black
}

// http 请求分形图
func handlerMandelbrot(w http.ResponseWriter, r *http.Request) {
	goplio_ch3_mandelbrot(w)
}

// 获取路径中文件的基本名
func goplio_ch3_basename() {
	s := "a/b/c.go"
	goplio_ch3_basename1(s)
	goplio_ch3_basename2(s)
}

// basename removes directory components and a .suffix.
// e.g., a => a, a.go => a, a/b/c.go => c, a/b.c.go => b.c
func goplio_ch3_basename1(s string) string {
	fmt.Println(s)
	// discard last '/' and everything before
	for i := len(s) - 1; i >= 0; i-- {
		if '/' == s[i] {
			s = s[i+1:]
			break
		}
	}
	// preserve everything before last '.'
	for i := len(s) - 1; i >= 0; i-- {
		if '.' == s[i] {
			s = s[:i]
			break
		}
	}
	fmt.Println(s)
	return s
}

// use strings.LastIndex to realize basename
func goplio_ch3_basename2(s string) string {
	fmt.Println(s)
	slash := strings.LastIndex(s, "/") // -1 if "/" not found
	s = s[slash+1:]
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	fmt.Println(s)
	return s
}

// 数字型字符串每隔3位插入.
func goplio_ch3_comma() {
	s := "123456789"
	fmt.Println(goplio_ch3_comma1(s))
	d := "1123456789.123456789"
	fmt.Println(goplio_ch3_comma2(d))
	fmt.Println(goplio_ch3_comma3(d))
}

// comma1 inserts commas in a non-negative decimal integer string
func goplio_ch3_comma1(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return goplio_ch3_comma1(s[:n-3]) + "," + s[n-3:]
}

// comma2 inserts commas in a non-negative decimal string
func goplio_ch3_comma2(s string) string {
	var pointString string
	if point := strings.LastIndex(s, "."); point >= 0 {
		pointString = s[point:]
		s = s[:point]
	}
	n := len(s)
	if n <= 3 {
		return s
	}
	return goplio_ch3_comma2(s[:n-3]) + "," + s[n-3:] + pointString
}

// comma3 inserts commas in a non-negative decimal string without recursion
func goplio_ch3_comma3(s string) string {
	var pointString string
	if point := strings.LastIndex(s, "."); point >= 0 {
		pointString = s[point:]
		s = s[:point]
	}
	n := len(s)
	var preS string
	if preNum := n % 3; preNum > 0 {
		preS = s[:preNum] + ","
		s = s[preNum:]
	}
	n = len(s)
	var buf bytes.Buffer
	buf.WriteString(preS)
	for i := 0; i < n; i = i + 3 {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(s[i : i+3])
	}
	buf.WriteString(pointString)
	return buf.String()
}

// practical func from strings bytes package
/*
strings package
	func Contains(s, substr string) bool
	func Count(s, sep string) int
	func Fields(s string) []string
	func HasPrefix(s, prefix string) bool
	func Index(s, sep string) int
	func Join(a []string, sep string) string
*/
/*
bytes package
	func Contains(b, subslice []byte) bool
	func Counts(s, sep []byte) int
	func Fields(s []byte) [][]byte
	func HasPrefix(s, prefix []byte) bool
	func Index(s, sep []byte) int
	func Join(s [][]byte, sep []byte) []byte
*/

// bytes.Buffer
func goplio_ch3_printints() {
	fmt.Println(intsToString([]int{1, 2, 3}))
}

// intsToString is like fmt.Sprint(values) but adds commas
func intsToString(values []int) string {
	// 当向bytes.Buffer添加任意字符的UTF8编码时，最好使用bytes.Buffer的WriteRune方法，WriteByte方法对于写入类似'['和']'等ASCII字符则更有效
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteByte(']')
	return buf.String()
}

// 测试 same string but out-of-order 函数
func testSameStrOutOrder() {
	s1 := "abceeef, 世界"
	s2 := "abeeecf, 世界"
	s3 := "eeeafcc, 世界"
	fmt.Printf("is same string but out of order?: %s, %s, %v\n", s1, s2, sameStrOutOrder(s1, s2))
	fmt.Printf("is same string but out of order?: %s, %s, %v\n", s1, s3, sameStrOutOrder(s1, s3))
}

// same string but out-of-order
func sameStrOutOrder(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	var buf1, buf2 []string
	for _, v := range s1 {
		buf1 = append(buf1, string(v))
	}
	for _, v := range s2 {
		buf2 = append(buf2, string(v))
	}
	sort.Strings(buf1)
	sort.Strings(buf2)
	return strings.Join(buf1, "") == strings.Join(buf2, "")
}

// 字符串和数字的转换
func interConvStrNum() {
	x := 123
	y := fmt.Sprint("%d", x)
	fmt.Println("y:", y, strconv.Itoa(x))
	fmt.Println(strconv.FormatInt(int64(x), 2))
	s := fmt.Sprintf("x=%b", x)
	fmt.Println("s: ", s)
	if x, err := strconv.Atoi("123"); nil == err {
		fmt.Println(x)
	}
	if y, err := strconv.ParseInt("123", 10, 64); nil == err {
		fmt.Println(y) // base 10, up to 64 bits
	}
	var x1, x2 int
	fmt.Sscanf("x=100, y=100", "x=%d, y=%d", &x1, &x2)
	fmt.Println(x1, x2)
}

type Flags uint

const (
	FlagUp           Flags = 1 << iota //is up
	FlagBroadcast                      // supports broadcast access capability
	FlagLoopback                       // is a loopback interface
	FlagPointToPoint                   // belongs to a point-to-point link
	FlagMulticast                      // supports multicast access capability
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	GiB // 1073741824
	TiB // 1099511627776 (exceeds 1 << 32)
	PiB // 1125899906842624
	EiB // 1152921504606846976
	ZiB // 1180591620717411303424 (exceeds 1 << 64)
	YiB // 1208925819614629174706176
)

func IsUp(v Flags) bool {
	return v&FlagUp == FlagUp
}

func TurnDown(v *Flags) {
	*v &^= FlagUp
}

func SetBroadcast(v *Flags) {
	*v |= FlagBroadcast
}

func IsCast(v Flags) bool {
	return v&(FlagBroadcast|FlagMulticast) != 0
}

// 网络标志位设置函数
func goplio_ch3_netflag() {
	var v Flags = FlagMulticast | FlagUp
	fmt.Printf("%b %t\n", v, IsUp(v)) // "10001 true"
	TurnDown(&v)
	fmt.Printf("%b %t\n", v, IsUp(v)) // "10000 false"
	SetBroadcast(&v)
	fmt.Printf("%b %t\n", v, IsUp(v))   // "10010 false"
	fmt.Printf("%b %t\n", v, IsCast(v)) // "10010 true"
	fmt.Println(YiB / ZiB)
}

// 生成sha256码
func goplio_ch4_sha256() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))
	fmt.Printf("%x\n%x\n%t\n%T\n%v\n", c1, c2, c1 == c2, c1, c1)
}

// 测试数组和切片
func testReverse() {
	// 数组
	// go中数组特殊, 数组的传值就可以直接改变数组元素, 相当于引用类型
	a := [...]int{0, 1, 2, 3, 4, 5}
	reverse(a[:]) // 这里传入的是切片, 然后使用切片改变了数组, 切片是引用类型
	fmt.Printf("%T, %v\n", a, a)

	// 切片 Rotate s left by two positions
	s := []int{0, 1, 2, 3, 4, 5}
	reverse(s[:2])
	reverse(s[2:])
	fmt.Printf("%T, %v\n", s, s)

	// 切片 测试旋转
	s = rotateCir(s, 1)
	fmt.Printf("%T, %v\n", s, s)

	// 数组
	arr := [...]int{0, 1, 2, 3, 4, 5}
	arr = reverseArray(arr, len(arr))
	fmt.Printf("%T, %v\n", arr, arr)

	// 切片
	str := []string{"value", "value", "value", "new"}
	str = eliminateRepeat(str)
	fmt.Printf("%T, %v\n", str, str)

}

// 翻转(传入切片类型), 切片通过引用改变了底层的数组
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// 数组作为参数, 必须标明数组类型 [xxx]int , 数组的长度是类型的一部分
// 函数参数值传递, 数组不会改变
func reverseArray(arr [6]int, size int) [6]int {
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// slice之间不能比较, 不能作为map的key, (要比较的话-->需要自定义比较函数, 不推荐, 并且不能作为map的键值)
// 原因1: 一个slice的元素是间接引用的，一个slice甚至可以包含自身
// 原因2: 一个固定值的slice在不同的时间可能包含不同的元素，因为底层数组的元素可能会被修改
func equalArr(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// appendInt 添加int元素进入slice
func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	if zlen <= cap(x) {
		// There is room to grow. Extend the slice.
		z = x[:zlen]
	} else {
		// There is insufficient space. Allocate a new array.
		// Grow by doubling, for amortized linear complexity.
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x) // a built-in function; see text
	}
	z[len(x)] = y
	return z
}

// 测试自定义的appendInt函数, 了解slice底层数据的扩展
func testAppendInt() {
	var x, y []int
	for i := 0; i < 10; i++ {
		y = appendInt(x, i)
		fmt.Printf("%2d cap=%3d\t\t%v\n", i, cap(x), y)
		x = y
	}
}

// rotate函数，通过一次循环完成旋转
func rotateCir(s []int, pos int) []int {
	if pos < 0 || pos > len(s) {
		return s
	}
	s1 := s[:pos]
	s = s[pos:]
	for _, x := range s1 {
		s = append(s, x)
	}
	return s
}

// 原地完成消除 []string 中相邻重复字符串
func eliminateRepeat(s []string) []string {
	var new = 0
	for _, value := range s {
		if value != s[new] {
			new++
			s[new] = value
		}
	}
	new++
	return s[:new]
}

// 按顺序遍历map
func traverseMap() {
	//ages := map[string]int{}
	ages := make(map[string]int)
	ages["alice"] = 22
	ages["charile"] = 30
	ages["lian"] = 50
	ages["zhoujielun"] = 37
	ages["hxm"] = 32
	ages["bob"] = ages["bob"] + 1 // 安全, 返回对应value类型的零值
	delete(ages, "hxm")
	fmt.Printf("%T, %v\n", ages, ages)
	//_ = &ages["bob"] // 禁止 map的值禁止取地址. map的值其实也是引用

	for name, age := range ages {
		fmt.Printf("name: %10s,\t\tage: %3d\n", name, age)
	}

	var names []string
	for name := range ages {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("name: %10s,\t\tage: %3d\n", name, ages[name])
	}
}

// 判断两个map是否相等
func equalMap(x, y map[string]int) bool {
	if len(x) != len(y) {
		return false
	}
	for k, xv := range x {
		if yv, ok := y[k]; !ok || yv != xv {
			return false
		}
	}
	return true
}

// 读取多行输入, 只打印第一次出现的行
func goplio_ch4_dedup() {
	seen := make(map[string]bool) // a set of strings
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if !seen[line] {
			seen[line] = true
			fmt.Println(line)
		}
	}
	if err := input.Err(); nil != err {
		fmt.Fprintf(os.Stderr, "dedup: %v\n", err)
		os.Exit(1)
	}
}

var m = make(map[string]int)

// 创造key值, %q参数忠实地记录每个字符串元素的信息
func k(list []string) string  { return fmt.Sprintf("%q", list) }
func Add(list []string)       { m[k(list)]++ }
func Count(list []string) int { return m[k(list)] }

// 使用转化技术将不可比较类型(比如slice)转化为string, 创造map
func testConvKey() {
	list1 := []string{"jack", "harry"}
	list2 := []string{"swift"}
	Add(list1)
	Add(list2)
	Add(list1)
	fmt.Println(m)
	fmt.Println(Count(list1), Count(list2))
}

// 统计输入中每个Unicode码点出现的次数
func goplio_ch4_charcount() {
	counts := make(map[rune]int)    // counts of Unicode characters
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters

	file, err := os.Open("config.xml")
	if nil != err {
		log.Fatal(err)
	}
	in := bufio.NewReader(file) // 读取文件
	//in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune()
		if io.EOF == err {
			break
		}
		if nil != err {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}

		// 如果输入是无效字符, 那么返回 unicode.ReplacementChar, 并且长度为1
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

// 嵌套map
var graph = make(map[string]map[string]bool)

//函数惰性初始化map
func addEdge(from, to string) {
	edges := graph[from]
	if nil == edges {
		edges = make(map[string]bool)
		graph[from] = edges
	}
	edges[to] = true
}

func hasEdge(from, to string) bool {
	return graph[from][to]
}

// 输入文本中每个单词出现的频率
func goplio_ch4_wordfreq() {
	counts := map[string]int{}

	file, err := os.Open("config.xml")
	if nil != err {
		log.Fatal(err)
	}
	in := bufio.NewScanner(file)
	in.Split(bufio.ScanWords) // 设置分割函数
	for in.Scan() {
		counts[in.Text()]++
	}
	fmt.Println(counts)
	for i, v := range counts {
		fmt.Printf("%s: %d\n", i, v)
	}
}

// 结构体实现二叉树用来插入排序
type tree struct {
	value int // 小写, 不可被导出
	left  *tree
	right *tree
}

// Sort sorts values in place
func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	fmt.Println(appendValues(values[:0], root))
}

// appendValues appends the elements of t to values in order and returns the resulting slice
// 中序遍历
func appendValues(values []int, t *tree) []int {
	if nil != t {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

// add node to the tree
func add(t *tree, value int) *tree {
	if nil == t {
		// Equivalent to return &tree{value: value}
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

// 中序遍历
func show(t *tree) {
	if nil == t {
		return
	}
	show(t.left)
	fmt.Println(t.value)
	show(t.right)
}

// 测试二叉树
func testBinaryTree() {
	var root *tree
	for i := 0; i < 10; i = i + 3 {
		root = add(root, i)
	}
	for i := 1; i < 10; i = i + 3 {
		root = add(root, i)
	}

	show(root)

	// 结构体面值
	fmt.Println("结构体面值:")
	var test1 = tree{6, nil, nil}
	var test2 = tree{left: nil, right: nil, value: 6}
	fmt.Println(test1, test2)
	fmt.Println("test1 == test2 : ", test1 == test2)

	values := []int{1, 4, 5, 3, 2}
	Sort(values)
}

// json 格式
func goplio_ch4_movie() {
	type Movie struct {
		Title  string
		Year   int  `json:"released"`
		Color  bool `json:"color,omitempty"`
		Actors []string
	}
	var movies = []Movie{
		{Title: "Casablanca", Year: 1942, Color: false,
			Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
		{Title: "Cool Hand Luke", Year: 1967, Color: true,
			Actors: []string{"Paul Newman"}},
		{Title: "Bullitt", Year: 1968, Color: true,
			Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
		// ...
	}
	// 返回编码后的字节slice
	data, err := json.Marshal(movies)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)

	// 带缩进的json格式
	data, err = json.MarshalIndent(movies, "", " ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)

	var titles []struct{ Title string }
	if err = json.Unmarshal(data, &titles); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}
	fmt.Println(titles) // "[{Casablanca} {Cool Hand Luke} {Bullitt}]"
}

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}
type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *GUser
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}
type GUser struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

// 查询Go语言项目中和JSON解码相关的问题
func goplio_ch4_issues() {
	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
}

// 模板 text/template
const templ = `{{.TotalCount}} issues:
{{range .Items}}----------------------------------------
Number: {{.Number}}
User: {{.User.Login}}
Title: {{.Title | printf "%.64s"}}
Age: {{.CreatedAt | daysAgo}} days
{{end}}`

// 模板 html/template
const htmlTempl = `
<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-align: left'>
<th>#</th>
<th>State</th>
<th>User</th>
<th>Title</th>
</tr>
{{range .Items}}
<tr>
<td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
<td>{{.State}}</td>
<td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
<td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>
`

// 通过time.Since函数将CreatedAt成员转换为过去的时间长度
func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

var report = texttemp.Must(texttemp.New("issuelist").
	Funcs(texttemp.FuncMap{"daysAgo": daysAgo}).
	Parse(templ))

var issueList = htmltemp.Must(htmltemp.New("issuelist").Parse(htmlTempl))

// 参数 repo:golang/go is:open json decoder
// 测试模板参数打印报告
func goplio_ch4_issuesreport() {
	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := report.Execute(os.Stdout, result); nil != err {
		log.Fatal(err)
	}
}

// 参数 repo:golang/go commenter:gopherbot json encoder >issues.html
// 测试html模板参数打印报告
func goplio_ch4_issueshtml() {
	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := issueList.Execute(os.Stdout, result); nil != err {
		log.Fatal(err)
	}
}

//使用 golang.org/x/net/html package, 需要翻墙下载
//递归遍历整个HTML结点树，并输出树的结构
/*
func goplio_ch5_outline() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
		os.Exit(1)
	}
	outline(nil, doc)
}

func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data) // push tag
		fmt.Println(stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}
*/

/*
func goplio_ch5_findlinks2() {
	for _, url := range os.Args[1:] {
		links, err := findLinks(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "findlinks2: %v\n", err)
			continue
		}
		for _, link := range links {
			fmt.Println(link)
		}
	}
}

// findLinks performs an HTTP GET request for url, parses the
// response as HTML, and extracts and returns the links.
func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}

// visit appends to links each link found in n and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

// bare return 一个函数将所有的返回值都显示的变量名，那么该函数的return语句可以省略操作数
// CountWordsAndImages does an HTTP GET request for the HTML
// document url and returns the number of words and images in it.
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}
func countWordsAndImages(n *html.Node) (words, images int) {
// ...
}

*/

// 错误处理机制, 传播错误, 重试, 输出错误信息并结束, 输出错误信息, 忽略
// WaitForServer attempts to contact the server of a URL.
// It tries for one minute using exponential back-off.
// It reports an error if all attempts fail.
func WaitForServer(url string) error {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil // success
		}
		log.Printf("server not responding (%s);retrying…", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}

// gopl.io/ch5/outline2
// forEachNode针对每个结点x,都会调用pre(x)和post(x)。
// pre和post都是可选的。
// 遍历孩子结点之前,pre被调用
// 遍历孩子结点之后，post被调用
/*
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

var depth int

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
		depth++
	}
}
func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}

// 调用
//forEachNode(doc, startElement, endElement)
*/

// 匿名函数  go语言圣经 5.6节, 需要重点学习
//strings.Map(func(r rune) rune { return r + 1 }, "HAL-9000")

// squares返回一个匿名函数。
// 该匿名函数每次被调用时都会返回下一个数的平方。
func squares() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}

// 闭包
func testAnonymousfunc() {
	f := squares()   // f在赋值时, squares函数中的x为0, 以后每次f(), 就在x的基础上++. x是squares()第一次调用是传递过来的, 所以值会累加, 但在此调用squares不会累加
	fmt.Println(f()) // "1"
	fmt.Println(f()) // "4"
	fmt.Println(f()) // "9"
	fmt.Println(f()) // "16"

	// 变量的生命周期不由它的作用域决定：squares返回后，变量x仍然隐式的存在于f中

	ff := squares() // 新的匿名函数
	fmt.Println(ff())
	fmt.Println(ff())
	fmt.Println(ff())
	fmt.Println(ff())
}

//给定一些计算机课程，每个课程都有前置课程，只有完成了前置课程才可以开始当前课程的学习
// 闭包, 拓扑排序
// prereqs记录了每个课程的前置课程
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

// 拓扑排序
func goplio_ch5_topsort() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

//深度优先搜索整张图，获得符合要求的课程序列
func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)

	//匿名函数需要被递归调用时，必须首先声明一个变量，再将匿名函数赋值给这个变量。如果不分成两部，函数字面量无法与visitAll绑定，无法递归调用该匿名函数
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item]) // 一层一层向深处递归遍历, 找到最后的不依赖其他课程的最前置课程. 然后在递归返回的时候, 就append进了反向顺序
				fmt.Println(item)
				order = append(order, item)
			}
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order
}

/*
// Package links provides a link-extraction function.
package links
import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
)
// Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}
*/
// 广度优先算法。调用者需要输入一个初始的待访问列表和一个函数f
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool) // 避免同一个元素被访问两次
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

//crawl函数会将URL输出，提取其中的新链接，并将 这些新链接返回
/*
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}
*/

// 运行抓取器, 参数 https://golang.org
/*
func spiderHtml() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl, os.Args[1:])
}
*/

// 警告:捕获迭代变量
// 在循环中, 使用循环的迭代变量时要注意: 在整个循环体作用域中, 共享同一个循环变量, 使用闭包(函数值)中记录的是循环变量的内存地址，而不是循环变量某一时刻的值, 会导致循环迭代失效(都是操作最后一次结果)
// 闭包(函数值), go语句, defer语句, (range的循环, 循环变量i)会经常遇到此类问题, 这些都会等待循环结束再执行函数值

func tempDirs() []string {
	var dirs []string
	dirs = append(dirs, "/temp/test")
	return dirs
}

// 循环作用域陷阱
func goplio_ch4_561() {
	var rmdirs []func()
	for _, d := range tempDirs() {
		dir := d // NOTE: necessary!
		// creates parent directories too
		os.MkdirAll(dir, 0755)
		rmdirs = append(rmdirs, func() {
			os.RemoveAll(dir)
		})
	}
	// ...do some work…
	for _, rmdir := range rmdirs {
		rmdir() // clean up
	}
}

// 可变参数
func sum(vals ...int) int {
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}

func ff(...int) {}
func gg([]int)  {}

func errorf(linenum int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Line %d: ", linenum)
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
}

// 测试可变参数
func testVariablePara() {

	fmt.Println(sum())  // "0"
	fmt.Println(sum(3)) // "3"
	// 隐式创建原始数据的数组, 再将数据的切片传入被调用函数
	fmt.Println(sum(1, 2, 3, 4)) // "10"

	fmt.Printf("%T\n", ff) // "func(...int)"
	fmt.Printf("%T\n", gg) // "func([]int)"

	linenum, name := 12, "count"
	errorf(linenum, "undefined: %s", name) // "Line 12: undefined: count"
}

// 可变参数版本 join
func multiJoin(seg string, values ...string) string {
	var s string
	for i, str := range values {
		if i > 0 {
			s += seg
		}
		s += str
	}
	return s
}

// 可变参数版本 join, 使用 bytes.Buffer
func muitiJoinBuf(seg string, values ...string) string {

	// from string.
	//var r io.Reader = strings.NewReader(string("hello, world"))
	// from bytes.
	//var r io.Reader = bytes.NewReader([]byte("hello, world!"))
	// from bytes buffer.
	//var r io.Reader = bytes.NewBuffer([]byte("hello, world"))
	// buffer reader
	//var r io.Reader = bufio.NewReader(strings.NewReader(string("hello, world")))
	//func (b *Reader) ReadString(delim byte) (line string, err error)

	var buf bytes.Buffer
	for i, str := range values {
		if i > 0 {
			buf.WriteString(seg)
		}
		buf.WriteString(str)
	}
	return buf.String()
}

// 测试可变参数版本Join
func testJoin() {
	fmt.Println(multiJoin(",-- ", "123", "789", "456"))
	fmt.Println(muitiJoinBuf(",-- ", "123", "789", "456"))

}

// bufio.NewScanner 扫描例子
func bufioParseJson() {
	s := bufio.NewScanner(strings.NewReader("/*block comments*///line comments\n{}"))
	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// read more.
		return 0, nil, nil
	})
	for s.Scan() {
		fmt.Println(s.Text())
	}
	fmt.Println("err is", s.Err())

}

// T类型的值不拥有所有*T指针的方法, 使用T类型的变量调用*T类型的方式是语法糖, 编译器给隐式取址了, 如果是不可取值的值(InSet{}.String()), 那么就会编译失败

// 一个类型如果拥有一个接口需要的所有方法，那么这个类型就实现了这个接口
// 接口类型就是约定了函数调用方式(参数, 返回值和函数名一样), 只关心类型可以做什么, 不关心其内部的实现. 就是需要一个打字的, 传过来一个码农可以, 传过来打字员也可以.
// interface{}被称为空接口类型, 对实现它的类型没有要求，所以可以将任意一个值赋给空接口类型, 空接口类型值不能直接操作, 需要使用类型断言来获取interface{}中值

// 接口值，由两个部分组成，一个具体的类型和那个类型的值, 称为接口的动态类型和动态值
