package storage

var (
	// BitUSD 美元
	BitUSD = "bitUSD"
	// BitCNY 人民币
	BitCNY = "bitCNY"
)

var (
	// BitUSDSymbol 美元符号
	BitUSDSymbol = "USD"
	// BitCNYSymbol 人民币符号
	BitCNYSymbol = "CNY"
)

// GetAsset 获取资产名
func GetAsset(symbol string) string {
	switch symbol {
	case BitUSDSymbol:
		return BitUSD
	case BitCNYSymbol:
		return BitCNY
	default:
		return ""
	}
}

// GetAssetSymbol 获取资产符号
func GetAssetSymbol(asset string) string {
	switch asset {
	case BitUSD:
		return BitUSDSymbol
	case BitCNY:
		return BitCNYSymbol
	default:
		return ""
	}
}

// Asset 资产信息
type Asset struct {
	Asset  string `json:"asset"`  // 资产名称
	Amount uint32 `json:"amount"` // 资产总额
	Freeze uint32 `json:"freeze"` // 冻结资产
}

// RedEnvelope 红包信息
type RedEnvelope struct {
	ID         uint64 `json:"id"`          // 红包ID
	GroupID    int64  `json:"group_id"`    // 群组ID
	MessageID  int32  `json:"message_id"`  // 消息ID
	SenderID   int64  `json:"sneder_id"`   // 发送者
	SenderName string `json:"sneder_name"` // 发送者名字
	Asset      string `json:"asset"`       // 资产类型
	Amount     uint32 `json:"amount"`      // 红包总额
	Received   uint32 `json:"received"`    // 领取金额
	Number     uint32 `json:"number"`      // 红包个数
	Lucky      bool   `json:"lucky"`       // 是否随机
	Value      uint32 `json:"value"`       // 单个价值
	Active     bool   `json:"active"`      // 是否激活
	Memo       string `json:"memo"`        // 红包留言
	Timestamp  int64  `json:"timestamp"`   // 时间戳
}

// RedEnvelopeUser 红包用户
type RedEnvelopeUser struct {
	UserID    int32  `json:"user_id"`    // 用户ID
	FirstName string `json:"first_name"` // 用户名
}

// RedEnvelopeRecord 红包记录
type RedEnvelopeRecord struct {
	Value int              `json:"value"`          // 红包金额
	User  *RedEnvelopeUser `json:"user,omitempty"` // 用户信息
}
