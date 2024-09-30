package router

import (
	"github.com/gin-gonic/gin"
	"sduonline-training-backend/middleware"
)

func Setup(engine *gin.Engine) {
	// 测试 上线后注释掉
	/*	test := engine.Group("/test")
		{
			//测试panic
			test.GET("/panic", func(c *gin.Context) {
				panic("test panic")
			})
			//初始化数据库
			test.GET("/database_initialization", func(c *gin.Context) {
				aw := app.NewWrapper(c)
				err := model.Database_initialization()
				if err != nil {
					aw.Error(err.Error())
				}
				aw.Success("success!")
			})
		}*/

	// 用户
	user := engine.Group("/user")
	{

	}
	user.Use(middleware.JWT(1))
	{

	}

}
