package godis


import
(
     "sync"
     "gopkg.in/redis.v3"
)


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

	conn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pubConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	subConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})


	//	pubsub, _ := redis.Dial("tcp", host)
	client := & RedisClient{
		conn:conn,
		pubConn:pubConn,
		subConn:subConn,
		}
	return client
}




func (redisClient *RedisClient) getActiveSubPub() {

//	var i []string
//	err := redisClient.conn.Command(&i, "PUBSUB", "CHANNELS")
//	if err != nil {
//		log.Println("err %v",err)
//	}

}

