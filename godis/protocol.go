package godis

import
(
	"errors"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/ugorji/go/codec"
	"log"
)




type ProtoType map[interface{}]interface{}
type HandleServerFunc *func(args ProtoType) (ProtoType, error)
type HandleClientFunc *func(args ProtoType) (interface{}, error)
// Event representation
type Event struct {
	MsgId  string
	Args   ProtoType
}


type Resp  struct {
	Code   int64
	Data   ProtoType
	ErrMsg string
}


// Returns a pointer to a new event,
// a UUID V4 message_id is generated
func newEvent(args ProtoType) (Event, error) {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	e := Event{
		MsgId:  id.String(),
		Args:   args,
	}
	return e, nil
}


func newResp(code int64, err string, data ProtoType) (Resp){
	r := Resp{
		Code: code,
        Data: data,
	}
    if err!="" {
		r.ErrMsg=err
	}
	log.Printf("r: %v\n",r)
	return r
}


// Packs an event into MsgPack bytes
func (r *Resp) packBytes() ([]byte, error) {
	data := make([]interface{}, 2)
	data[0] = r.Code
	data[1] = r.Data
	if len(r.ErrMsg)>0{
		data = append(data, r.ErrMsg)
	}
	log.Printf("data: %v\n",data)
   return encode(data)
}



// Packs an event into MsgPack bytes
func (e *Event) packBytes() ([]byte, error) {
	data := make([]interface{}, 2)
	data[0] = e.MsgId
	data[1] = e.Args
    return encode(data)
}

func encode(data []interface{}) ([]byte, error) {

	var buf []byte

	enc := codec.NewEncoderBytes(&buf, &codec.MsgpackHandle{})
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return buf, nil
}

func decode(b []byte) (interface{}, error) {

	var mh codec.MsgpackHandle
	var v interface{}
	dec := codec.NewDecoderBytes(b, &mh)

	err := dec.Decode(&v)
	if err != nil {
		return nil, err
	}
	return v, nil
}


//// Unpacks an event fom MsgPack bytes
func unPackEventBytes(b []byte) (*Event, error) {

	v, err :=decode(b)
	if err != nil {
		return nil, err
	}

	// get the event headers
	h, ok := v.([]interface{})[0].([]byte)
	if !ok {
		return nil, errors.New("zerorpc/event interface conversion error")
	}
	// get the event args
	args :=convertValue(v.([]interface{})[1])

	e := Event{
		MsgId: string(h),
		Args:  args.(map[interface{}]interface{}),
	}

	return &e, nil
}


func unPackRespByte(b []byte) (*Resp, error) {

	v, err := decode(b)
	if err != nil {
		return nil, err
	}
	log.Printf("v: %v\n",v)
    code := convertValue(v.([]interface{})[0])
	data := convertValue(v.([]interface{})[1])
	r := Resp{
		Code: code.(int64),
	}
	if (data!=nil){
		r.Data = data.(map[interface{}]interface{})
	}
	if (len(v.([]interface{}))>2){
		er := convertValue(v.([]interface{})[2])
		r.ErrMsg = er.(string)
	}

	return &r, nil
}
//
//// Returns a pointer to a new heartbeat event
//func newHeartbeatEvent() (*Event, error) {
//	ev, err := newEvent("_zpc_hb", nil)
//	if err != nil {
//		return nil, err
//	}
//
//	return ev, nil
//}

// converts an interface{} to a type
func convertValue(v interface{}) interface{} {
	var out interface{}

	switch t := v.(type) {
	case []byte:
		out = string(t)

	case []interface{}:
		for i, x := range t {
			t[i] = convertValue(x)
		}

		out = t

	case map[interface{}]interface{}:
		for key, val := range v.(map[interface{}]interface{}) {
			t[key] = convertValue(val)
		}
		out = t

	default:
		out = t
	}

	return out
}
