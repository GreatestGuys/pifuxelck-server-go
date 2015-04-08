package common

type Scannable interface {
	Scan(dest ...interface{}) error
}
