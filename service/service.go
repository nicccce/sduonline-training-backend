package service

import "sduonline-training-backend/model"

var userModel model.UserModel

func Setup() {
	userModel = model.UserModel{AbstractModel: model.AbstractModel{Tx: model.DB}}

}
