package godis


import (
	"sync"
	"log"
	uuid "github.com/nu7hatch/gouuid"
//	"strings"
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
		rec := make(chan []string)
		go s.redisClient.subConn.Subscribe(rec, name)
		i:=0
		var ls []string
		c <- 1
		for {
			ls = <-rec
			if ls[0]=="message" && len(ls)>2 {
				i=i+1
				go s.ProcessFunc(handlerFunc, ls[2], i)
			}
		}
	}()
}

var d=0
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
    d=d+1
//
// 	go s.redisClient.getActiveSubPub()
    h,err:= s.redisClient.pubConn.Publish(ev.MId, msg)
	if err != nil {
		log.Printf("err: %v\n", h)
	}else {

	 log.Printf("re i: %v\n", d)
	}
}