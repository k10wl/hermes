package db

type Client interface{}

type DB struct{}

func NewDB() *DB {
	return &DB{}
}
