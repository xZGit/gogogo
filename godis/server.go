package godis


import(
	"sync"
	"log"
	"github.com/garyburd/redigo/redis"
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



func (s *Server) RegisterTask(name string, handlerFunc *func(args []interface{}) (interface{}, error)) {
	go func() {
		for _, h := range s.handleFuncs {
			if h.TaskName == name {
				return
			}
		}
		psc := s.redisClient.pubSubConn
		psc.Subscribe(name)
		for {
			switch v := psc.Receive().(type) {
				case redis.Message:
				log.Printf(" message: %s\n", v)
				case redis.Subscription:
				log.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
				break
				case error:
				return
			}
		}
		s.handleFuncs = append(s.handleFuncs, &taskHandler{TaskName: name, HandlerFunc: handlerFunc})

		log.Printf("server registered handler for task %s", name)

	}()
}

