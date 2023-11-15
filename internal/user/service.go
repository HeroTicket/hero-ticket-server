package user

import "context"

type Service interface {
	CreateUser(ctx context.Context, u *User) (*User, error)
	UpdateUser(ctx context.Context, u *User) (*User, error)
	DeleteUser(ctx context.Context, did string) error
	FindUsers(ctx context.Context) ([]*User, error)
	FindUserByDID(ctx context.Context, did string) (*User, error)
	FindUserByWalletAddress(ctx context.Context, walletAddress string) (*User, error)
}

type userService struct {
	repo Repository
}

func New(repo Repository) Service {
	return &userService{repo}
}

func (s *userService) CreateUser(ctx context.Context, u *User) (*User, error) {
	return s.repo.CreateUser(ctx, u)
}

func (s *userService) UpdateUser(ctx context.Context, u *User) (*User, error) {
	return s.repo.UpdateUser(ctx, u)
}

func (s *userService) DeleteUser(ctx context.Context, did string) error {
	return s.repo.DeleteUser(ctx, did)
}

func (s *userService) FindUsers(ctx context.Context) ([]*User, error) {
	return s.repo.FindUsers(ctx)
}

func (s *userService) FindUserByDID(ctx context.Context, did string) (*User, error) {
	return s.repo.FindUserByDID(ctx, did)
}

func (s *userService) FindUserByWalletAddress(ctx context.Context, walletAddress string) (*User, error) {
	return s.repo.FindUserByWalletAddress(ctx, walletAddress)
}
