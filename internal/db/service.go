package db

type Service interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, bool, error)
}
