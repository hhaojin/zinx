package zinx

type HandleFunc func(req *Request)

type Router struct {
	Routes map[int]HandleFunc
}

func (r *Router) Add(typ int, fn HandleFunc) {
	r.Routes[typ] = fn
}

func (r *Router) Handle(req *Request) {
	req.Conn.WriteChan <- []byte("hello world")
}
