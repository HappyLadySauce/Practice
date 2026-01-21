package handler

import (
	"happyladysauce/database/gorm"
	"happyladysauce/handler/model"
	"happyladysauce/utils/validation"
	"net/http"

	"happyladysauce/utils/crypt"

	"github.com/gin-gonic/gin"
)

// RegistUser 用户注册接口
func RegistUser(ctx *gin.Context) {
	var request model.User
	// 绑定并验证请求参数
	if err := ctx.ShouldBind(&request); err != nil {
		respondWithError(ctx, http.StatusBadRequest, validation.ProcessErr(err))
		return
	}

	// 调用数据库注册用户
	id, err := database.RegisterUser(request.Name, request.PassWord)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 生成并设置JWT令牌
	if err := generateAndSetToken(ctx, id); err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	respondWithSuccess(ctx, "注册成功")
}

// Login 用户登录接口
func Login(ctx *gin.Context) {
	var request model.User
	// 绑定并验证请求参数
	if err := ctx.ShouldBind(&request); err != nil {
		respondWithError(ctx, http.StatusBadRequest, validation.ProcessErr(err))
		return
	}

	// 根据用户名查询用户
	validationUser, err := database.GetUserByName(request.Name)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if validationUser == nil {
		respondWithError(ctx, http.StatusBadRequest, "用户不存在")
		return
	}

	// 验证密码
	if !crypt.CheckPasswordHash(request.PassWord, validationUser.PassWord) {
		respondWithError(ctx, http.StatusBadRequest, "密码错误")
		return
	}

	// 生成并设置JWT令牌
	if err := generateAndSetToken(ctx, validationUser.Id); err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	respondWithSuccess(ctx, "登录成功")
}

// UpdatePassword 更新用户密码接口
func UpdatePassword(ctx *gin.Context) {
	// 从上下文获取用户ID
	uid, ok := ctx.Value("uid").(int)
	if !ok {
		respondWithError(ctx, http.StatusBadRequest, "请先登录")
		return
	}

	var request model.UpdatePassword
	// 绑定并验证请求参数
	if err := ctx.ShouldBind(&request); err != nil {
		respondWithError(ctx, http.StatusBadRequest, validation.ProcessErr(err))
		return
	}

	// 验证新旧密码是否相同
	if request.OldPassword == request.NewPassword {
		respondWithError(ctx, http.StatusBadRequest, "新密码不能与旧密码相同")
		return
	}

	// 更新密码
	if err := database.UpdatePassword(uid, request.OldPassword, request.NewPassword); err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	respondWithSuccess(ctx, "更新用户密码成功")
}

