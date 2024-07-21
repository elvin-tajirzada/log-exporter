package database

type Client interface {
	Push(log []byte) error
}
