package services

import "github.com/AndreyShep2012/go-company-handler/internal/app/v1/repositories"

type Company struct {
	ID                string
	Name              string
	Description       string
	AmountOfEmployees int
	Registered        bool
	Type              string
}

type CompanyUpdate struct {
	ID                string
	Name              string
	Description       *string
	AmountOfEmployees int
	Registered        *bool
	Type              string
}

func CompanyFromRepository(company repositories.Company) Company {
	return Company{
		ID:                company.ID,
		Name:              company.Name,
		Description:       company.Description,
		AmountOfEmployees: company.AmountOfEmployees,
		Registered:        company.Registered,
		Type:              company.Type,
	}
}

func RepositoryCompany(company Company) repositories.Company {
	return repositories.Company{
		ID:                company.ID,
		Name:              company.Name,
		Description:       company.Description,
		AmountOfEmployees: company.AmountOfEmployees,
		Registered:        company.Registered,
		Type:              company.Type,
	}
}

func RepositoryCompanyUpdate(company CompanyUpdate) repositories.CompanyUpdate {
	return repositories.CompanyUpdate{
		ID:                company.ID,
		Name:              company.Name,
		Description:       company.Description,
		AmountOfEmployees: company.AmountOfEmployees,
		Registered:        company.Registered,
		Type:              company.Type,
	}
}
