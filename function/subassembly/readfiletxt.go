package subassembly

import (
	"bufio"
	"fmt"
	log "main/function/logger"
	"os"
)

// FindTextIpPort 获取本地text中 ip:port 地址列表
func FindTextIpPort(filepath string) []string {
	filePath := filepath // 替换为实际的txt文件路径

	// 打开文件
	file, err := os.Open(filePath)
	log.Debug(fmt.Sprint("判断当前是否读取到数据验证：", file))
	if err != nil {
		log.DebugError(err)
		return []string{}
	}
	defer file.Close()

	ipports := []string{} // 用于存储IP:PORT地址的列表

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 提取IP:PORT地址
		ipport := examineIpPort(line)
		if ipport != "" {
			log.Verbose(fmt.Sprint("当前分析文本里读取到的值为：", ipport))
			ipports = append(ipports, ipport)
		}
	}

	if err := scanner.Err(); err != nil {
		log.DebugError(err)
		return []string{}
	}
	log.Debug(fmt.Sprint("当前全部读取并确认的文本里数据长度为：", len(ipports)))
	return ipports
}

// examineIpPort 分析字符串参数，若为IP:PORT格式的则，并生成url地址
func examineIpPort(line string) string {

	// 判断当前输入的ip地址或者域名地址是否包含 http 或者 https
	if IsIPAddressWithPort(line) || IsDomainNameWithPort(line) {
		return line
	}

	return ""
}
