package entity

type Request struct {
	Header any
	Body   string
}

type Response struct {
	Header any
	Data   string
	Code   int
}
