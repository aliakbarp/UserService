// This file contains types that are used in the repository layer.
package repository

type InsertUsersDataInput struct {
	PhoneNumber    string
	FullName       string
	HashedPassword string
}

type InsertUsersDataOutput struct {
	Id int32
}

type GetUserByPhoneNumberInput struct {
	PhoneNumber string
}

type GetUserByPhoneNumberOutput struct {
	Id          int32
	FullName    string
	PhoneNumber string
	HashedPass  string
	Count       int32
}

type GetUserByIdInput struct {
	Id int32
}

type GetUserByIdOutput struct {
	Id          int32
	FullName    string
	PhoneNumber string
	HashedPass  string
	Count       int32
}

type UpdateFullNameByIdInput struct {
	Id       int32
	FullName string
}

type UpdateCountByIdInput struct {
	Id    int32
	Count int32
}
