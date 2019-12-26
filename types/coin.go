package types

import (
	"github.com/shopspring/decimal"
)

// 币种表
type Coin struct {
	Id              int64           `json:"id"`
	FullName        string          `json:"full_name"`        // 全名
	ResumeName      string          `json:"resume_name"`      // 简称
	Rate            decimal.Decimal `json:"rate"`             // 汇率
	Name            string          `json:"name"`             // 显示名称
	Number          uint8           `json:"number"`           // 小数位数
	IconUrl         string          `json:"icon_url"`         // 图标地址
	Remark          string          `json:"remark"`           // 备注
	Status          uint8           `json:"status"`           // 是否启用(0.停用 1.启用)
	Sort            uint8           `json:"sort"`             // 排序
	Fee             float64         `json:"fee"`              // 手续费
	ContractAddress string          `json:"contract_address"` // 合约地址
	IsContract      uint8           `json:"is_contract"`      // 是否为合约币(0.不是 1.是)
	Auto            uint8           `json:"auto"`             // 用户创建时是否自动分配地址
	CreateTime      int64           `json:"create_time"`      // 创建时间
	UpdateTime      int64           `json:"update_time"`      // 更新时间
}

type CoinStatus uint8

const (
	Able CoinStatus = iota
	DisAble
)

type CoinIsContract uint8

const (
	No CoinIsContract = iota
	Yes
)