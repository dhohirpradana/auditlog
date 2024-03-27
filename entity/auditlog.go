package entity

type Request struct {
	Header any
	Body   any
	IP     string
}

type Response struct {
	Header any
	Body   any
	Code   int
}
