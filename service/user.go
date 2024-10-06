package service

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"sduonline-training-backend/model"
	"sduonline-training-backend/pkg/app"
	"sduonline-training-backend/pkg/conf"
	"sduonline-training-backend/pkg/util"
)

type UserService struct {
}
type UserVO struct {
	ID     int `json:"id"`
	RoleID int `json:"role_id"`
	model.UserInfo
	Token *string `json:"token,omitempty"`
}

func (receiver UserService) TestGetJWT(c *gin.Context) {
	aw := app.NewWrapper(c)
	type getJWTReq struct {
		UserID    int    `form:"user_id" binding:"required"`
		JWTSecret string `form:"jwt_secret" binding:"required"`
	}
	var req getJWTReq
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	if req.JWTSecret != conf.Conf.JWTSecret {
		aw.Error("jwtSecret不正确")
		return
	}
	user := userModel.FindUserByID(req.UserID)
	if user == nil {
		aw.Error("userID不存在")
		return
	}
	aw.Success(util.GenerateJWT(user.ID, user.RoleID))
}
func (receiver UserService) Login(c *gin.Context) {
	aw := app.NewWrapper(c)
	type loginReq struct {
		StudentID string `form:"sid" json:"sid" binding:"required"`
		Password  string `form:"password" json:"password" binding:"required"`
	}
	var req loginReq
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	user := userModel.FindUserByStudentID(req.StudentID)
	if user == nil {
		aw.Error("用户名不存在", 404)
		return
	}
	if !util.CheckPassword(req.Password, user.Password) {
		aw.Error("用户名或密码错误", 401)
		return
	}
	token := util.GenerateJWT(user.ID, user.RoleID)
	userVO := UserVO{
		ID:       user.ID,
		RoleID:   user.RoleID,
		UserInfo: user.UserInfo,
		Token:    &token,
	}
	aw.Success(userVO)
}
func (receiver UserService) Register(c *gin.Context) {
	aw := app.NewWrapper(c)
	type registerReq struct {
		Password string `form:"password" json:"password" binding:"required"`
		model.UserInfo
	}
	var req registerReq
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	pattern := `^202\d{9}$`
	rgxp := regexp.MustCompile(pattern)
	if !rgxp.MatchString(req.StudentID) {
		aw.Error("学号不合法", 400)
		return
	}
	user := userModel.FindUserByStudentID(req.StudentID)
	if user != nil {
		aw.Error("账号已被注册", 400)
		return
	}
	encryptPassword, err := util.EncryptPassword(req.Password)
	if err != nil {
		aw.Error(err.Error())
		return
	}
	user = &model.User{
		RoleID:   1,
		Password: encryptPassword,
		UserInfo: req.UserInfo,
	}
	userModel.CreateUser(user)
	aw.OK()
}
func (receiver UserService) DeleteUser(c *gin.Context) {
	aw := app.NewWrapper(c)
	studentID := c.Param("sid")
	user := userModel.FindUserByStudentID(studentID)
	if user == nil {
		aw.Error("账号不存在", 403)
		return
	}
	userClaim := util.ExtractUserClaims(c)
	if userClaim.RoleID < 1 && user.StudentID != studentID {
		aw.Error("无权限", 403)
		return
	}
	userModel.DeleteUser(user)
	aw.OK()
}
