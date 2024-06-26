package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-crud/cmd/api"
	"time"
)

var logger *zap.Logger

func Init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize a zap logger: %v", err))
	}
	defer logger.Sync()
}

func Get() *zap.Logger {
	return logger
}

func CreateUser(ctx context.Context, usc api.UserServiceClient, name, email string) {
	r, err := usc.CreateUser(ctx, &api.CreateUserRequest{
		Name:  name,
		Email: email,
	})
	if err != nil {
		if status.Code(err) == codes.Unknown {
			logger.Error("failed to create user", zap.Error(err))
		} else {
			logger.Fatal("failed to create user", zap.Error(err))
		}
	}
	logger.Info("user created", zap.String("email", r.GetEmail()), zap.String("id", r.GetId()))
}

func GetUser(ctx context.Context, usc api.UserServiceClient, id string) {
	r, err := usc.GetUser(ctx, &api.GetUserRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.Unknown {
			logger.Error("failed to get user", zap.Error(err))
		} else {
			logger.Fatal("failed to get user", zap.Error(err))
		}
	} else if r != nil {
		logger.Info("User retrieved", zap.String("email", r.GetEmail()), zap.String("id", r.GetId()))
	}
	logger.Info("user retrieval response is nil")
}

func UpdateUser(ctx context.Context, usc api.UserServiceClient, id, name, email string) (*api.User, error) {
	req := &api.UpdateUserRequest{
		Id:    id,
		Name:  name,
		Email: email,
	}
	updateUser, err := usc.UpdateUser(ctx, req)
	if err != nil {
		logger.Error("failed to update user", zap.Error(err))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	if updateUser == nil {
		logger.Error("failed to update user", zap.Error(err))
		return nil, fmt.Errorf("update user response is nil")
	}
	logger.Info("User update", zap.String("email", updateUser.GetEmail()), zap.String("id", updateUser.GetId()))
	//logger.Info("User updated", zap.String("email", updateUser.GetEmail()))
	return updateUser, nil
}

func DeleteUser(ctx context.Context, usc api.UserServiceClient, id string) {
	req := &api.DeleteUserRequest{Id: id}
	r, err := usc.DeleteUser(ctx, &api.DeleteUserRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.Unknown {
			logger.Error("failed to delete user", zap.Error(err))
		}
	} else if r != nil {
		logger.Info("User deleted", zap.String("email", r.GetEmail()), zap.String("id", req.GetId()))
	}
	logger.Info("User deletion response is nil")
}

func main() {
	Init()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	defer conn.Close()
	if err != nil {
		logger.Fatal("did not connect", zap.Error(err))
	}
	c := api.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//CreateUser(ctx, c, "Oleg Gorlopanov", "oleg.gorlopanov@example.com")
	//GetUser(ctx, c, "67e44047-6eef-45fd-845a-ec53abc89b55")
	_, _ = UpdateUser(ctx, c, "5339d08c-4825-4b9e-ae71-ac2f194b8a8b", "maxon4ik@example.com", "Maxim Karasik")
	//GetUser(ctx, c, "1dbf99bf-f560-45b9-895b-e5c945ad6b46")
	//DeleteUser(ctx, c, "67e44047-6eef-45fd-845a-ec53abc89b55")
	//GetUser(ctx, c, "67e44047-6eef-45fd-845a-ec53abc89b55")
	//DeleteUser(ctx, c, "1dbf99bf-f560-45b9-895b-e5c945ad6b46")
}
