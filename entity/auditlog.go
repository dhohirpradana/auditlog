package entity

type Request struct {
	Header any
	Body   any
}

type Response struct {
	Header any
	Body   any
	Code   int
}
