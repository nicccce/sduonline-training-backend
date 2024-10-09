package model

import (
	"github.com/google/uuid"
	"time"
)

type Homework struct {
	ID          int       `json:"id"`
	HomeworkID  string    `json:"hid"`
	StudentID   string    `json:"sid"`
	TaskID      int       `json:"tid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Note        *string   `json:"note,omitempty"`
	Display     bool      `json:"display"`
	CreatedAt   time.Time `json:"created_at"`
}
type HomeworkModel struct {
	AbstractModel
}

func (receiver HomeworkModel) CreateHomework(homework *Homework) error {
	homework.HomeworkID = uuid.New().String()
	if homework.Note == nil {
		note := ""
		homework.Note = &note
	}
	result := receiver.Tx.Create(homework)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (receiver HomeworkModel) GetHomeworks(taskID int) ([]Homework, error) {
	var homeworks []Homework
	result := receiver.Tx.Where("task_id = ? AND display = true", taskID).Find(&homeworks)
	if result.Error != nil {
		return nil, result.Error
	}
	return homeworks, nil
}
func (receiver HomeworkModel) GetHomeworksAdmin(taskID int) ([]Homework, error) {
	var homeworks []Homework
	result := receiver.Tx.Where("task_id = ?", taskID).Find(&homeworks)
	if result.Error != nil {
		return nil, result.Error
	}
	return homeworks, nil
}
func (receiver HomeworkModel) GetHomeworksByStudentID(studentID string) ([]Homework, error) {
	var homeworks []Homework
	result := receiver.Tx.Where("student_id = ?", studentID).Find(&homeworks)
	if result.Error != nil {
		return nil, result.Error
	}
	return homeworks, nil
}
func (receiver HomeworkModel) DisplayHomework(homeworkID string) error {
	result := receiver.Tx.Model(&Homework{}).Where("homework_id = ?", homeworkID).Update("display", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (receiver HomeworkModel) FindHomeworkByHomeworkID(homework *Homework) error {
	result := receiver.Tx.Where("homework_id = ?", homework.HomeworkID).First(&homework)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (receiver HomeworkModel) GetHomeworkIDsByTaskID(taskID int) ([]string, error) {
	var homeworkList []Homework
	var homeworkIDs []string

	// 查询符合条件的作业记录
	if err := receiver.Tx.Where("task_id = ?", taskID).Find(&homeworkList).Error; err != nil {
		return nil, err
	}

	// 提取作业记录中的 HomeworkID 到列表中
	for _, homework := range homeworkList {
		homeworkIDs = append(homeworkIDs, homework.HomeworkID)
	}

	return homeworkIDs, nil
}
func (receiver HomeworkModel) GetHomeworkIDsByStudentID(studentID string) ([]string, error) {
	var homeworkList []Homework
	var homeworkIDs []string

	// 查询符合条件的作业记录
	if err := receiver.Tx.Where("student_id = ?", studentID).Find(&homeworkList).Error; err != nil {
		return nil, err
	}

	// 提取作业记录中的 HomeworkID 到列表中
	for _, homework := range homeworkList {
		homeworkIDs = append(homeworkIDs, homework.HomeworkID)
	}

	return homeworkIDs, nil
}
func (receiver HomeworkModel) CheckHomeworkExistence(userID int, taskID int) bool {
	var count int64
	var studentID string
	receiver.Tx.Table("users").Select("student_id").Where("id = ?", userID).Pluck("student_id", &studentID)
	receiver.Tx.Model(&Homework{}).Where("student_id = ? AND task_id = ?", studentID, taskID).Count(&count)
	return count > 0
}
