package model

import (
	"gorm.io/gorm"
	"sduonline-training-backend/pkg/util"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	RoleID   int    `json:"role_id"`
	Password string `json:"password"`
	UserInfo
}
type UserInfo struct {
	UserName  string    `json:"username" form:"username" binding:"required"`
	StudentID string    `json:"sid" form:"sid" binding:"required"`
	Qq        string    `json:"qq" form:"qq" binding:"required"`
	Info      string    `json:"info" form:"info" binding:"required"`
	CreatedAt time.Time `json:"created_at" binging:"-"`
}

type UserModel struct {
	AbstractModel
}

func (receiver UserModel) FindUserByID(id int) *User {
	var user User
	err := receiver.Tx.Take(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	util.ForwardOrPanic(err)
	return &user
}
func (receiver UserModel) UpdateUser(user *User) {
	err := receiver.Tx.Save(user).Error
	util.ForwardOrPanic(err)
}
func (receiver UserModel) CreateUser(user *User) {
	err := receiver.Tx.Create(user).Error
	util.ForwardOrPanic(err)
}
func (receiver UserModel) FindAllUsers() []User {
	var users []User
	err := receiver.Tx.Order("id desc").Find(&users).Error
	util.ForwardOrPanic(err)
	return users
}
func (receiver UserModel) DeleteUser(user *User) {
	err := receiver.Tx.Unscoped().Delete(user).Error
	util.ForwardOrPanic(err)
}
func (receiver UserModel) FindUserByStudentID(studentID string) *User {
	var user User
	err := receiver.Tx.Where("student_id=?", studentID).Take(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	util.ForwardOrPanic(err)
	return &user
}
