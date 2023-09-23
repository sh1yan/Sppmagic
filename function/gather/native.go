package gather

import (
	"flag"
	"fmt"
	"main/function"
	log "main/function/logger"
	"main/function/subassembly"
	"os"
)

// GetNativeSocks5Data 获取本地代理地址文件数据
func GetNativeSocks5Data() {

	log.Debug(fmt.Sprint("当前传入文件路径为：", function.Txtfilepath)) // debug 查看下是否有路径参数传过来
	rst := subassembly.FindTextIpPort(function.Txtfilepath)
	if len(rst) == 0 {
		if len(function.Txtfilepath) == 0 {
			log.Error("未输入 -f 参数的值，请按照flag提示进行输入")
			flag.PrintDefaults()
			os.Exit(0)
		}
		log.Error(fmt.Sprintf("当前导入的路径文件内容为空，请在文本里填入socks5代理地址：%v", function.Txtfilepath))
		os.Exit(0)
	}
	function.Address = rst
}
