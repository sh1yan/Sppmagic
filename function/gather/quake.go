package gather

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"main/function"
	log "main/function/logger"
	"main/function/subassembly"
	"net/http"
	"strconv"
	"time"
)

// GetQuakeSocks5Data 用于获取quake里的socks5数据
func GetQuakeSocks5Data() {

	GETRES := []string{}
	SearchSyntax := function.Conf.Section("quake").Key("rule").Value()
	quake_timeout := strconv.Itoa(120) // 因为作者是终身会员，这里为防止抓取1万条数据时卡住而导致失败
	postdata := quake_api_search(SearchSyntax, 0, function.Quake_size, function.Quake_permissions)
	res := postQuakeHttp(postdata, function.Quake_token, quake_timeout)
	log.Debug(fmt.Sprint(res.Value()))
	for i := range res.Get("data").Array() {
		log.Verbose(fmt.Sprint(res.Get("data").Array()[i].Get("ip"), ":", res.Get("data").Array()[i].Get("port")))
		ip_port := fmt.Sprintf("%v:%v", res.Get("data").Array()[i].Get("ip"), res.Get("data").Array()[i].Get("port"))
		GETRES = append(GETRES, ip_port)
	}
	function.Address = GETRES
	subassembly.MakeApiResultDown("Quake", GETRES)
	log.Debug(fmt.Sprintf("获取成功，总查询数量：%v", len(function.Address)))
}

var (
	Exclude_Senior_Member      = []string{"hostname", "transport", "asn", "org", "service.name", "location.country_cn", "location.province_cn", "location.city_cn", "service.http.host", "time", "service.http.title", "service.response", "service.cert", "components.product_catalog", "components.product_type", "components.product_level", "components.product_vendor", "location.country_en", "location.province_en", "location.city_en", "location.district_en", "location.district_cn", "location.isp", "service.http.body", "components.product_name_cn", "components.version", "service.http.infomation.mail", "service.http.favicon.hash", "service.http.favicon.data", "domain", "service.http.status_code"}
	Exclude_Registered_Members = []string{"hostname", "transport", "asn", "org", "service.name", "location.country_cn", "location.province_cn", "location.city_cn", "service.http.host", "service.http.title", "service.http.server"}
)

// QuakeData 该参数需要保留，用于post数据包查询时使用
type quakeData struct {
	Query   string   `json:"query"`
	Start   int      `json:"start"`
	Size    int      `json:"size"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// quake_api_search 用于填充JSON格式的API数据
func quake_api_search(keyword string, start int, size int, grade int) quakeData {
	// grade 如果结果大于0，则是高级会员，如果等于0则是普通会员

	reqData := quakeData{}
	reqData.Query = keyword
	reqData.Start = start
	reqData.Size = size
	reqData.Include = []string{"ip", "port"}
	if grade > 0 {
		reqData.Exclude = Exclude_Senior_Member
	} else {
		reqData.Exclude = Exclude_Registered_Members
	}

	return reqData
}

// quakehttp 核心代码块，用于API请求数据
func postQuakeHttp(postdata quakeData, key string, timeout string) gjson.Result {
	var itime, err = strconv.Atoi(timeout)
	if err != nil {
		log.Error(fmt.Sprintf("设置 Quake 超时参数错误: %v", err))
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Timeout:   time.Duration(itime) * time.Second,
		Transport: transport,
	}
	url := "https://quake.360.cn/api/v3/search/quake_service"
	payload, _ := json.Marshal(postdata)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Error(fmt.Sprint(err))
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Set("X-QuakeToken", key)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(fmt.Sprint(err))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	jsoninfostr := string(body)
	log.Debug(jsoninfostr)
	res := gjson.Parse(jsoninfostr)
	return res

}

func GetQuakeAccountInfo() {

	quake_token := function.Quake_token

	var API_Remain = "0" // 用于存放实时的API剩余次数的参数

	// 创建一个自定义请求对象
	req, err := http.NewRequest("GET", "https://quake.360.net/api/v3/user/info", nil)
	if err != nil {
		log.Error(fmt.Sprintf("创建GET形式请求请求失败:%v", err))
		return
	}
	// 添加自定义头部信息
	req.Header.Add("X-QuakeToken", quake_token)
	// 创建HTTP客户端并发送请求
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Error(fmt.Sprintf("进行GET数据请求失败:%v", err))
		return
	}
	defer response.Body.Close()
	// 读取响应的内容
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(fmt.Sprintf("进行读取API响应信息失败:%v", err))
		return
	}

	// 打印响应内容
	log.Debug("响应内容:")
	log.Debug(string(body))
	jsoninfostr := string(body)

	result := gjson.Parse(jsoninfostr)
	f_query_api_count := result.Get("data").Get("f_query_api_count").Value()
	if f_query_api_count != nil {
		log.Debug(fmt.Sprintf("当前API剩余次数_注册用户：%v", f_query_api_count))
		API_Remain = fmt.Sprint(f_query_api_count)
	} else {
		free_query_api_count := result.Get("data").Get("free_query_api_count").Value()
		log.Debug(fmt.Sprintf("当前API剩余次数_会员用户：%v", free_query_api_count))
		API_Remain = fmt.Sprint(free_query_api_count)
	}

	if API_Remain == "0" {
		log.Info("当前免费API次数已经用完了，本次将会使用月度积分或长效积分进行数据获取")
		log.Info(fmt.Sprintf("当前配置文件中设置的数据范围最大值为：%v", strconv.Itoa(function.Quake_size)))
		log.Info(fmt.Sprintf("当前月度剩余积分数：%v", result.Get("data").Get("month_remaining_credit").Value()))
		log.Info(fmt.Sprintf("当前长效积分剩余数：%v", result.Get("data").Get("persistent_credit").Value()))
	} else {
		log.Info(fmt.Sprintf("当前免费API次数还剩余：%v", API_Remain))
		if function.Quake_permissions == 0 {
			if function.Quake_size > 500 {
				log.Warning("您当前Quake账号为注册用户，默认数据长度为500条，超出部分将会以积分形式扣除")
			}
		}
		log.Info(fmt.Sprintf("当前配置文件中设置的数据范围最大值为：", strconv.Itoa(function.Quake_size)))
	}
}
