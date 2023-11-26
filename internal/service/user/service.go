package user

import "context"

type Service interface {
	CreateUser(ctx context.Context, params CreateUserParams) (*User, error)
	UpdateUser(ctx context.Context, params UpdateUserParams) error
	DeleteUser(ctx context.Context, id string) error
	FindAdmin(ctx context.Context) (*User, error)
	FindUsers(ctx context.Context) ([]*User, error)
	FindUserByID(ctx context.Context, id string) (*User, error)
	FindUserByAccountAddress(ctx context.Context, accountAddress string) (*User, error)
	FindUserByName(ctx context.Context, name string) (*User, error)
}

type userService struct {
	repo Repository
}

func New(repo Repository) Service {
	return &userService{repo}
}

func (s *userService) CreateUser(ctx context.Context, params CreateUserParams) (*User, error) {
	return s.repo.CreateUser(ctx, params)
}

func (s *userService) UpdateUser(ctx context.Context, params UpdateUserParams) error {
	return s.repo.UpdateUser(ctx, params)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *userService) FindAdmin(ctx context.Context) (*User, error) {
	return s.repo.FindAdmin(ctx)
}

func (s *userService) FindUsers(ctx context.Context) ([]*User, error) {
	return s.repo.FindUsers(ctx)
}

func (s *userService) FindUserByID(ctx context.Context, id string) (*User, error) {
	return s.repo.FindUserByID(ctx, id)
}

func (s *userService) FindUserByAccountAddress(ctx context.Context, accountAddress string) (*User, error) {
	return s.repo.FindUserByAccountAddress(ctx, accountAddress)
}

func (s *userService) FindUserByName(ctx context.Context, name string) (*User, error) {
	return s.repo.FindUserByName(ctx, name)
}
