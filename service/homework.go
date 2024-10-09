package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sduonline-training-backend/model"
	"sduonline-training-backend/pkg/app"
	"sduonline-training-backend/pkg/conf"
	"sduonline-training-backend/pkg/util"
	"strings"
	"time"
)

type HomeworkService struct{}
type HomeworkSimpleVO struct {
	HomeworkID  string    `json:"hid"`
	StudentID   string    `json:"sid"`
	TaskID      int       `json:"tid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Display     bool      `json:"display"`
	CreatedAt   time.Time `json:"create_at"`
}
type HomeworkVO struct {
	HomeworkID  string    `json:"hid"`
	StudentID   string    `json:"sid"`
	TaskID      int       `json:"tid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Note        *string   `json:"note,omitempty"`
	Display     bool      `json:"display"`
	CreatedAt   time.Time `json:"create_at"`
}
type HomeworkDTO struct {
	StudentID     string                  `form:"-"`
	TaskID        int                     `form:"tid" binding:"required"`
	Title         string                  `form:"title" binding:"required"`
	Description   string                  `form:"description" binding:"required"`
	Note          *string                 `form:"note"`
	Files         []*multipart.FileHeader `form:"files" binding:"required"`
	RelativePaths []string                `form:"relative_paths" binding:"required"`
}

func (service HomeworkService) DisplayHomework(c *gin.Context) {
	aw := app.NewWrapper(c)
	if err := homeworkModel.DisplayHomework(c.Param("hid")); err != nil {
		aw.Error(err.Error())
		return
	}
	aw.OK()
}
func (service *HomeworkService) UploadHomework(c *gin.Context) {
	aw := app.NewWrapper(c)
	uc := util.ExtractUserClaims(c)

	var req HomeworkDTO
	if err := c.ShouldBind(&req); err != nil {
		aw.Error(err.Error())
		return
	}

	if homeworkModel.CheckHomeworkExistence(uc.UserID, req.TaskID) {
		aw.Error("不可重复上传", 403)
		return
	}

	// 验证文件和路径数量一致
	if len(req.Files) != len(req.RelativePaths) {
		aw.Error("files and relative_paths count mismatch")
		return
	}

	// 获取学生ID并验证用户权限
	user := userModel.FindUserByID(uc.UserID)
	req.StudentID = user.StudentID

	homework := model.Homework{
		StudentID:   req.StudentID,
		TaskID:      req.TaskID,
		Title:       req.Title,
		Description: req.Description,
		Note:        req.Note,
		Display:     false,
	}

	// 使用事务确保一致性
	tx := homeworkModel.Tx.Begin()
	homework.HomeworkID = uuid.New().String()
	if homework.Note == nil {
		note := ""
		homework.Note = &note
	}
	if err := tx.Create(&homework).Error; err != nil {
		aw.Error(err.Error())
		tx.Rollback()
		return
	}

	for i, file := range req.Files {
		relativePath := filepath.Clean(req.RelativePaths[i])
		if strings.Contains(relativePath, "..") {
			aw.Error("invalid file path")
			DeleteHomeworkFile(homework.HomeworkID)
			tx.Rollback()
			return
		}

		savePath := filepath.Join(conf.Conf.UploadDir, homework.HomeworkID, relativePath)
		if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
			aw.Error(err.Error())
			DeleteHomeworkFile(homework.HomeworkID)
			tx.Rollback()
			return
		}

		uploadedFile, err := file.Open()
		if err != nil {
			aw.Error(err.Error())
			DeleteHomeworkFile(homework.HomeworkID)
			tx.Rollback()
			return
		}
		defer uploadedFile.Close()

		// 验证文件大小和类型
		const MaxFileSize = 10 << 21 // 10 MB
		if file.Size > MaxFileSize {
			aw.Error("file size exceeds limit")
			DeleteHomeworkFile(homework.HomeworkID)
			tx.Rollback()
			return
		}

		/*allowedTypes := []string{"application/pdf", "image/jpeg", "image/png"}
		fileType := file.Header.Get("Content-Type")
		if !contains(allowedTypes, fileType) {
			aw.Error("unsupported file type")
			tx.Rollback()
			return
		}*/

		fileData, err := io.ReadAll(uploadedFile)
		if err != nil {
			aw.Error(err.Error())
			DeleteHomeworkFile(homework.HomeworkID)
			tx.Rollback()
			return
		}

		if err := os.WriteFile(savePath, fileData, 0644); err != nil {
			aw.Error(err.Error())
			DeleteHomeworkFile(homework.HomeworkID)
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	aw.OK()
}
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
func DeleteHomeworkFile(homeworkID string) error {
	path := filepath.Join(conf.Conf.UploadDir, "/", homeworkID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return nil
}
func DeleteHomeworkByHomeworkID(homeworkID string) error {
	homework := model.Homework{
		HomeworkID: homeworkID,
	}
	if err := homeworkModel.FindHomeworkByHomeworkID(&homework); err != nil {
		return err
	}
	tx := homeworkModel.Tx.Begin()
	if err := tx.Where("homework_id = ?", homework.HomeworkID).Delete(&model.Homework{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除作业文件
	if err := DeleteHomeworkFile(homeworkID); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (service HomeworkService) DeleteHomework(c *gin.Context) {
	aw := app.NewWrapper(c)
	homeworkID := c.Param("hid")
	if err := DeleteHomeworkByHomeworkID(homeworkID); err != nil {
		aw.Error(err.Error())
		return
	}
	aw.OK()
}
