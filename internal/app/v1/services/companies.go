package services

import (
	"context"

	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/repositories"
)

type CompaniesRepository interface {
	Create(ctx context.Context, company repositories.Company) (repositories.Company, error)
	Get(ctx context.Context, id string) (repositories.Company, error)
	Update(ctx context.Context, company repositories.CompanyUpdate) error
	Delete(ctx context.Context, id string) error
}

type CompaniesService struct {
	repo CompaniesRepository
}

func NewCompaniesService(repo CompaniesRepository) *CompaniesService {
	return &CompaniesService{repo: repo}
}

func (s CompaniesService) Create(ctx context.Context, company Company) (Company, error) {
	res, err := s.repo.Create(ctx, RepositoryCompany(company))
	return CompanyFromRepository(res), handleError(err)
}

func (s CompaniesService) Get(ctx context.Context, id string) (Company, error) {
	res, err := s.repo.Get(ctx, id)
	return CompanyFromRepository(res), handleError(err)
}

func (s CompaniesService) Update(ctx context.Context, update CompanyUpdate) error {
	return handleError(s.repo.Update(ctx, RepositoryCompanyUpdate(update)))
}

func (s CompaniesService) Delete(ctx context.Context, id string) error {
	return handleError(s.repo.Delete(ctx, id))
}
