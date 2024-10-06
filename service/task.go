package service

import (
	"github.com/gin-gonic/gin"
	"sduonline-training-backend/model"
	"sduonline-training-backend/pkg/app"
	"strconv"
)

type TaskService struct{}

func (service *TaskService) GetAllTasks(c *gin.Context) {
	aw := app.NewWrapper(c)
	tasks, err := taskModel.GetAllTasks()
	if err != nil {
		aw.Error(err.Error())
		return
	}
	aw.Success(tasks)
}
func (service *TaskService) AddTask(c *gin.Context) {
	aw := app.NewWrapper(c)
	var req model.Task
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	if err := taskModel.AddTask(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	aw.OK()
}
func (service *TaskService) DeleteTask(c *gin.Context) {
	aw := app.NewWrapper(c)
	taskID, err := strconv.Atoi(c.Param("tid"))
	if err != nil {
		aw.Error(err.Error())
		return
	}
	if !taskModel.CheckTaskByID(taskID) {
		aw.Error("培训不存在", 404)
		return
	}
	if err := taskModel.DeleteTask(taskID); err != nil {
		aw.Error(err.Error())
		return
	}
	aw.OK()
}
func (service *TaskService) UpdateTask(c *gin.Context) {
	aw := app.NewWrapper(c)
	taskID, err := strconv.Atoi(c.Param("tid"))
	if err != nil {
		aw.Error(err.Error())
		return
	}
	if !taskModel.CheckTaskByID(taskID) {
		aw.Error("培训不存在", 404)
		return
	}
	var req model.Task
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	req.ID = taskID
	if err := taskModel.UpdateTask(&req); err != nil {
		aw.Error(err.Error())
		return
	}
	aw.OK()
}
