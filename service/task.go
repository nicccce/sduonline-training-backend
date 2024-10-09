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
	studentIDs, err := homeworkModel.GetHomeworkIDsByTaskID(taskID)
	if err != nil {
		aw.Error(err.Error())
	}
	for _, studentID := range studentIDs {
		err = DeleteHomeworkByHomeworkID(studentID)
	}
	if err != nil {
		aw.Error(err.Error())
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
func (service TaskService) GetHomeworks(c *gin.Context) {
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
	homeworks, err := homeworkModel.GetHomeworks(taskID)
	if err != nil {
		aw.Error(err.Error())
		return
	}
	var homeworkDTOs []HomeworkSimpleVO
	for _, homework := range homeworks {
		homeworkDTO := HomeworkSimpleVO{
			HomeworkID:  homework.HomeworkID,
			StudentID:   homework.StudentID,
			TaskID:      homework.TaskID,
			Title:       homework.Title,
			Description: homework.Description,
			Display:     homework.Display,
			CreatedAt:   homework.CreatedAt,
		}

		homeworkDTOs = append(homeworkDTOs, homeworkDTO)
	}
	aw.Success(homeworkDTOs)
}
func (service TaskService) GetHomeworksAdmin(c *gin.Context) {
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
	homeworks, err := homeworkModel.GetHomeworksAdmin(taskID)
	if err != nil {
		aw.Error(err.Error())
		return
	}
	var homeworkDTOs []HomeworkVO
	for _, homework := range homeworks {
		homeworkDTO := HomeworkVO{
			HomeworkID:  homework.HomeworkID,
			StudentID:   homework.StudentID,
			TaskID:      homework.TaskID,
			Title:       homework.Title,
			Description: homework.Description,
			Display:     homework.Display,
			CreatedAt:   homework.CreatedAt,
			Note:        homework.Note,
		}

		homeworkDTOs = append(homeworkDTOs, homeworkDTO)
	}
	aw.Success(homeworkDTOs)
}
