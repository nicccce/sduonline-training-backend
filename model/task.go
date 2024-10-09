package model

import "gorm.io/gorm"

type Task struct {
	ID       int    `json:"tid" binding:"-"`
	Name     string `json:"tname" binding:"required"`
	Content  string `json:"content" binding:"required"`
	DeadLine string `json:"ddl" binding:"required"`
}
type TaskModel struct {
	AbstractModel
}

func (receiver TaskModel) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := receiver.Tx.Find(&tasks).Error
	if err == gorm.ErrRecordNotFound {
		return tasks, nil
	}
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
func (receiver TaskModel) AddTask(task *Task) error {
	return receiver.Tx.Create(task).Error
}
func (receiver TaskModel) CheckTaskByID(taskID int) bool {
	var task Task
	if err := receiver.Tx.First(&task, taskID).Error; err != nil {
		if err != nil {
			return false
		}
	}
	return true
}
func (receiver TaskModel) UpdateTask(task *Task) error {
	return receiver.Tx.Save(task).Error
}
func (receiver TaskModel) DeleteTask(id int) error {
	return receiver.Tx.Where("id = ?", id).Delete(&Task{}).Error
}
