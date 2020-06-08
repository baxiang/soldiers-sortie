package codec


type Serialization interface {
	Marshal(interface{})([]byte,error)
	Unmarshal([]byte,interface{})error
}

const (
	Proto   = "proto"   // protobuf
	MsgPack = "msgpack" // msgpack
	Json    = "json"    // json
)

var serializationMap = make(map[string]Serialization)