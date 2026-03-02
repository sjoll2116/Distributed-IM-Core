package respond

type LoginRespond struct {
	Uuid      string `json:"uuid"`
	Nickname  string `json:"nickname"`
	Telephone string `json:"telephone"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	Gender    int8   `json:"gender"`
	Birthday  string `json:"birthday"`
	Signature string `json:"signature"`
	IsAdmin   int8   `json:"is_admin"`
	Status    int8   `json:"status"`    // 状态
	CreatedAt string `json:"createdAt"` // 注册时间
	Token     string `json:"token"`     // JWT Token
}
