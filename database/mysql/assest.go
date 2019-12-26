package mysql

import (
	"etherscan-go/types"
	"github.com/pkg/errors"
	"time"
)

func (m *mysql) GetAssestByAddress(to string) (*types.UserAssest, error) {
	var result types.UserAssest
	err := m.db.Where("address = ?", to).First(&result).Error
	return &result, handleErr(err, "coin")
}

func (m *mysql) UpdateAssestById(assest *types.UserAssest) (err error) {
	now := time.Now().Unix()
	tx := m.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if err := tx.Model(&assest).Update(map[string]interface{}{"ava_balance": assest.AvaBalance, "update_time": now}).Error; err != nil {
		return errors.WithMessage(err, "update assest")
	}
	return nil
}
