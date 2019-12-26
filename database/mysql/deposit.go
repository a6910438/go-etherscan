package mysql

import (
	"etherscan-go/types"
	"github.com/pkg/errors"
	"time"
)

func (m *mysql) Add(deposit types.Deposit) (id int64, err error) {
	deposit.CreateTime = time.Now().Unix()
	tx := m.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if err := tx.Create(&deposit).Error; err != nil {
		return 0, errors.WithMessage(err, "create information")
	}
	return deposit.Id, nil
}

func (m *mysql) GetDepositByHash(hash string) (*types.Deposit, error) {
	var result types.Deposit
	err := m.db.Where("hash = ?", hash).First(&result).Error
	return &result, handleErr(err, "deposit")
}
