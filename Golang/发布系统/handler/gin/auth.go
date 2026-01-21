package handler

import (
	"time"
	"happyladysauce/utils/jwt"

	"github.com/gin-gonic/gin"
)

// 常量定义
const (
	JWTSecretKey = "news"    // JWT签名密钥
	JWTSubject   = "news"    // JWT签发者
	CookieDomain = "localhost" // Cookie域名
	TokenCookieName = "token"  // Token Cookie名称
)

// generateAndSetToken 生成JWT令牌并设置到Cookie中
func generateAndSetToken(ctx *gin.Context, uid int) error {
	// 生成JWT令牌
	header := jwt.DefaultJwtHeader()
	payload := jwt.JwtPayload{
		Issue:       JWTSubject,
		IssueAt:     time.Now().Unix(),
		Expiration:  time.Now().Add(24 * time.Hour).Unix(),
		UserDefined: map[string]any{"uid": uid},
	}

	token, err := jwt.GenJWT(*header, payload, JWTSecretKey)
	if err != nil {
		return err
	}

	// 设置Cookie
	expirationTime := time.Until(time.Unix(payload.Expiration, 0))
	ctx.SetCookie(
		TokenCookieName,
		token,
		int(expirationTime.Seconds()),
		"/",
		CookieDomain,
		false, // 开发环境可设为false，生产环境应设为true（HTTPS）
		true,  // 禁止JavaScript访问
	)

	return nil
}

// GetUidFromJWT 从JWT令牌中提取用户ID
func GetUidFromJWT(token string) int {
	_, payload, err := jwt.ParseJWT(token, "news")
	if err != nil {
		return 0
	}
	for k, v := range payload.UserDefined {
		if k == "uid" {
			if uid, ok := v.(int); ok {
				return uid
			}
		}
	}
	return 0
}

// GetLoginUid 从请求的Cookie中获取当前登录用户的ID
func GetLoginUid(ctx *gin.Context) int {
	// 修复：从名为'token'的cookie中获取，与user.go中设置的保持一致
	token, err := ctx.Cookie("token")
	if err != nil {
		return 0
	}
	return GetUidFromJWT(token)
}

// Auth 认证中间件，验证用户是否已登录
func Auth(ctx *gin.Context) {
	loginUid := GetLoginUid(ctx)
	if loginUid == 0 {
		ctx.JSON(401, gin.H{
			"message": "未登录",
		})
		ctx.Abort()
		return
	}
	ctx.Set("uid", loginUid)
	ctx.Next()
}

