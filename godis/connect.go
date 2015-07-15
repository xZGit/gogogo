package godis


import
(
//    "github.com/garyburd/redigo/redis"
     "sync"
//	"fmt"
	"menteslibres.net/gosexy/redis"
	"log"
)


var host = "127.0.0.1"
var port = uint(6379)

type RedisClient struct {
	conn  *redis.Client
//	pubSubConn redis.PubSubConn
	mutex sync.Mutex
}


type Value struct {
	value interface{}
}


func NewRedisClient(host string) *RedisClient {
//	host = fmt.Sprintf("%s:6379", host)
	var conn *redis.Client

	conn = redis.New()
	err := conn.Connect(host,port)

	if err != nil {
		log.Fatalf(" failed to connect: %s\n", err.Error())
		panic(err)
	}

	//	pubsub, _ := redis.Dial("tcp", host)
	client := & RedisClient{
		conn:conn,
		}
	return client
}




func (redisClient *RedisClient) getActiveSubPub() {

//	v := bytes.Index(n[0], []byte{0})
//	fmt.Printf("n %v",v)
//	var buf bytes.Buffer
//	enc := gob.NewEncoder(&buf)
//	err = enc.Encode(n[0])
//	if err != nil {
//
//	}
//
//	g := CToGoString(buf.Bytes()[:])
//	fmt.Printf("last %s",g)

}

