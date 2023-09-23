package function

import (
	"github.com/go-ini/ini"
	"main/function/balance"
	"sync"
)

// 当前版本信息
var Version = "1.0.1"

// logo
var Slogan = `

 ________  ________  ________  _____ ______   ________  ________  ___  ________     
|\   ____\|\   __  \|\   __  \|\   _ \  _   \|\   __  \|\   ____\|\  \|\   ____\    
\ \  \___|\ \  \|\  \ \  \|\  \ \  \\\__\ \  \ \  \|\  \ \  \___|\ \  \ \  \___|    
 \ \_____  \ \   ____\ \   ____\ \  \\|__| \  \ \   __  \ \  \  __\ \  \ \  \       
  \|____|\  \ \  \___|\ \  \___|\ \  \    \ \  \ \  \ \  \ \  \|\  \ \  \ \  \____  
    ____\_\  \ \__\    \ \__\    \ \__\    \ \__\ \__\ \__\ \_______\ \__\ \_______\
   |\_________\|__|     \|__|     \|__|     \|__|\|__|\|__|\|_______|\|__|\|_______|
   \|_________|                                                                     

			Sppmagic version: ` + Version + `

`

var (
	Txtfilepath string // 用于存放免费代理路径的参数
	Model       string // 用于选择数据获取模式的参数
)

var (
	Conf              *ini.File // 用于后期存在 config.ini 中的参数信息
	Fofa_email        string    // 用于存放fofa的邮箱信息
	Fofa_key          string    // 用于存放fofa的token信息
	Thread            int       // 这个是存活校验的线程参数
	Quake_token       string    // 用于存放quake的token
	Quake_size        int       // 用于设置每次quake查询的返回最大条数
	Quake_permissions int       // 用于判断当前用户是注册用户还是会员用户，做出一些提醒
)

var (
	Address            []string                       // 用于存在通过采集模块获取到的socks5的IP:PORT的数据
	Alive_address      []string                       // 用于存在存活的socks5的地址栏
	Alive_address_time [][]string                     // 用于存放这三个信息 代理IP、响应时间、物理位置
	B                  = &balance.RoundRobinBalance{} // 这是一个数组循环使用的结构体
	Wg                 sync.WaitGroup                 // 定义一个等待组
)

var (
	Tmp      int
	Tmp_try  int    // 该参数表示使用代理比较快的前三名的哪一个，参数固定为  0,1,2  默认初始化为 0 ，使用最快的呢个
	Tmp_addr string // 临时IP代理使用地址
)

var (
	SetProxy = false // 设置固定代理地址
	UseProxy = ""
)

var (
	LogLevel       int    // log等级,默认设置3级
	NoColor        bool   // 是否开启log输出非颜色版设置
	OutputFileName string // 用于设置log输出名称设置
	NoSave         bool   // not save file // logsync.go 中设置不进行日志写入的设置, 注：在常规的logger中并没有设置该参数
)
