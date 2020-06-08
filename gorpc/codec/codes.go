package codec


type Codec interface {
	Encode([]byte)([]byte,error)
	Decode([]byte)([]byte,error)
}

const FrameHeadLen = 15

const Magic = 0x11
const Version  = 0