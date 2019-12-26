package types

import (
	"github.com/shopspring/decimal"
)

// 充值表
type Deposit struct {
	Id         int64           `json:"id"`
	UserId     int64           `json:"user_id"`     // 用户ID
	CoinId     int64           `json:"coin_id"`     // 币种ID
	Receice    string          `json:"receice"`     // 接收地址
	Fee        decimal.Decimal `json:"fee"`         // 旷工费
	Amount     decimal.Decimal `json:"amount"`      // 金额
	Send       string          `json:"send"`        // 发送地址
	CreateTime int64           `json:"create_time"` // 创建时间
	Hash       string          `json:"hash"`        // 哈希
}
