package model

type Image struct {
	Id       int64
	Filename string
	Data     []byte
	TempDir  string
}
