package server

import (
	"context"
	"go.uber.org/mock/gomock"
	"grpc-crud/cmd/api"
	"grpc-crud/internal/server/mocks"
	"reflect"
	"testing"
)

//func TestNewServer(t *testing.T) {
//	type args struct {
//		cache Cache
//		db    DB
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Server
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewServer(tt.args.cache, tt.args.db); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewServer() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestServer_CreateServer(t *testing.T) {
	type fields struct {
		UnimplementedUserServiceServer api.UnimplementedUserServiceServer
		cache                          *mocks.MockCache
		db                             *mocks.MockDB
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCache(ctrl)
	mockDB := mocks.NewMockDB(ctrl)

	type args struct {
		ctx context.Context
		req *api.GetUserRequest
	}

	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *api.User
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				cache:                          tt.fields.cache,
				db:                             tt.fields.db,
			}
			got, err := s.GetUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_DeleteServer(t *testing.T) {}

func TestServer_UpdateServer(t *testing.T) {}

func TestServer_GetServer(t *testing.T) {}

func TestServer_userExists(t *testing.T) {}
