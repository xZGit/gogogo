package godis


import (
	"sync"
	"log"
	uuid "github.com/nu7hatch/gouuid"
	"strings"
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



func (s *Server) RegisterTask(name string, handlerFunc *func(args []string) (interface{}, error), c chan int) {
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



func (s *Server) ProcessFunc(handlerFunc *func(args []string) (interface{}, error), ls string) {


	l := len(ls)
	data := ls[1:l-1]
	strArray := strings.Fields(data)
	sl := len(strArray)
	var arg []string
	arg =strArray[2:sl]
	arg[0]=arg[0][1:]
	last := len(arg)-1
	arg[last]=arg[last][:len(arg[last])-1]

	log.Printf("process func: %v\n", arg)

	v, err := (*handlerFunc)(arg)
	if err != nil {
		panic(err)
	}


 	go s.redisClient.getActiveSubPub()
	log.Printf("result: %v\n", strArray[0])
	s.redisClient.pubConn.Publish(strArray[0],v)
}