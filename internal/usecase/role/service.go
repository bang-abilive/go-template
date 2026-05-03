package role

import (
	"context"
	"ndinhbang/go-template/internal/domain/entity"
	"ndinhbang/go-template/internal/domain/values"
)

type Service interface {
	Create(ctx context.Context, in CreateRequest) (CreateResponse, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, in CreateRequest) (CreateResponse, error) {
	// 1. Validate & Khởi tạo Value Object
	slugVO, err := values.NewSlug(in.Slug)
	if err != nil {
		return CreateResponse{}, err
	}

	// 3. Khởi tạo Entity
	r := &entity.Role{
		Slug:        slugVO,
		Name:        in.Name,
		Lv:          in.Lv,
		Permissions: in.Permissions,
	}

	// 4. Lưu trữ
	if err := s.repository.Create(ctx, r); err != nil {
		return CreateResponse{}, err
	}

	return CreateResponse{
		ID:        r.ID,
		Slug:      r.Slug.Value(),
		CreatedAt: r.CreatedAt,
	}, nil
}
