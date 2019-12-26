package mysql

import "etherscan-go/types"

func (m *mysql) GetCoinByAddress(contractAddress string) (*types.Coin, error) {
	var result types.Coin
	err := m.db.Where("status = 1 and contract_address = ?", contractAddress).First(&result).Error
	return &result, handleErr(err, "coin")
}
