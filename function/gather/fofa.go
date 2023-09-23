package gather

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/valyala/fastjson"
	"io"
	"main/function"
	log "main/function/logger"
	"main/function/subassembly"
	"net/http"
	"time"
)

// GetFofaSocks5Data 通过fofa获取到现有的socks5不需要认证代理的国内数据，fofa是每60条获取下最新的
func GetFofaSocks5Data() error {

	function.Tmp = 1
	for { // 外循环使用了并发标识进行并发，下列又实用了tmp进行结果延缓等待，避免应并发过快，导致程序结束
		if function.Tmp > 1 { // 第一次运行不需要挂在60秒，第二次for循环时将挂在1分钟后再执行
			time.Sleep(60 * time.Second)
		}
		function.Tmp = function.Tmp + 1
		req, err := http.NewRequest("GET", "https://fofa.info/api/v1/search/all", nil) // 获取fofa API 接口数据
		if err != nil {
			log.LogError(err)
			return err
		}

		// 下行代码创建了一个自定义的http.Transport对象，其中包含一个TLS客户端配置，该配置允许忽略服务器证书验证。
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		// 参数介绍 rule := `protocol=="socks5" && "Version:5 Method:No Authentication(0x00)" && country="CN"`
		rule := function.Conf.Section("fofa").Key("rule").Value() // 获取下fofa的socks5不需账号验证的搜索语法
		rule = base64.StdEncoding.EncodeToString([]byte(rule))    // 对搜索语法进行base64位编码
		r := req.URL.Query()
		r.Add("email", function.Fofa_email)
		r.Add("key", function.Fofa_key)
		r.Add("qbase64", rule)
		r.Add("size", "2000")         // 获取每页查询数量，这个参数感觉可以修改到配置文件里，进行手工修改
		req.URL.RawQuery = r.Encode() // 组合上述添加的URL参数
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
		resp, err := (&http.Client{Transport: tr}).Do(req) // 获取到fofa上前2000个socks5代理地址信息
		if err != nil {
			log.LogError(err)
			return err
		}
		defer resp.Body.Close()                // 在函数结束后自动关闭响应体
		body, _ := io.ReadAll(resp.Body)       // 获取到比特形式的界面返回结果
		var p fastjson.Parser                  // 定义一个 json 格式的参数
		v, _ := p.Parse(string(body))          // 分析界面返回的请求数据
		if v.GetStringBytes("errmsg") != nil { // 判断是否获取失败，若获取失败，则返回报告信息
			log.Error(fmt.Sprint(string(body)))
			return err
		}
		var rst []string                          // 定义一个存放结果的字符串数组
		for _, i := range v.GetArray("results") { // 获取JSON中results对应的值
			ipaddr := string(i.GetStringBytes("1")) + ":" + string(i.GetStringBytes("2")) // 获取 IP:port 数据
			rst = append(rst, ipaddr)                                                     // 保存到字符串数据组
		}
		function.Address = rst
		subassembly.MakeApiResultDown("Fofa", rst)
		log.Debug(fmt.Sprint("获取成功，总查询数量：", len(function.Address)))
	}
}
