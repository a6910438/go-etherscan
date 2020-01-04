package mysql

import "etherscan-go/types"

func (m *mysql) GetConfigByWithdrawNumber() (*types.Config, error) {
	var result types.Config
	err := m.db.Where("c_key = 'than_withdraw_number'").First(&result).Error
	return &result, handleErr(err, "config")
}
