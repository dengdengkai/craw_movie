// movie_info
package models

import (
	"github.com/astaxie/goredis"
)

const (
	URL_QUEUE     = "url_queue"
	URL_VISIT_SET = "url_visit_set"
)

//定义一个连接实例
var (
	client goredis.Client
)

//redis高速内存key-value数据库连接地址
func ConnectRedis(addr string) {
	client.Addr = addr
}

//放入队列以key-value形式
func PutinQueue(url string) {
	client.Lpush(URL_QUEUE, []byte(url))
}

//获取消息
func PopfromQueue() string {
	res, err := client.Rpop(URL_QUEUE)
	if err != nil {
		panic(err)
	}
	return string(res)
}

func AddToSet(url string) {
	//添加地址，将url转换为数组类型
	client.Sadd(URL_VISIT_SET, []byte(url))
}

func GetQueueLength() int {

	length, err := client.Llen(URL_QUEUE)
	if err != nil {
		return 0
	}
	return length
}

func IsVisit(url string) bool {

	bIsVisit, err := client.Sismember(URL_VISIT_SET, []byte(url))
	if err != nil {
		return false
		//没访问过
	}
	return bIsVisit
}
