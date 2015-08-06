package godis


import (
	"sync"
	"log"
	uuid "github.com/nu7hatch/gouuid"
	"time"
	"gopkg.in/redis.v3"
)


type Server struct {
	redisClient *RedisClient
	mutex       sync.Mutex
	uuid        string
}

func NewServer(host string) (*Server) {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return &Server{
		redisClient:NewRedisClient(host),
		uuid :id.String(),
	}
}



func (s *Server) RegisterTask(name string, handlerFunc HandleServerFunc, c chan int) {
	go func() {
		pubsub, err := s.redisClient.subConn.Subscribe(name)
		if err != nil {
			panic(err)
		}
		for {
			msgi, err := pubsub.ReceiveTimeout(100 * time.Millisecond)
			if err != nil {
				err := pubsub.Ping("")
				if err != nil {
					panic(err)
				}
				continue
			}
            i:=0
			switch msg := msgi.(type) {
				case *redis.Subscription:
				log.Println(msg.Kind, msg.Channel)
				case *redis.Message:
				go s.ProcessFunc(handlerFunc, msg.Payload, i)
				case *redis.Pong:
//				log.Println(msg)
				default:
				log.Println("unknown message: %#v", msgi)
			}
		}
	}()
}

var i=0
func (s *Server) ProcessFunc(handlerFunc HandleServerFunc, ls string,i int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()


	mySlice := []byte(ls)
	ev, err := unPackEventBytes(mySlice)
	if err != nil {
		panic(err)
	}
	v, err := (*handlerFunc)(ev.Args)
	var resp Resp
	if err != nil {
		resp = newResp(ev.MsgId, 1, err.Error(), nil)
	}else {
		resp = newResp(ev.MsgId, 0, "", v)
	}

	msg, err :=resp.packBytes()
    s.redisClient.pubConn.Publish(ev.MId, string(msg[:]))
	log.Println("i %v",i)
}