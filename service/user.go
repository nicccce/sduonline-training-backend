package service

import (
	"github.com/gin-gonic/gin"
	"sduonline-training-backend/pkg/app"
	"sduonline-training-backend/pkg/conf"
	"sduonline-training-backend/pkg/util"
)

type UserService struct {
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
	user, _ := userModel.FindUserByID(req.UserID)
	if user == nil {
		aw.Error("userID不存在")
		return
	}
	aw.Success(util.GenerateJWT(user.ID, user.RoleID))
}
