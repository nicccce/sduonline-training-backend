package service

import "sduonline-training-backend/model"

var userModel model.UserModel
var taskModel model.TaskModel

func Setup() {
	userModel = model.UserModel{AbstractModel: model.AbstractModel{Tx: model.DB}}
	taskModel = model.TaskModel{AbstractModel: model.AbstractModel{Tx: model.DB}}
}
