package entities

// User 用于在 handler ---use case --- repository 之间传递数据的实体类
// 目的是方便 use case 操作对应的业务
type User struct {
	Uid      int64
	UserName string
	Password string
	Email    string
	Phone    string
	// AvatarURL string
}

// IsValidEmail TODO: 根据正则匹配 ？来判断是否是合法的邮箱
func (u *User) IsValidEmail() bool {
	return true
}

// EncryptPassword TODO: 加密密码, 直接修改 User 中的 Password 字段
func (u *User) EncryptPassword() error {
	return nil
}

// CheckPassword TODO: 检查密码是否正确
func (u *User) CheckPassword() bool {
	return true
}
