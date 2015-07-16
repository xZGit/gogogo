package godis


import
(
     "sync"
	"menteslibres.net/gosexy/redis"
	"log"
)


var host = "127.0.0.1"
var port = uint(6379)

type RedisClient struct {
	conn  *redis.Client
	pubConn  *redis.Client
	subConn  *redis.Client
	mutex sync.Mutex
}


type Value struct {
	value interface{}
}


func NewRedisClient(host string) *RedisClient {
//	host = fmt.Sprintf("%s:6379", host)
	var conn, pubConn, subConn *redis.Client

	conn = redis.New()
	pubConn = redis.New()
	subConn = redis.New()
	err := conn.Connect(host,port)
	err = pubConn.Connect(host,port)
	err =subConn.Connect(host,port)
	if err != nil {
		log.Fatalf(" failed to connect: %s\n", err.Error())
		panic(err)
	}

	//	pubsub, _ := redis.Dial("tcp", host)
	client := & RedisClient{
		conn:conn,
		pubConn:pubConn,
		subConn:subConn,
		}
	return client
}




func (redisClient *RedisClient) getActiveSubPub() {

	var i []string
	log.Printf("GET hello\n")
	err := redisClient.conn.Command(&i, "PUBSUB", "CHANNELS")
	if err != nil {
		log.Println("err %v",err)
	}
    log.Println("sssh %v",i)

}

