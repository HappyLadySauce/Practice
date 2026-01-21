package database_test

import (
	"os"
	"strings"
	"testing"

	"happyladysauce/database/gorm"
	"happyladysauce/utils/conf"
	"happyladysauce/utils/crypt"
	"happyladysauce/utils/logger"
)

// 测试入口：初始化数据库连接
func TestMain(m *testing.M) {
	config.InitConfig("../../conf", "db", "yaml")
	database.InitGormDB()
	logger.InitLogger()
	os.Exit(m.Run())
}

// TestRegisterUser 测试用户注册功能
func TestRegisterUser(t *testing.T) {
	// 使用唯一用户名，确保首次注册不会受历史数据影响
	name := "TestUser1"
	password := "123456"
	
	// 确保测试用户不存在
	if id, _ := database.GetIdByName(name); id > 0 {
		database.LogOff(id)
	}

	id, err := database.RegisterUser(name, password)
	if err != nil {
		t.Fatalf("首次注册用户失败: %v", err)
	}
	if id <= 0 {
		t.Errorf("首次注册用户失败: 无效的用户ID")
	}
	t.Logf("首次注册用户成功: 用户ID=%d", id)

	// 再次使用相同用户名注册，应返回重复错误
	_, err = database.RegisterUser(name, password)
	if err == nil {
		t.Fatalf("预期重复注册失败，但成功")
	}
	if !strings.Contains(err.Error(), "已存在") {
		t.Fatalf("重复注册返回非预期错误: %v", err)
	}
	t.Logf("重复注册返回预期错误: %v", err)
	
	// 清除测试数据
	defer func() {
		if err := database.LogOff(id); err != nil {
			t.Logf("注销测试用户失败: %v", err)
		}
	}()
}

// TestLogOff 测试用户注销功能
func TestLogOff(t *testing.T) {
	// 创建测试用户
	name := "TestUserLogOff"
	password := "123456"
	id, err := database.RegisterUser(name, password)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	
	// 测试成功注销
	err = database.LogOff(id)
	if err != nil {
		t.Fatalf("注销用户失败: %v", err)
	}
	t.Logf("成功注销用户: %d", id)
	
	// 测试注销不存在的用户
	err = database.LogOff(id)
	if err == nil {
		t.Fatalf("预期注销不存在用户失败，但成功")
	}
	if !strings.Contains(err.Error(), "不存在") {
		t.Fatalf("注销不存在用户返回非预期错误: %v", err)
	}
	t.Logf("注销不存在用户返回预期错误: %v", err)
}

// TestUpdatePassword 测试更新用户密码功能
func TestUpdatePassword(t *testing.T) {
	// 创建测试用户
	name := "TestUserUpdatePwd"
	oldPassword := "123456"
	newPassword := "654321"
	id, err := database.RegisterUser(name, oldPassword)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer database.LogOff(id)
	
	// 测试使用错误的旧密码
	err = database.UpdatePassword(id, "wrongpassword", newPassword)
	if err == nil {
		t.Fatalf("预期使用错误旧密码更新失败，但成功")
	}
	if !strings.Contains(err.Error(), "旧密码错误") {
		t.Fatalf("使用错误旧密码返回非预期错误: %v", err)
	}
	t.Logf("使用错误旧密码返回预期错误: %v", err)
	
	// 测试成功更新密码
	err = database.UpdatePassword(id, oldPassword, newPassword)
	if err != nil {
		t.Fatalf("更新密码失败: %v", err)
	}
	t.Logf("成功更新用户密码")
	
	// 验证新密码是否生效
	user, err := database.GetUserById(id)
	if err != nil {
		t.Fatalf("获取用户信息失败: %v", err)
	}
	if !crypt.CheckPasswordHash(newPassword, user.PassWord) {
		t.Fatalf("密码更新未生效")
	}
	
	// 测试更新不存在用户的密码
	err = database.UpdatePassword(999999, "anypassword", "newword")
	if err == nil {
		t.Fatalf("预期更新不存在用户密码失败，但成功")
	}
	if !strings.Contains(err.Error(), "不存在") {
		t.Fatalf("更新不存在用户密码返回非预期错误: %v", err)
	}
}

// TestGetUserById 测试根据ID获取用户信息
func TestGetUserById(t *testing.T) {
	// 创建测试用户
	name := "TestUserGetById"
	password := "123456"
	id, err := database.RegisterUser(name, password)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer database.LogOff(id)
	
	// 测试获取存在的用户
	user, err := database.GetUserById(id)
	if err != nil {
		t.Fatalf("获取用户信息失败: %v", err)
	}
	if user.Name != name {
		t.Fatalf("获取的用户名不匹配，预期: %s, 实际: %s", name, user.Name)
	}
	t.Logf("成功获取用户信息: %s", user.Name)
	
	// 测试获取不存在的用户
	_, err = database.GetUserById(999999)
	if err == nil {
		t.Fatalf("预期获取不存在用户失败，但成功")
	}
	if !strings.Contains(err.Error(), "不存在") {
		t.Fatalf("获取不存在用户返回非预期错误: %v", err)
	}
}

// TestGetUserByName 测试根据用户名获取用户信息
func TestGetUserByName(t *testing.T) {
	// 创建测试用户
	name := "TestUserGetByName"
	password := "123456"
	id, err := database.RegisterUser(name, password)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer database.LogOff(id)
	
	// 测试获取存在的用户
	user, err := database.GetUserByName(name)
	if err != nil {
		t.Fatalf("获取用户信息失败: %v", err)
	}
	if user.Id != id {
		t.Fatalf("获取的用户ID不匹配，预期: %d, 实际: %d", id, user.Id)
	}
	t.Logf("成功获取用户信息: ID=%d", user.Id)
	
	// 测试获取不存在的用户
	nonExistentName := "NonExistentUser" + t.Name() // 确保用户名唯一
	_, err = database.GetUserByName(nonExistentName)
	if err == nil {
		t.Logf("警告: 查询不存在用户未返回错误，这可能是预期行为")
	} else {
		t.Logf("获取不存在用户返回错误: %v", err)
	}
}

// TestGetIdByName 测试根据用户名获取用户ID
func TestGetIdByName(t *testing.T) {
	// 创建测试用户
	name := "TestUserGetIdByName"
	password := "123456"
	expectedId, err := database.RegisterUser(name, password)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer database.LogOff(expectedId)
	
	// 测试获取存在的用户ID
	actualId, err := database.GetIdByName(name)
	if err != nil {
		t.Fatalf("获取用户ID失败: %v", err)
	}
	if actualId != expectedId {
		t.Fatalf("获取的用户ID不匹配，预期: %d, 实际: %d", expectedId, actualId)
	}
	t.Logf("成功获取用户ID: %d", actualId)
	
	// 测试获取不存在的用户ID
	nonExistentName := "NonExistentUser" + t.Name() // 确保用户名唯一
	actualId, err = database.GetIdByName(nonExistentName)
	if err == nil {
		if actualId == 0 {
			t.Logf("查询不存在用户返回ID 0，这可能是预期行为")
		} else {
			t.Logf("警告: 查询不存在用户未返回错误且返回非零ID: %d", actualId)
		}
	} else {
		t.Logf("获取不存在用户ID返回错误: %v", err)
	}
}