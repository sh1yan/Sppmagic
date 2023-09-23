package logger

import (
	"fmt"
	"github.com/gookit/color"
	fc "main/function"
	"os"
	"path"
	"regexp"
	"runtime"
	"time"
)

var ( // 输出颜色设置
	Red         = color.Red.Render
	Cyan        = color.Cyan.Render
	Yellow      = color.Yellow.Render
	White       = color.White.Render
	Blue        = color.Blue.Render
	Purple      = color.Style{color.Magenta, color.OpBold}.Render
	LightRed    = color.Style{color.Red, color.OpBold}.Render
	LightGreen  = color.Style{color.Green, color.OpBold}.Render
	LightWhite  = color.Style{color.White, color.OpBold}.Render
	LightCyan   = color.Style{color.Cyan, color.OpBold}.Render
	LightYellow = color.Style{color.Yellow, color.OpBold}.Render
)

var (
	defaultLevel = LevelWarning // 输出等级
	noWrite      int            // 不进行写入log
	Num          int64
	End          int64
	LogSucTime   int64
	LogErrTime   int64
	WaitTime     int64
)

// SetLevel 设置输出等级
func SetLevel(l Level) {
	defaultLevel = l
}

func getCallerInfo(skip int) (info string) {
	_, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		info = "runtime.Caller() failed"
	}

	fileName := path.Base(file) // Base函数返回路径的最后一个元素
	return fmt.Sprintf("%s line:%d", fileName, lineNo)
}

// log 日志输出效验判断
func log(l Level, nw int, detail string) {
	switch fc.LogLevel { // 判断配置文件中的日志输出等级，并设置到日志等级
	case 0:
		SetLevel(0)
	case 1:
		SetLevel(1)
	case 2:
		SetLevel(2)
	case 3:
		SetLevel(3)
	case 4:
		SetLevel(4)
	case 5:
		SetLevel(5)
	}

	if l > defaultLevel { // 判断输入等级是否大于设置等级，若大于则当前日志则不进行输出
		return
	}

	if nw == 0 { // 在该项目中，只有logwriteinfo函数才进行写入，剩下的均不写入本地文件
		// 目前只写入 info 信息 和 Success 的信息
		//strTrim := fmt.Sprintf("[%s] ", Cyan(getDate())) // 匹配log日志前面的日期
		//detail := strings.TrimPrefix(detail, strTrim)    // 去除日期信息
		writeLogFile(clean(detail), fc.OutputFileName) // 写入到本地文件
		return
	}

	if fc.NoColor { // 判断是否关闭颜色输出
		fmt.Println(clean(detail))
		return
	} else {
		fmt.Println(detail)
	}

	if l == LevelFatal {
		os.Exit(0)
	}
}

// Fatal 严重级别日志 log等级：0
func Fatal(detail string) {
	noWrite = 1
	log(LevelFatal, noWrite, fmt.Sprintf("[%s] [%s] %s", Cyan(getDate()), LightRed("FATAL"), detail))
}

// LogWriteInfo log写入本地消息日志 log等级：2
func LogWriteInfo(detail string) {
	noWrite = 0
	log(LevelInfo, noWrite, detail)
}

// Error 错误日志 log等级：1
func Error(detail string) {
	noWrite = 1
	log(LevelError, noWrite, fmt.Sprintf("[%s] [%s] %s", Cyan(getDate()), LightRed("ERROR"), detail))
}

// Info 消息日志 log等级：2
func Info(detail string) {
	noWrite = 1
	log(LevelInfo, noWrite, fmt.Sprintf("[%s] [%s] %s", Cyan(getDate()), LightGreen("INFO"), detail))
}

// Warning 告警日志 log等级：3
func Warning(detail string) {
	noWrite = 1
	log(LevelWarning, noWrite, fmt.Sprintf("[%s] [%s] %s", Cyan(getDate()), LightYellow("WARNING"), detail))
}

// Debug 调试日志 log等级：4
func Debug(detail string) {
	noWrite = 1
	log(LevelDebug, noWrite, fmt.Sprintf("[%s] [%s] [%s] %s", Cyan(getDate()), LightWhite("DEBUG"), Yellow(getCallerInfo(2)), detail))
}

// Verbose 详细调试信息日志 log等级：5
func Verbose(detail string) {
	noWrite = 1
	log(LevelVerbose, noWrite, fmt.Sprintf("[%s] [%s] %s", Cyan(getDate()), LightCyan("VERBOSE"), detail))
}

// Success 成功信息日志 log等级：2
func Success(detail string) {
	noWrite = 1
	log(LevelInfo, noWrite, fmt.Sprintf("[%s] [%s] %s", Cyan(getDate()), LightGreen("+"), detail))
}

// Failed 失败信息日志 log等级：2
func Failed(detail string) {
	noWrite = 1
	log(LevelInfo, noWrite, fmt.Sprintf("[%s] [%s] %s", Cyan(getDate()), LightRed("-"), detail))
}

// Common 普通信息日志 log等级：2
func Common(detail string) {
	noWrite = 1
	log(LevelInfo, noWrite, fmt.Sprintf("[%s] [%s] %s", Cyan(getDate()), LightGreen("*"), detail))
}

// Command 接收命令信息日志，该参数log不换行输出 log等级：2
func Command(detail string) {
	// 该命令接收参数信息，主要用于人工的参数输入，任何情况都不写入文件，在log函数中
	cmdInfo := fmt.Sprintf("[%s] [%s] %s", Cyan(getDate()), LightGreen("✎"), detail)
	fmt.Print(cmdInfo)
}

func LogError(errinfo interface{}) {
	if WaitTime == 0 {
		fmt.Println(fmt.Sprintf(" %v/%v %v", End, Num, errinfo))
	} else if (time.Now().Unix()-LogSucTime) > WaitTime && (time.Now().Unix()-LogErrTime) > WaitTime {
		fmt.Println(fmt.Sprintf(" %v/%v %v", End, Num, errinfo))
		LogErrTime = time.Now().Unix()
	}
}

func getTime() string {
	return time.Now().Format("15:04:05")
}

func getDate() string {
	return time.Now().Format("2006.1.2")
}

func DebugError(err error) bool {
	/* Processing error display */
	if err != nil {
		pc, _, line, _ := runtime.Caller(1)
		Debug(fmt.Sprintf("%s%s%s",
			White(runtime.FuncForPC(pc).Name()),
			LightWhite(fmt.Sprintf(" line:%d ", line)),
			White(err)))
		return true
	}
	return false
}

// Clean by https://github.com/acarl005/stripansi/blob/master/stripansi.go
func clean(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)
	return re.ReplaceAllString(str, "")
}

func writeLogFile(result string, filename string) {
	var text = []byte(result + "\n")
	fl, err := os.OpenFile(fmt.Sprintf("./result/%v", filename), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Open %s error, %v\n", filename, err)
		return
	}
	_, err = fl.Write(text)
	fl.Close()
	if err != nil {
		fmt.Printf("Write %s error, %v ", filename, err)
	}
}
