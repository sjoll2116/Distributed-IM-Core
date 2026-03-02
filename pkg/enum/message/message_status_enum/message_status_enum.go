package message_status_enum

const (
	// 未发送
	Unsent = iota
	// 已发送（服务端已通过 WebSocket 下发）
	Sent
	// 已确认（客户端已回传 ACK 回执）
	Acked
)
