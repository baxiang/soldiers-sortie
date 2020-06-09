package transport

type serverTransport struct {
	opts *ServerTransportOptions
}

var serverTransportMap = make(map[string]ServerTransportOption)

func init() {
	serverTransportMap["default"] = DefaultServerTransport
}

var DefaultServerTransport = NewServerTransport()

var NewServerTransport = func() ServerTransport {
	return &serverTransport{
		opts: &ServerTransportOptions{},
	}
}
