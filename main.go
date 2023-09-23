package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"main/function"
	"main/function/command"
	"main/function/gather"
	"main/function/inspect"
	log "main/function/logger"
	"main/function/proxyService"
	"main/function/subassembly"
	"net"
	"os"
	"strconv"
)

// init 初始配置文件加载显示
func init() {
	fmt.Print(function.Slogan) // logo 输出
	fmt.Println("原作者：Ggasdfg321 By T00ls.Com")
	subassembly.Flag()

	tcfgs, err := ini.Load("config.ini") // 加载默认配置文件
	if err != nil {
		// 若不存在以上配置文件，则进行在当前文件夹下生成默认的config.ini配置文件  // 以下base64位字符就是配置文件信息
		tconf, _ := base64.StdEncoding.DecodeString("W2dsb2JhbF0KYmluZF9pcCA9IDEyNy4wLjAuMQpiaW5kX3BvcnQgPSAxMDgwCnRocmVhZCA9IDEwMDAKCltxdWFrZV0KcXVha2V0b2tlbiA9ICIiCnJ1bGUgPSBzZXJ2aWNlOiJzb2NrczUiIGFuZCByZXNwb25zZToiQWNjZXB0ZWQgQXV0aCBNZXRob2Q6IDB4MCIgYW5kIGNvdW50cnk6Q04Kc2l6ZSA9IDUwMAojIHBlcm1pc3Npb25zIOS4uui0puWPt+eahOadg+mZkO+8jOmcgOimgeaJi+W3peWhq+WGme+8jDDkuLrms6jlhoznlKjmiLfvvIwx5Li66auY57qn5Lya5ZGYL+e7iOi6q+S8muWRmApwZXJtaXNzaW9ucyA9IDAKCltmb2ZhXQplbWFpbCA9CmtleSA9IApydWxlID0gJ3Byb3RvY29sPT0ic29ja3M1IiAmJiAiVmVyc2lvbjo1IE1ldGhvZDpObyBBdXRoZW50aWNhdGlvbigweDAwKSIgJiYgY291bnRyeT0iQ04iJwoKW3J1bGVdCiMg5piv5ZCm5byA5ZCv6L2u6K+i562W55Wl77yM5Lmf5bCx5piv5q+P6K+35rGC5LiA5qyh5bCx5o2i5LiA5qyh5Luj55CG77yM5LiN5byA5ZCv55qE6K+d5bCx5piv5Zu65a6a6YCf5bqm5pyA5b+r55qE5Luj55CGaXAKcG9sbGluZz0gdHJ1ZQ==")
		ioutil.WriteFile("config.ini", tconf, 0666)
	}
	function.Conf = tcfgs // 把本地读取到的信息传递给提前定义的 conf 变量中
	if len(function.Model) != 0 {
		switch function.Model {
		case "native":
			log.Debug("当前 -m 输入的参数为：native")
			gather.GetNativeSocks5Data()
		case "quake":
			log.Debug("当前 -m 输入的参数为：quake")
			function.Quake_token = function.Conf.Section("quake").Key("quaketoken").Value()                         // 获取token
			function.Quake_size, _ = strconv.Atoi(function.Conf.Section("quake").Key("size").Value())               // 获取JSON数据查询最大值
			function.Quake_permissions, _ = strconv.Atoi(function.Conf.Section("quake").Key("permissions").Value()) // 获取当前用户的权限
			gather.GetQuakeAccountInfo()                                                                            // 获取当前quake账号的积分和API接口剩余次数信息
			gather.GetQuakeSocks5Data()                                                                             // 获取API查询获取到的信息
		case "fofa":
			log.Debug("当前 -m 输入的参数为：fofa")
			function.Fofa_email = function.Conf.Section("fofa").Key("email").Value() // 获取邮箱
			function.Fofa_key = function.Conf.Section("fofa").Key("key").Value()     // 获取秘钥
		default:
			log.Debug(fmt.Sprintf("当前 -m 输入的参数为：%v", function.Model))
			log.Error("输入的 -m 的参数错误，请重新输入")
			flag.PrintDefaults()
			os.Exit(0)
		}
	} else {
		log.Error("没有输入指定的参数，使用该工具必须输入 -m 的参数")
		flag.PrintDefaults()
		os.Exit(0)
	}

	function.Thread, _ = strconv.Atoi(function.Conf.Section("global").Key("thread").Value()) // 获取线程数
	function.Tmp_try = 0                                                                     // 初始化 tmp_try 参数，用于筛选使用测速最快的前三的哪个IP地址
}

// main 工具入口函数处
func main() {
	log.Info("正在获取socks5代理中")
	subassembly.MkdirResult() // 创建存放一些结果数据的目录

	if function.Model == "fofa" {
		go gather.GetFofaSocks5Data() // 通过fofa自动获取大量socks5代理地址
	}
	go inspect.CheckAlive() // 对大量socks5地址进行存活性探测
	go command.Command()    // 通过命令设置，使用固定代理还是使用随机代理地址

	// 获取配置文件中的本地监听地址 127.0.0.1 监听 1080
	add := function.Conf.Section("global").Key("bind_ip").Value() + ":" + function.Conf.Section("global").Key("bind_port").Value()
	server, err := net.Listen("tcp", add)
	if err != nil {
		log.LogError(err)
		return
	}
	for { // 持续监听当前端口
		client, err := server.Accept() // 接受等待并向监听器返回下一个连接
		if err != nil {
			log.Error(fmt.Sprintf("Accept failed : %v", err))
			continue
		}
		go proxyService.Process(client) // 循环监听数据之间的传递
	}
}
