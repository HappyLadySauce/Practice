package database

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"happyladysauce/database/model"
	"happyladysauce/utils/crypt"
	"happyladysauce/utils/logger"
)

func RegisterUser(name, password string) (int, error) {
	user := model.User{
		Name: name,
		PassWord: password,
	}

	// 密码加密
	user.PassWord = crypt.BcryptPassword(user.PassWord)

	if err := PostDB.Create(&user).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return 0, fmt.Errorf("用户名[%s]已存在", name)
			}
		}
		logger.Log.Error("注册用户失败", "name", name, "err", err)
		return 0, errors.New("注册用户失败")
	}
	logger.Log.Info("注册用户成功", "name", name, "id", user.Id)
	return user.Id, nil
}

// LogOff 注销用户
func LogOff(uid int) error {
	tx := PostDB.Delete(&model.User{}, "id = ?", uid)
	if err := tx.Error; err != nil {
		logger.Log.Error("注销用户失败", "uid", uid, "err", err)
		return errors.New("注销用户失败")
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("用户uid:[%d]不存在", uid)
	}
	logger.Log.Info("注销用户成功", "uid", uid)
	return nil
}

// 更新用户密码
func UpdatePassword(uid int, oldPassword, newPassword string) error {
	// 检查旧密码是否匹配
	user, err := GetUserById(uid)
	if err != nil {
		return fmt.Errorf("根据uid查询用户失败, uid:[%d], err:[%w]", uid, err)
	}
	if !crypt.CheckPasswordHash(oldPassword, user.PassWord) {
		return fmt.Errorf("用户uid:[%d]旧密码错误", uid)
	}

	// 密码加密
	newPassword = crypt.BcryptPassword(newPassword)

	tx :=PostDB.Model(&model.User{}).Where("id=?", uid).
	Update("password", newPassword)
	if tx.Error != nil {
		logger.Log.Error("更新用户密码失败", "uid", uid, "error", tx.Error)
		return errors.New("更新用户密码失败, 请稍后重试")
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("用户uid:[%d]不存在", uid)
	}
	logger.Log.Info("更新用户密码成功", "uid", uid)
	return nil
}

// GetUserById 根据uid查询用户
func GetUserById(uid int) (*model.User, error) {
	user := model.User{Id: uid}
	tx := PostDB.Select("*").First(&user)
	if tx.Error != nil {	
		// Error 两种情况: 1.系统出现问题数据库连接不上
		// 2.根据条件找不到
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound)	{	
			// 如果不是 ErrRecordNotFound 根据条件查询不到结果
			// 则判断为系统出现问题, 需要打印错误日志进行记录
			logger.Log.Error("系统出现问题, 根据uid查询用户失败", "uid", uid, "err", tx.Error)
			return nil, errors.New("系统出现问题, 根据uid查询用户失败")
		}
		// 如果是 ErrRecordNotFound 根据条件查询不到结果
		// 则判断为用户不存在, 不需要打印错误日志
		return nil, fmt.Errorf("用户uid:[%d]不存在", uid)
	}
	// 如果查询到用户, 则返回用户信息
	return &user, nil
}

func GetUserByName(name string) (*model.User, error) {
	user := model.User{Name: name}
	tx := PostDB.Select("*").First(&user)
	if tx.Error != nil {	
		// Error 两种情况: 1.系统出现问题数据库连接不上
		// 2.根据条件找不到
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound)	{	
			// 如果不是 ErrRecordNotFound 根据条件查询不到结果
			// 则判断为系统出现问题, 需要打印错误日志进行记录
			logger.Log.Error("系统出现问题, 根据name查询用户失败", "name", name, "err", tx.Error)
			return nil, errors.New("系统出现问题, 根据name查询用户失败")
		}
		// 如果是 ErrRecordNotFound 根据条件查询不到结果
		// 则判断为用户不存在, 不需要打印错误日志
		return nil, fmt.Errorf("用户name:[%s]不存在", name)
	}
	// 如果查询到用户, 则返回用户信息
	return &user, nil
}

func GetIdByName(name string) (int, error) {
	user := model.User{Name: name}
	tx := PostDB.Select("id").First(&user)
	if tx.Error != nil {	 
		// Error 两种情况: 1.系统出现问题数据库连接不上
		// 2.根据条件找不到
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound)	{	
			// 如果不是 ErrRecordNotFound 根据条件查询不到结果
			// 则判断为系统出现问题, 需要打印错误日志进行记录
			logger.Log.Error("系统出现问题, 根据用户名查询用户ID失败", "name", name, "err", tx.Error)
			return 0, errors.New("系统出现问题, 根据用户名查询用户ID失败")
		}
		// 如果是 ErrRecordNotFound 根据条件查询不到结果
		// 则判断为用户不存在, 不需要打印错误日志
		return 0, fmt.Errorf("用户[%s]不存在", name)
	}
	// 如果查询到用户, 则返回用户ID
	return user.Id, nil
}
