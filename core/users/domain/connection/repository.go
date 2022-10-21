package connection

type Repository interface {
	Ping() error
}
