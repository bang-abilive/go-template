package user

import (
	"context"
	"errors"
	"ndinhbang/go-template/internal/domain/entity"
	"ndinhbang/go-template/internal/domain/values"
)

type Service interface {
	Register(ctx context.Context, in RegisterRequest) (RegisterResponse, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Register(ctx context.Context, in RegisterRequest) (RegisterResponse, error) {
	// 1. Validate & Khởi tạo Value Object
	emailVO, err := values.NewEmail(in.Email)
	if err != nil {
		return RegisterResponse{}, err
	}

	// 2. Kiểm tra nghiệp vụ
	user, err := s.repository.FindByEmail(ctx, emailVO.Value())
	if err != nil {
		return RegisterResponse{}, err
	}
	if user != nil {
		return RegisterResponse{}, errors.New("user already exists")
	}

	// 3. Khởi tạo Entity
	u := &entity.User{
		Email:    emailVO,
		Password: in.Password, // Thực tế cần hash password ở đây
	}

	// 4. Lưu trữ
	if err := s.repository.Create(ctx, u); err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{
		ID:        u.ID,
		Email:     u.Email.Value(),
		CreatedAt: u.CreatedAt,
	}, nil
}
