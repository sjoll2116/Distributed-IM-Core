package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义全局的 Secret (在生产环境中建议写在 config.yaml 配置文件中读取)
var jwtSecret = []byte("Kamachat_Super_Secret_Key_2026")

// Claims 定义了我们将存入 JWT 载荷的自定义信息
type Claims struct {
	Uuid string `json:"uuid"`
	jwt.RegisteredClaims
}

// GenerateToken 给定用户的 uuid，生成带有效期的 JWT Token
func GenerateToken(uuid string) (string, error) {
	nowTime := time.Now()
	// 设置 Token 有效期为 7 天
	expireTime := nowTime.Add(7 * 24 * time.Hour)

	claims := Claims{
		Uuid: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "kamachat_server",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 验证 JWT Token，如果合法则提取其中的 Claims (例如可获取 user.Uuid)
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, errors.New("invalid or expired token")
}
