package zinx

import (
	"fmt"
	"net"
)

type Server struct {
	Name   string
	Router *Router
}

func New(name string) *Server {
	return &Server{
		Name:   name,
		Router: &Router{Routes: make(map[int]HandleFunc)},
	}
}

func resolveAddr(addr ...string) string {
	if len(addr) == 0 {
		return ":8999"
	}
	return addr[0]
}

func (s *Server) Run(addr ...string) {
	address := resolveAddr(addr...)
	fmt.Println("zinx listening ", address)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		panic(err)
	}
	l, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		panic(err)
	}
	var connId uint32
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			continue
		}
		connId++
		go s.serve(connId, conn)
	}
}

func (s *Server) Close() {

}

func (s *Server) serve(cid uint32, c *net.TCPConn) {
	dialConn := NewConnection(c, cid, s.Router)
	go dialConn.ReadLoop()
	go dialConn.WriteLoop()
}

func (s *Server) AddRouter(typ int, handler HandleFunc) {
	s.Router.Add(typ, handler)
}
