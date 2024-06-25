package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"grpc-crud/cmd/api"
	"time"
)

//go:generate mockgen -destination=mocks/DBMock.go -package=mocks grpc-crud/internal/server DB
type DB interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

//go:generate mockgen -destination=mocks/CacheMock.go -package=mocks grpc-crud/internal/server Cache
type Cache interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type Server struct {
	api.UnimplementedUserServiceServer
	cache Cache
	db    DB
}

func NewServer(cache Cache, db DB) *Server {
	return &Server{cache: cache, db: db}
}

func (s *Server) userExists(id string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Server) CreateServer(ctx context.Context, req *api.CreateUserRequest) (*api.User, error) {
	userID := uuid.New().String()
	_, err := s.db.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", userID, req.Name, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	createdUser := &api.User{Id: userID, Name: req.Name, Email: req.Email}
	return createdUser, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *api.UpdateUserRequest) (*api.User, error) {
	exists, err := s.userExists(req.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("user does not exist")
	}
	_, err = s.db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", req.Name, req.Email, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	updatedUser := &api.User{Id: req.Id, Name: req.Name, Email: req.Email}
	return updatedUser, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) (*api.User, error) {
	result, err := s.db.Exec("DELETE FROM users WHERE id = $1", req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("user does not exist")
	}
	s.cache.Delete(req.Id)
	return nil, nil
}

func (s *Server) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.User, error) {
	item, found := s.cache.Get(req.Id)
	if found {
		userToGet, ok := item.(*api.User)
		if !ok {
			return nil, fmt.Errorf("failed to convert to user")
		}
		return userToGet, nil
	}
	userToGet := &api.User{}
	err := s.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = $1", req.Id).Scan(&userToGet.Id, &userToGet.Name, &userToGet.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %s does not exist", req.Id)
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	s.cache.Set(req.Id, userToGet, time.Minute*5)
	return userToGet, nil
}

func (s *Server) UpdateUser(ctx context.Context)
