package mysql

import (
	"etherscan-go/types"
	"github.com/pkg/errors"
	"time"
)

func (m *mysql) AddRecord(record types.Record, tableName string) (id int64, err error) {
	record.CreateTime = time.Now().Unix()
	tx := m.db.Table(tableName).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if err := tx.Create(&record).Error; err != nil {
		return 0, errors.WithMessage(err, "create information")
	}
	return record.Id, nil
}
