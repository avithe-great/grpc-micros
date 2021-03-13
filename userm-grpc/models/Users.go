package models

import (
	"userm-grpc/proto"
)

type Users struct {
	Id      int64 `gorm:"type:int;primary_key"`
	Fname   string
	City    string
	Phone   int64
	Height  int64
	Married bool
}

func (em *Users) GetgRPCModel() proto.User {
	return proto.User{
		Id:      em.Id,
		Fname:   em.Fname,
		City:    em.City,
		Phone:   em.Phone,
		Height:  em.Height,
		Married: em.Married,
	}
}

func (em *Users) From(gRPCModel proto.User) {

	em.Id = gRPCModel.Id
	em.Fname = gRPCModel.Fname
	em.Phone = gRPCModel.Phone
	em.City = gRPCModel.City
	em.Height = gRPCModel.Height
	em.Married = gRPCModel.Married
}
