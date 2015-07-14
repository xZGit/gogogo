package godis


import
(
    "github.com/garyburd/redigo/redis"
     "sync"
	"fmt"
)




type RedisClient struct {
	conn  redis.Conn
	pubSubConn redis.PubSubConn
	mutex sync.Mutex
}


func NewRedisClient(host string) *RedisClient {
	host = fmt.Sprintf("%s:6379", host)
	conn, _ := redis.Dial("tcp", host)
	pubsub, _ := redis.Dial("tcp", host)
	client := & RedisClient{
		conn:conn,
		pubSubConn:redis.PubSubConn{pubsub}}
	return client
}


