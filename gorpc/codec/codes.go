package codec

import (
	"bytes"
	"encoding/binary"
	"math"
	"sync"
	"github.com/golang/protobuf/proto"
)

//编解码
type Codec interface {
	Encode([]byte)([]byte,error) // 编码
	Decode([]byte)([]byte,error) //解码
}

const FrameHeadLen = 15
const Magic = 0x11
const Version  = 0

type FrameHeader struct {
	Magic uint8
	Version uint8
	MsgType uint8
	ReqType uint8
	CompressType uint8
	StreamID uint8
	Length uint32
	Reserved uint32
}

var codecMap = make(map[string]Codec)
// GetCodec get a Codec by a codec name
func GetCodec(name string) Codec {
	if codec, ok := codecMap[name]; ok {
		return codec
	}
	return DefaultCodec
}
var DefaultCodec = NewCodec()

var NewCodec = 	func () Codec {
	return &defaultCodec{}
}
type defaultCodec struct {
}

func(c *defaultCodec)Encode(data []byte)([]byte,error){
	totalLen :=FrameHeadLen+len(data)
	buffer := bytes.NewBuffer(make([]byte,0,totalLen))
	frame := FrameHeader{
		Magic:        Magic,
		Version:      Version,
		MsgType:      0x0,
		ReqType:      0x0,
		CompressType: 0x0,
		Length:       uint32(len(data)),
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.Magic); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Version); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.MsgType); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.ReqType); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.CompressType); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.StreamID); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Length); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Reserved); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (c *defaultCodec) Decode(frame []byte) ([]byte,error) {
	return frame[FrameHeadLen:], nil
}

func init(){
	RegisterCodec("proto",DefaultCodec)
}

func RegisterCodec(name string,codec Codec){
	if codecMap ==nil{
		codecMap = map[string]Codec{}
	}
	codecMap[name] = codec
}

func upperLimit(val int) uint32 {
	if val > math.MaxInt32 {
		return uint32(math.MaxInt32)
	}
	return uint32(val)
}

var bufferPool = &sync.Pool{
	New : func() interface {} {
		return &cachedBuffer {
			Buffer : proto.Buffer{},
			lastMarshaledSize : 16,
		}
	},
}

type cachedBuffer struct {
	proto.Buffer
	lastMarshaledSize uint32
}