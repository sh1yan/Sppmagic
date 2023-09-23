package inspect

import (
	"crypto/tls"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/valyala/fastjson"
	"io"
	"main/function"
	log "main/function/logger"
	"main/function/subassembly"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// CheckAlive 持续对socks5代理地址进行存活性探测，最后再放入循环结构体中
func CheckAlive() {
	lock := &sync.Mutex{} // 生成一个互斥锁，用来保护并发时产生的各种问题
	var fastj fastjson.Parser
	for {
		if len(function.Address) == 0 { // 这里因为是采用了并发模式，可能代理池还没有接收到到参数,所以挂在等待地址里存在数据
			time.Sleep(1 * time.Second)
		}
		// log.Debug(fmt.Sprint("正在过滤代理中"))

		// 目前ini配置中的thread线程是1000
		// 该并发池主要用于批量判断传入的socks5代理是否为存活的，同时获取代理的延迟和物理位置
		// 原理就是创建一个HTTP客户端的代理，然后使用该代理去请求百度open来获得物理位置，若存活则正常存放，若失败得从历史池中取出
		p, _ := ants.NewPoolWithFunc(function.Thread, func(i interface{}) { // 使用ants生成一个高效的 goroutine 池

			// 以下程序均是在代理池中运行

			socks5 := i.(string)

			// 当前函数结束时，该函数会判断程序结束时是否因为panic goroutine的程序错误导致的终端，若是的话则删除当前这个无效的socks5代理地址
			defer func() {
				if err := recover(); err != nil { // 判断是否因为panic的程序错误
					if subassembly.SlicesFind(function.Alive_address, socks5) == true { // 若socks5参数在存活的地址栏中则进入下列代码块
						function.Alive_address = subassembly.SliceDelete(function.Alive_address, socks5).([]string)             // 在存活列表中删除无用的socks5地址
						function.Alive_address_time = subassembly.SliceDelete(function.Alive_address_time, socks5).([][]string) // 删除多重数组中信息
					}
					function.Wg.Done() // 计数器减一
					return
				}
			}()
			socksProxy := "socks5://" + socks5
			proxy := func(_ *http.Request) (*url.URL, error) {
				return url.Parse(socksProxy) // 该代理函数通常用于将HTTP请求路由到代理服务器，以便在客户端和目标服务器之间建立代理连接。
			}
			tr := &http.Transport{ // 创建一个自定义的 HTTP 客户端
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Proxy:           proxy,
			}

			// 获取当前代理IP地址的信息
			urlProxy := fmt.Sprintf("https://opendata.baidu.com/api.php?query=%s&co=&resource_id=6006", strings.Split(socks5, ":")[0])
			start := time.Now().UnixNano() // 创建一个开始的时间戳
			req, _ := http.NewRequest("GET", urlProxy, nil)
			req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
			req.Header.Add("Connection", "close")
			resp, _ := (&http.Client{Transport: tr, Timeout: 8 * time.Second}).Do(req)
			stop := time.Now().UnixNano() - start // 创建一个开始的时间戳
			if resp.StatusCode == 200 {           // 判断是否正常获取到返回信息
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				v, _ := fastj.Parse(string(body)) // 以JSON格式解析返回的信息内容

				// 获取socks5的IP地址归属地址
				location := subassembly.ConvertByte2String(v.GetStringBytes("data", "0", "location"), "GB18030")
				if location != "" { // 判断socks5代理地址的物理位置是否为空
					if subassembly.SlicesFind(function.Alive_address, socks5) == false { // 判断当前代理地址是否不在socks5存活代理池中，若不在，则加入代理池中
						sliceAlive := []string{socks5, strconv.FormatInt(stop, 10), location} // 放入数组中三个值：代理IP、响应时间、物理位置
						lock.Lock()                                                           // 以下操作加锁
						function.Alive_address = append(function.Alive_address, socks5)       // 往存活socks5代理池中存放socks5代理地址
						log.Debug(fmt.Sprintf("当前socks5验活过程中，该地址为存活地址：%v", socks5))
						function.Alive_address_time = append(function.Alive_address_time, sliceAlive) // 往[][]string 中存放一个数组
						lock.Unlock()                                                                 // 以上操作解锁
						log.LogWriteInfo(socks5)                                                      // 将存活的socks5代理地址放入报本地保存记录中
					}
				}
			}
			function.Wg.Done() // 并发线程减一
		})
		defer p.Release()                    // 函数结束时，Release 关闭该池并释放工人队列。
		for _, i := range function.Address { // 获取从fofa中搜集到的socks5代理地址
			function.Wg.Add(1) // 工作池进程加一
			_ = p.Invoke(i)    // 向任务池提交任务，提交socks5地址
		}
		function.Wg.Wait() // 等待工作池全部运行完毕
		log.Debug("程序当前所有导入到验活功能里socks5地址以全部验活完成")
		function.B.Set(function.Alive_address) // 将存活socks5代理地址都放入循环体中，进行循环使用
		log.Debug(fmt.Sprintf("截止到本轮验活中，共有 %v 条存活地址", len(function.Alive_address)))
		// fmt.Println("[*]", "过滤完成，一共有", len(alive_address), "/", len(alive_address_time), "个代理可用")
		time.Sleep(10 * time.Second) // 挂载十秒钟
	}
}
