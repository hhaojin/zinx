package zinx

import "C"
import (
	"fmt"
	"io"
	"net"
	"time"
)

type Connection struct {
	isClosed  bool
	ConnID    uint32
	ExitChan  chan bool
	WriteChan chan []byte
	Conn      *net.TCPConn
	r         *Router
}

func NewConnection(conn *net.TCPConn, connId uint32, r *Router) *Connection {
	return &Connection{
		isClosed:  false,
		ConnID:    connId,
		ExitChan:  make(chan bool, 1),
		WriteChan: make(chan []byte),
		Conn:      conn,
		r:         r,
	}
}

func (c *Connection) ReadLoop() {
	for {
		buf := make([]byte, 1024)
		n, err := c.Conn.Read(buf)

		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				fmt.Printf("conn read error: %v, retrying\n", ne.Error())
				time.Sleep(time.Millisecond * 200)
				continue
			}
			if err != io.EOF {
				fmt.Println("conn read err: ", err)
			}
			c.ExitChan <- true
			c.Conn.Close()
			break
		}
		go func() {
			fmt.Println("read conn: ", string(buf[:n]))
			req := &Request{
				Conn: c,
				data: buf,
			}
			c.r.Handle(req)
			c.WriteChan <- buf[:n]
		}()
	}

}

func (c *Connection) WriteLoop() {
	for {
		select {
		case buf := <-c.WriteChan:
			_, err := c.Conn.Write(buf)
			if err != nil {
				fmt.Println("write conn err: ", err)
			}
		case <-c.ExitChan:
			fmt.Println("write loop exit")
			close(c.WriteChan)
			return
		}
	}

}
