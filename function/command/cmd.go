package command

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"main/function"
	log "main/function/logger"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
)

// cmd 存放输入的命令
type cmd struct {
	command string
}

var (
	c = cmd{} // 实例化一个cmd结构体
)

// Show 显示当前代理地址延迟情况或者显示当前存活代理池里的IP地址信息
func (c cmd) Show(args []string) {
	if len(args) > 1 {
		if args[1] == "ip" {
			for _, i := range function.Alive_address_time {
				if i[0] == function.Tmp_addr {
					defer func() {
						if err := recover(); err != nil {
							log.Info(fmt.Sprint("当前使用的IP地址是 ", function.Tmp_addr, " ", i[2], " 延迟", "错误"))
							return
						}
					}()
					socksProxy := "socks5://" + function.Tmp_addr
					proxy := func(_ *http.Request) (*url.URL, error) {
						return url.Parse(socksProxy)
					}
					tr := &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
						Proxy:           proxy,
					}
					start := time.Now()
					url := "https://opendata.baidu.com/api.php" // {"status":1,"msg":"\u53c2\u6570\u9519\u8bef","data":[]}
					req, _ := http.NewRequest("GET", url, nil)
					req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
					(&http.Client{Transport: tr, Timeout: 8 * time.Second}).Do(req)
					end := time.Since(start)
					log.Info(fmt.Sprint("当前使用的IP地址是", function.Tmp_addr, i[2], "延迟", end)) // 延迟测试
					return
				}
			}
			fmt.Println("请先使用代理后在执行这条命令")
			return
		} else if args[1] == "all" { // 显示当前所有socks5代理资源池里的IP地址
			log.Info("正在输出全部IP地址")
			log.Info("---------------------------------")
			for _, i := range function.Alive_address_time {
				log.Info(fmt.Sprint(i[0], "  ", i[1], "  ", i[2]))
			}
			log.Info("---------------------------------")
			log.Info(fmt.Sprint("一共有 ", len(function.Alive_address_time), " 代理"))
			return
		} else {
			log.Error(fmt.Sprint(args[0], " ", args[1], ":参数错误！"))
		}
	}

}

// Use 设置固定代理地址或使用随机代理地址
func (c cmd) Use(args []string) {
	if len(args) > 1 {
		if strings.Contains(args[1], ":") {
			function.SetProxy = true
			function.UseProxy = args[1]
			log.Success(fmt.Sprint("设置代理 ", function.UseProxy, " 成功"))
			return
		} else if args[1] == "random" {
			function.SetProxy = false
			function.UseProxy = ""
			log.Success("设置random成功")
			return
		}
	}
	log.Error(fmt.Sprint(args[0], ":参数错误！"))
}

// Exit 退出命令
func (c cmd) Exit(args []string) {
	os.Exit(0)
}

// Command 接收命令并分析命令再执行命令
func Command() {

	for {
		if len(function.Alive_address) > 0 { // 判断是否存在存活地址
			c.command = "" // 初始化命令
			defer func() { // 若程序崩溃，则输出错误，并继续执行 Command() 函数
				if err := recover(); err != "" {
					// log.LogError(err) // debug的时候使用
					log.Error(fmt.Sprint(c.command+" : ", "指令错误"))
					Command()
				}
			}()
			log.Command("-> ")                       // 提示符
			reader := bufio.NewReader(os.Stdin)      // 接收外部参数传入
			c.command, _ = reader.ReadString('\n')   // 判断是否接收到换行符了
			c.command = strings.TrimSpace(c.command) // 去掉前后无用空白
			if c.command == "" {                     // 若当前命令为空，则跳出当次循环
				continue
			}
			funcs := reflect.ValueOf(&c)                           // 传入的参数为结构体的函数名
			comm := strings.ToUpper(c.command[:1]) + c.command[1:] // 小写变大写 组合成字符串命令
			var args []reflect.Value                               // 定义一个  Value 是 Go 值的反射接口。
			if len(strings.Split(c.command, " ")) > 1 {            // 判断是否存在输入命令
				comm = strings.Split(comm, " ")[0]                                     // 获取到参数的名称
				args = []reflect.Value{reflect.ValueOf(strings.Split(c.command, " "))} // 获取到参数的值
			} else {
				args = []reflect.Value{reflect.ValueOf([]string{c.command})} // 默认获取参数的名称
			}
			funcs.MethodByName(comm).Call(args) // show ip    or    show all
		}

	}
}
