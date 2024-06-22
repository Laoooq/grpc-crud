package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-crud/cmd/api"
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

func GetUser(ctx context.Context, usc *api.UserServiceClient, id string) {}

func UpdateUser(ctx context.Context, usc *api.UserServiceClient, name, email string) {}

func DeleteUser(ctx context.Context, usc *api.UserServiceClient, name string) {}
