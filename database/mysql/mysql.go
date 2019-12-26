package mysql

import (
	"fmt"

	"etherscan-go/cmd/app/config"
	mErrors "etherscan-go/error"
	log "github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

type mysql struct {
	db  *gorm.DB
	log log.Logger
}

func NewMysql(config config.DbConfig) (*mysql, error) {
	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?sql_notes=false&parseTime=true&loc=Local&charset=utf8mb4,utf8",
			config.User, config.Password, config.Host, config.Port, config.Dbname),
	)
	if err != nil {
		return nil, errors.WithMessage(err, "mysql open")
	}
	db.SingularTable(true) //设置框架表明为单数
	return &mysql{
		db:  db,
		log: log.New("method", "mysql"),
	}, nil
}

func handleErr(err error, msg string) error {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return mErrors.ErrRecordNotFound
	default:
		return errors.WithMessage(err, msg)
	}
}
