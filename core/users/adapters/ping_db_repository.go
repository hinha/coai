package adapters

import "github.com/hinha/coai/internal/store/gorm/mysql"

type PingMysqlRepository struct {
	db *mysql.DB
}

func NewPingMysqlRepository(db *mysql.DB) *PingMysqlRepository {
	return &PingMysqlRepository{db: db}
}

func (p *PingMysqlRepository) Ping() error {
	return p.db.Ping()
}
