package redisUtil

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/core/redis"

	redisTrue "github.com/go-redis/redis/v8"
	"time"
)

func Set(c *gin.Context, key string, value string, expiration time.Duration) error {
	err := redis.Instance().Set(c, key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
	return err
}

// 通过key获取值
// err == redis.Nil, "key不存在"
// err != nil,  "错误"
// val == "", "值是空字符串"
func Get(c *gin.Context, key string) (string, error) {
	value, err := redis.Instance().Get(c, key).Result()
	if err != nil && value != "" {
		panic(err)
	}
	return value, err
}

// 删除key
func Del(c *gin.Context, keys ...string) error {
	err := redis.Instance().Del(c, keys...).Err()
	if err != nil {
		panic(err)
	}
	return err
}

// 指定key的值+1
func Incr(c *gin.Context, key string) error {
	err := redis.Instance().Incr(c, key).Err()
	if err != nil {
		panic(err)
	}
	return err
}

// 指定key的值+1
func Decr(c *gin.Context, key string) error {
	err := redis.Instance().Decr(c, key).Err()
	if err != nil {
		panic(err)
	}
	return err
}

// 批量设置，没有过期时间
func MSet(c *gin.Context, values ...interface{}) error {
	err := redis.Instance().MSet(c, values...).Err()
	return err
}

// 批量设置取数据
// 示例：values, err := MGet(key1, key2)
// for i, _ := range values {
// fmt.Println(values[i])
// }
func MGet(c *gin.Context, keys ...string) ([]interface{}, error) {
	values, err := redis.Instance().MGet(c, keys...).Result()
	return values, err
}

// 执行命令
// 返回结果
// s, err := cmd.Text()
// flag, err := cmd.Bool()
// num, err := cmd.Int()
// num, err := cmd.Int64()
// num, err := cmd.Uint64()
// num, err := cmd.Float32()
// num, err := cmd.Float64()
// ss, err := cmd.StringSlice()
// ns, err := cmd.Int64Slice()
// ns, err := cmd.Uint64Slice()
// fs, err := cmd.Float32Slice()
// fs, err := cmd.Float64Slice()
// bs, err := cmd.BoolSlice()
func Do(c *gin.Context, args ...interface{}) *redisTrue.Cmd {
	cmd := redis.Instance().Do(c, args)
	return cmd
}

// 清空缓存
func FlushDB(c *gin.Context) error {
	err := redis.Instance().FlushDB(c).Err()
	return err
}

// 发布
// 示例Publish("mychannel1", "payload").Err()
func Publish(c *gin.Context, channel string, msg string) error {
	err := redis.Instance().Publish(c, channel, msg).Err()
	return err
}

// 订阅
func Subscribe(c *gin.Context, channel string, subscribe func(msg *redisTrue.Message, err error)) {
	pubsub := redis.Instance().Subscribe(c, channel)
	// 使用完毕，记得关闭
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(c)
		subscribe(msg, err)
	}
}

// 列表的头部（左边）,尾部（右边）
// 列表左边插入
func LPust(c *gin.Context, channel string, values ...interface{}) error {
	return redis.Instance().LPush(c, channel, values...).Err()
}

// 列表从左边开始取出start至stop位置的数据
func LRange(c *gin.Context, key string, start, stop int64) error {
	return redis.Instance().LRange(c, key, start, stop).Err()
}

// 列表左边取出
func LPop(c *gin.Context, key string) *redisTrue.StringCmd {
	return redis.Instance().LPop(c, key)
}

// 列表右边插入
func RPust(c *gin.Context, channel string, values ...interface{}) error {
	return redis.Instance().RPush(c, channel, values...).Err()
}

// 列表右边取出
func RPop(c *gin.Context, key string) error {
	return redis.Instance().RPop(c, key).Err()
}

// 列表哈希插入
func HSet(c *gin.Context, key string, field string, values ...interface{}) error {
	return redis.Instance().HSet(c, key, field, values).Err()
}

// 列表哈希取出
func HGet(c *gin.Context, key, field string) *redisTrue.StringCmd {
	return redis.Instance().HGet(c, key, field)
}

// 列表哈希批量插入
func HMSet(c *gin.Context, key string, values map[string]interface{}) error {
	return redis.Instance().HMSet(c, key, values).Err()
}

// 列表哈希批量取出
func HMGet(c *gin.Context, key string, fields ...string) []interface{} {
	return redis.Instance().HMGet(c, key, fields...).Val()
}

// 列表无序集合插入
func SAdd(c *gin.Context, key string, members ...interface{}) error {
	return redis.Instance().SAdd(c, key, members...).Err()
}

// 列表无序集合，返回所有元素
func SMembers(c *gin.Context, key string) []string {
	return redis.Instance().SMembers(c, key).Val()
}

// 列表无序集合，检查元素是否存在
func SIsMember(c *gin.Context, key string, member interface{}) bool {
	b, err := redis.Instance().SIsMember(c, key, member).Result()
	if err != nil {
		panic(err)
	}
	return b
}
