package model

type LoginPassSecret struct {
	Id         int
	Title      string
	RecordType int
	Login      string
	Password   string
}

type TextSecret struct {
	Id         int
	Title      string
	RecordType int
	Text       string
}

type FileSecret struct {
	Id         int
	Title      string
	RecordType int
	Path       string
	Binary     []byte
}

type CardSecret struct {
	Id         int
	Title      string
	RecordType int
	CardNumber string
	CVV        string
	Due        string
}
