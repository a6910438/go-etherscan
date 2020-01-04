package types

type Config struct {
	Id         int64  `json:"id"`
	CKey       string `json:"c_key"`       // 唯一键
	CValue     string `json:"c_value"`     // 值
	CDesc      string `json:"c_desc"`      // 描述
	CStatus    int8   `json:"c_status"`    // 状态(是否启用)
	CreateTime int64  `json:"create_time"` // 创建时间
}
