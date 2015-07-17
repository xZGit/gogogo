package godis


import(
	"sync"
    "log"
	"time"
)



// Task handler representation
type taskHandler struct {
	TaskId    string
	HandlerFunc HandleClientFunc
}

type Client struct {
	id string
	redisClient   *RedisClient
	mutex         sync.Mutex
	HandleTasks [] *taskHandler
	hasListen      bool
}



func NewClient(id string, host string) (*Client) {
	client:=Client{
		id:id,
		redisClient:NewRedisClient(host),
		HandleTasks :make([]*taskHandler, 0),
		hasListen:false,
	}

	return &client
}



func (c *Client) Call(name string, handlerFunc HandleClientFunc, args ProtoType,n int)  error {
    c.mutex.Lock()
	defer c.mutex.Unlock()
	if(!c.hasListen){
		go c.Listen()
		c.hasListen = true
	}
	go func() {
		c.redisClient.getActiveSubPub()
		event, err := newEvent(c.id, args)
		if err != nil {
			panic(err)
		}
		msg, err := event.packBytes()
		task:=taskHandler{
			TaskId:event.MsgId,
			HandlerFunc: handlerFunc,
		}
		c.HandleTasks=append(c.HandleTasks,&task)
		c.redisClient.pubConn.Publish(name, msg)
	}()
	return nil
}


func (c *Client) Listen() {
     	i:=0
		rec := make(chan []string)
		go c.redisClient.subConn.Subscribe(rec, c.id)
		var ls []string
		for {
			select {
			case ls = <-rec:
				if ls[0]=="message" && len(ls)>2 {
					i=i+1
					log.Printf(": %v\n", i)
					go c.ProcessFunc(ls[2])
				}
			case <-time.After(5 * time.Second):
//			 	println("timeout")
				break
			}
		}
}



func (c *Client) ProcessFunc(ls string) {

	mySlice := []byte(ls)

	resp, err := unPackRespByte(mySlice)
	if err != nil {
		log.Printf("err: %v\n", err)
	}


	for _, h := range  c.HandleTasks{
		if h.TaskId == resp.MsgId{
			(*h.HandlerFunc)(resp.RespInfo)
		}
	}

}

