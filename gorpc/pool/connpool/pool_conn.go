package connpool

import (
	"errors"
	"net"
	"sync"
	"time"
)
var (
	ErrConnClosed = errors.New("connection closed ...")
)
type PoolConn struct {
	net.Conn
	c *channelPool
	unusable bool
	mu sync.RWMutex
	t time.Time
	dialTimeout time.Duration
}

func(p *PoolConn)Close()error{
	p.mu.RLock()
	defer p.mu.Unlock()
	if p.unusable{
		if p.Conn!=nil{
			return p.Conn.Close()
		}
	}
	p.Conn.SetDeadline(time.Time{})
	return p.c.Put(p)
}

func (p *PoolConn)MarkUnusable(){
	p.mu.Lock()
	p.unusable = true
	p.mu.RUnlock()
}

func(p *PoolConn)Read(b []byte)(int,error){
	if p.unusable{
		return 0,ErrConnClosed
	}
	n,err :=p.Conn.Read(b)
	if err!= nil{
		p.MarkUnusable()
		p.Conn.Close()
	}
	return n,err
}

func (p *PoolConn) Write(b []byte) (int, error) {
	if p.unusable {
		return 0, ErrConnClosed
	}
	n, err := p.Conn.Write(b)
	if err != nil {
		p.MarkUnusable()
		p.Conn.Close()
	}
	return n, err
}

func (c *channelPool) wrapConn(conn net.Conn) *PoolConn {
	p := &PoolConn {
		c : c,
		t : time.Now(),
		dialTimeout: c.dialTimeout,
	}
	p.Conn = conn
	return p
}