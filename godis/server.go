package godis


import(
	"sync"
	"log"
	"strings"
)


type Server struct {
	redisClient   *RedisClient
	mutex         sync.Mutex
	handleFuncs   []*taskHandler
}


// Task handler representation
type taskHandler struct {
	TaskName    string
	HandlerFunc *func(args []interface{}) (interface{}, error)
}



func NewServer(host string) (*Server) {
	return  &Server{
		redisClient:NewRedisClient(host),
		handleFuncs:make([]*taskHandler, 0),
	}
}



func (s *Server) RegisterTask(name string, handlerFunc *func(args []interface{}) (interface{}, error), c chan int) {
	go func() {
		for _, h := range s.handleFuncs {
			if h.TaskName == name {
				return
			}
		}
		rec := make(chan []string)
		go s.redisClient.conn.Subscribe(rec, name)

		var ls []string
		c <- 1
		for {
			ls = <-rec
			log.Printf("Consumer received: %v\n", strings.Join(ls, ", "))
		}
		s.handleFuncs = append(s.handleFuncs, &taskHandler{TaskName: name, HandlerFunc: handlerFunc})

		log.Printf("server registered handler for task %s", name)

	}()
}

