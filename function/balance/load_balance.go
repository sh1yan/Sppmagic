package balance

import "errors"

// RoundRobinBalance 数组循环取值结构体，应用于轮询策略，也就是每请求一次就换一次代理
type RoundRobinBalance struct {
	curIndex int      // 索引
	rss      []string // 订阅
}

// Set 将字符串数组放入到循环平衡结构体的订阅数组中
func (r *RoundRobinBalance) Set(s []string) error {
	if len(s) == 0 {
		return errors.New("input []string")
	}

	r.rss = s
	return nil
}

// next 通过取余数操作来循环遍历整个数组中的地址信息
func (r *RoundRobinBalance) next() string {
	if len(r.rss) == 0 {
		return ""
	}
	lens := len(r.rss)      // 将结构体订阅数组中的长度整理成数字
	if r.curIndex >= lens { // 如果结构体索引大于或等于数组长度，则将索引值替换为空
		r.curIndex = 0
	}

	curAddr := r.rss[r.curIndex] // 通过索引进行获取地址信息

	// 这行代码的目的是将 r.curIndex 增加1，并在达到数组或切片的长度（lens）时回绕到0。这是通过取余数操作来实现的，因此无论 r.curIndex 增加到多少，都会保持在 0 到 lens-1 的范围内。
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}

// Get 执行 r.next() 函数，来获得地址信息
func (r *RoundRobinBalance) Get() string {
	return r.next()
}
