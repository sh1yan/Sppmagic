package subassembly

import (
	"flag"
	"golang.org/x/text/encoding/simplifiedchinese"
	"main/function"
	log "main/function/logger"
	"regexp"
)

func Flag() {
	flag.StringVar(&function.Txtfilepath, "f", "", "URL文件路径地址，请参照格式输入, -f D://proxy.txt")
	flag.StringVar(&function.Model, "m", "", "目前存在3种数据获取模式：native | quake | fofa , -m native")
	flag.IntVar(&function.LogLevel, "logl", 3, "设置日志输出等级，默认为3级，-logl 3")
	flag.StringVar(&function.OutputFileName, "o", "outcome.txt", "")
	flag.Parse()
}

// RemoveDuplicates 删除重复项
func RemoveDuplicates(arr []string) []string {
	encountered := map[string]bool{} // 用于记录已经遇到的元素
	result := []string{}             // 存储去重后的结果

	for _, value := range arr {
		if !encountered[value] {
			encountered[value] = true
			result = append(result, value)
		}
	}

	return result
}

// IsDomainNameWithPort 判断字符串是否是域名加端口形式
func IsDomainNameWithPort(str string) bool {
	domainPortPattern := `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}:\d+$` // 判断主域名加端口形式的正则
	match, _ := regexp.MatchString(domainPortPattern, str)
	if !match {
		subdomainsPortPattern := `^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z]{2,}:\d+$` // 判断子域名加端口形式的正则
		match1, _ := regexp.MatchString(subdomainsPortPattern, str)
		return match1
	}
	return match
}

// IsIPAddressWithPort 判断字符串是否是IP地址加端口号的格式
func IsIPAddressWithPort(str string) bool {
	ipPortPattern := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+$`
	match, _ := regexp.MatchString(ipPortPattern, str)
	return match
}

// IsCharacterEmpty 判断当前输入的文件是否为空
func IsCharacterEmpty(char string) bool {
	return len(char) != 0
}

// SliceDelete 接收传入的数组并删除制定的参数
func SliceDelete(tslice any, val any) any {

	if _, ok := tslice.([]string); ok { // 字符串数组 删除其中与指定值相等的元素，然后返回修改后的切片
		slice := tslice.([]string)
		for i := 0; i < len(slice); i++ {
			if slice[i] == val {
				slice = append(slice[:i], slice[i+1:]...)
				i--
			}
		}
		return slice

	} else if _, ok := tslice.([]int64); ok { // 数字数组 删除其中与指定值相等的元素，然后返回修改后的切片
		slice := tslice.([]int64)
		for i := 0; i < len(slice); i++ {
			if slice[i] == val {
				slice = append(slice[:i], slice[i+1:]...)
				i--
			}
		}
		return slice
	} else if _, ok := tslice.([][]string); ok { // 双重字符串数组 删除其中与指定值相等的元素，然后返回修改后的切片
		slice := tslice.([][]string)
		for i := 0; i < len(slice); i++ {
			if slice[i][0] == val {
				slice = append(slice[:i], slice[i+1:]...)
				i--
			}
		}
		return slice
	}
	log.Error("暂时不支持这种类型转换")
	panic("暂时不支持这种类型转换")

}

// SlicesFind 用于遍历参数2是否在参数1的数组中，若在则返回true
func SlicesFind(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func silcesIndex(slice []string, val string) int {
	for n, item := range slice {
		if item == val {
			return n
		}
	}
	return -1
}

func ConvertByte2String(byte []byte, charset string) string {

	var str string
	switch charset {
	case "GB18030":
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case "UTF8":
		fallthrough
	default:
		str = string(byte)
	}

	return str
}
