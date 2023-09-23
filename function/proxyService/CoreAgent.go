package proxyService

import (
	"fmt"
	"io"
	"main/function"
	log "main/function/logger"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Process 负责把代理服务端获取的数据互相传递给客户端地址
func Process(client net.Conn) {
	defer func() { // 若程序报错，则继续执行当前函数
		if err := recover(); err != nil {
			Process(client)
		}
	}()
	if len(function.Alive_address) == 0 { // 若存活socks5代理池没有存活地址了，则等待1秒再继续
		time.Sleep(1 * time.Second)
		Process(client)
	}
	defer client.Close()               // 程序结束时，关停监听客户端
	addr := getproxy(function.Tmp_try) // 用于获取当前需要使用的socks5代理地址
	function.Tmp_addr = addr           // 将当前需要使用的地址传递给 tmp_addr 参数
	log.Info(fmt.Sprint("当前使用用的ip是", addr))
	cc, err := net.DialTimeout("tcp", addr, 5*time.Second) // 建立一个5秒超时的tcp客户端
	if err != nil {
		log.Error(fmt.Sprint("connect error:", err))
		function.Tmp_try = 1 + function.Tmp_try // 则进行使用测速第二名的地址
		Process(client)
	}
	function.Tmp_try = 0   // 重置使用测速第一的代理地址
	defer cc.Close()       // 程序结束时，关闭tcp客户端
	go io.Copy(cc, client) // 代理地址端数据传递给客户端地址数据
	io.Copy(client, cc)    // 客户端地址数据传递给代理地址端数据
}

// getproxy 根据配置及人工输入，分配是使用固定代理IP，还是轮询策略代理IP，还是速度最快的代理地址IP
func getproxy(arg int) string {
	polling := function.Conf.Section("rule").Key("polling").Value() // 获取配置文件中轮询策略设置情况
	polling = strings.ToLower(polling)                              // 大小换成小写
	if function.SetProxy == true {                                  // 判断是否使用了固定代理地址，若使用了则返回代理地址
		return function.UseProxy
	} else {
		if polling == "true" {
			return function.B.Get() // 轮询策略，也就是每请求一次就换一次代理，默认是开启状态
		} else { // 如果不使用固定IP，也没有使用轮询，将使用下面步骤，也就是默认使用速度最快的
			var times []int64
			for _, i := range function.Alive_address_time {
				t, err := strconv.ParseInt(i[1], 10, 64)
				if err != nil {
					log.Error(fmt.Sprint("代理 ", function.Alive_address_time[0], " 时间转换失败"))
					return ""
				}
				times = append(times, t)                // 存放所有测试时间
				sort.Slice(times, func(i, j int) bool { // 时间测试排序
					return times[i] < times[j]
				})
			}
			one := times[0]
			two := times[1]
			three := times[2]
			for _, i := range function.Alive_address_time { // 筛选测速结果中比较快的前三名的socks5代理地址，并返回
				t, _ := strconv.ParseInt(i[1], 10, 64)
				switch arg {
				case 0:
					if t == one {
						return i[0]
					}
				case 1:
					if t == two {
						return i[0]
					}
				case 2:
					if t == three {
						return i[0]
					}
				default:
					if t == one {
						return i[0]
					}
				}

			}
			log.Failed("当前无代理")
			return ""
		}
	}
}
