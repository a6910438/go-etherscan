package types

import "github.com/shopspring/decimal"

// 用户资产表
type UserAssest struct {
	Id            int64           `json:"id"`
	UserId        int64           `json:"user_id"`        // 用户ID
	CoinId        int64           `json:"coin_id"`        // 币种ID
	AvaBalance    decimal.Decimal `json:"ava_balance"`    // 可用余额
	FreezeBalance decimal.Decimal `json:"freeze_balance"` // 冻结余额
	Address       string          `json:"address"`        // 地址
	CreateTime    int64           `json:"create_time"`    // 创建时间
	UpdateTime    int64           `json:"update_time"`    // 更新时间
}
