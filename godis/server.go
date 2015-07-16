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
	handleFuncs []*taskHandler
	uuid        string
}


// Task handler representation
type taskHandler struct {
	TaskName    string
	HandlerFunc *func(args []interface{}) (interface{}, error)
}



func NewServer(host string) (*Server) {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return &Server{
		redisClient:NewRedisClient(host),
		handleFuncs:make([]*taskHandler, 0),
		uuid :id.String(),
	}
}



func (s *Server) RegisterTask(name string, handlerFunc HandleServerFunc, c chan int) {
	go func() {
		for _, h := range s.handleFuncs {
			if h.TaskName == name {
				return
			}
		}
		rec := make(chan []string)
		go s.redisClient.subConn.Subscribe(rec, name)

		var ls []string
		c <- 1
		for {
			ls = <-rec
			if ls[0]=="message" && len(ls)>2 {
				go s.ProcessFunc(handlerFunc, ls[2])
			}

		}

	}()
}



func (s *Server) ProcessFunc(handlerFunc HandleServerFunc, ls string) {

	mySlice := []byte(ls)
	log.Printf("mySlice: %v\n", mySlice)
	ev, err := unPackEventBytes(mySlice)
	if err != nil {
		panic(err)
	}

	v, err := (*handlerFunc)(ev.Args)
	var resp Resp
	if err != nil {
		resp = newResp(1,err.Error(),nil)
	}else {
		resp = newResp(0,"",v)
	}
	msg, err :=resp.packBytes()

//
// 	go s.redisClient.getActiveSubPub()
	log.Printf("result: %v\n", ev.MsgId)
	s.redisClient.pubConn.Publish(ev.MsgId,msg)
}