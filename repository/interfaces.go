// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	InsertUsersData(ctx context.Context, input InsertUsersDataInput) (output InsertUsersDataOutput, err error)
	GetUserByPhoneNumber(ctx context.Context, input GetUserByPhoneNumberInput) (output GetUserByPhoneNumberOutput, err error)
	GetUserById(ctx context.Context, input GetUserByIdInput) (output GetUserByIdOutput, err error)
	UpdateFullNameByIdInput(ctx context.Context, input UpdateFullNameByIdInput) (err error)
	UpdateCountById(ctx context.Context, input UpdateCountByIdInput) (err error)
}
