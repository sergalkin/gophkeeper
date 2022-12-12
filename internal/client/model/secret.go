package model

import "time"

type SecretList struct {
	Id    int
	Title string
}

type LoginPassSecret struct {
	Id         int
	Title      string
	RecordType int
	Login      string
	Password   string
	UpdatedAt  time.Time
}

type TextSecret struct {
	Id         int `json:"id"`
	Title      string
	RecordType int
	Text       string
	UpdatedAt  time.Time
}

type FileSecret struct {
	Id         int
	Title      string
	RecordType int
	Path       string
	Binary     []byte
	UpdatedAt  time.Time
}

type CardSecret struct {
	Id         int
	Title      string
	RecordType int
	CardNumber string
	CVV        string
	Due        string
	UpdatedAt  time.Time
}
