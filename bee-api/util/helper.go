package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"time"
)

func ReGroup[T1 any, T2 comparable](list []T1, fn func(a T1) T2) map[T2][]T1 {
	ret := make(map[T2][]T1)
	for _, t1 := range list {
		key := fn(t1)
		if _, ok := ret[key]; !ok {
			ret[key] = make([]T1, 0, 2)
		}
		ret[key] = append(ret[key], t1)
	}
	return ret
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func IsEmptyArrayOrSlice(i interface{}) bool {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		// 如果不是切片或数组类型，直接返回false，因为不可能是空数组或切片
		return false
	}

	// 使用Len方法获取切片或数组的长度
	return v.Len() == 0
}

func GetFreePort() (int, error) {
	// 监听任意IP地址的一个空闲端口
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	// 获取监听的端口
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// / 生成随机订单号
func GenerateOrderNo(prefix string, uid int64) string {
	// 获取当前时间戳，格式化为字符串
	timestamp := time.Now().Format("20060102150405")

	// 生成随机数
	randomNum := rand.Intn(9000) + 1000 // 生成1000到9999之间的随机数

	// 拼接订单号
	orderNo := fmt.Sprintf("%s%s%d%d", prefix, timestamp, uid%10000, randomNum)
	return orderNo
}
