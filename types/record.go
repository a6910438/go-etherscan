package types

import "github.com/shopspring/decimal"

// 流水表
type Record struct {
	Id         int64           `json:"id"`
	Send       int64           `json:"send"`        // 发送用户ID
	Type       uint8           `json:"type"`        // 类型(1. 充值)
	Old        decimal.Decimal `json:"old"`         // 以前金额
	Now        decimal.Decimal `json:"now"`         // 现在金额
	Amount     decimal.Decimal `json:"amount"`      // 金额
	Receice    int64           `json:"receice"`     // 接受ID
	CreateTime int64           `json:"create_time"` // 创建时间
	Direction  string          `json:"direction"`   // 方向
}
