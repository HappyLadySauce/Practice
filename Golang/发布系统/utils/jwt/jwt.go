package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// JwtHeader 令牌头
type JwtHeader struct {
	Algo string `json:"alg"`	// 哈希算法, 默认为HMAC SHA256(写为 HS256)
	Type string `json:"typ"`	// 令牌(token)类型, 统一写为 JWT
}

// JwtPayload 令牌载荷
type JwtPayload struct {
	ID			string			`json:"jti"`	// 令牌ID, 用于唯一标识该令牌
	Issue		string			`json:"iss"`	// 令牌签发者
	Audience	string			`json:"aud"`	// 令牌受众
	Subject		string			`json:"sub"`	// 令牌主题
	IssueAt		int64			`json:"iat"`	// 令牌签发时间
	NotBefore	int64			`json:"nbf"`	// 令牌生效时间
	Expiration	int64			`json:"exp"`	// 令牌过期时间
	UserDefined	map[string]any	`json:"ud"`	// 用户自定义字段
}

func DefaultJwtHeader() *JwtHeader {
	return &JwtHeader{
		Algo: "HS256",
		Type: "JWT",
	}
}

func GenJWT(header JwtHeader, payload JwtPayload, secret string) (string, error) {
	var headerStr, payloadStr, signature string
	if bs1, err := json.Marshal(header); err != nil {
		return "", err
	} else {
		headerStr = base64.RawURLEncoding.EncodeToString(bs1)
	}

	if bs2, err := json.Marshal(payload); err != nil {
		return "", err
	} else {
		payloadStr = base64.RawURLEncoding.EncodeToString(bs2)
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(headerStr + "." + payloadStr))
	signature = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return headerStr + "." + payloadStr + "." + signature, nil
}

func ParseJWT(token string, secret string) (*JwtHeader, *JwtPayload, error) {
	// 解析token
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("token格式错误")
	}

	// 验证签名
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(parts[0] + "." + parts[1]))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if signature != parts[2] {
		return nil, nil, fmt.Errorf("token签名错误")
	}

	// 解析header
	var header JwtHeader
	var payload JwtPayload
	if bs1, err := base64.RawURLEncoding.DecodeString(parts[0]); err != nil {
		return &header, &payload, err
	} else {
		if err := json.Unmarshal(bs1, &header); err != nil {
			return &header, &payload, err	
		}
	}
	if bs2, err := base64.RawURLEncoding.DecodeString(parts[1]); err != nil {
		return &header, &payload, err
	} else {
		if err := json.Unmarshal(bs2, &payload); err != nil {
			return &header, &payload, err
		}
	}
	return &header, &payload, nil
}

