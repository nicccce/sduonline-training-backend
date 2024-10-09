package router

import (
	"github.com/gin-gonic/gin"
	"sduonline-training-backend/middleware"
	"sduonline-training-backend/service"
)

func Setup(engine *gin.Engine) {
	// 测试 上线后注释掉
	/*test := engine.Group("/test")
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
	user := engine.Group("/users")
	{
		hub := service.UserService{}
		user.GET("/test_get_jwt", hub.TestGetJWT)
		user.POST("/login", hub.Login)
		user.POST("/register", hub.Register)
	}
	user.Use(middleware.JWT(1))
	{
		hub := service.UserService{}
		user.DELETE("/:sid", hub.DeleteUser)
		user.GET("/me", hub.GetUser)
	}
	user.Use(middleware.JWT(2))
	{
		hub := service.UserService{}
		user.GET("/", hub.GetUsers)
	}
	task := engine.Group("/tasks")
	task.Use(middleware.JWT(1))
	{
		hub := service.TaskService{}
		task.GET("/", hub.GetAllTasks)
		task.GET("/:tid/homeworks", hub.GetHomeworks)
	}
	task.Use(middleware.JWT(2))
	{
		hub := service.TaskService{}
		task.POST("/", hub.AddTask)
		task.DELETE("/:tid", hub.DeleteTask)
		task.POST("/:tid", hub.UpdateTask)
		task.GET("/:tid/homeworks/admin", hub.GetHomeworksAdmin)
	}
	homework := engine.Group("/homeworks")
	homework.Use(middleware.JWT(1))
	{
		hub := service.HomeworkService{}
		homework.DELETE("/:hid", hub.DeleteHomework)
		homework.POST("/:hid/public", hub.DisplayHomework)
		homework.POST("/", hub.UploadHomework)
	}
}
