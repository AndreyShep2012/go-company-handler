package handlers

import "github.com/AndreyShep2012/go-company-handler/internal/app/v1/services"

type CreateCompanyRequest struct {
	Name              string `json:"name" validate:"required,max=15"`
	Description       string `json:"description" validate:"omitempty,max=3000"`
	AmountOfEmployees int    `json:"amount_of_employees" validate:"required,gte=0"`
	Registered        *bool  `json:"registered" validate:"required"`
	Type              string `json:"type" validate:"required,oneof=Corporations NonProfit Cooperative 'Sole Proprietorship'"`
}

type UpdateCompanyRequest struct {
	Name              string  `json:"name" validate:"required,max=15"`
	Description       *string `json:"description" validate:"omitempty,max=3000"`
	AmountOfEmployees int     `json:"amount_of_employees" validate:"required,gte=0"`
	Registered        *bool   `json:"registered" validate:"required"`
	Type              string  `json:"type" validate:"required,oneof=Corporations NonProfit Cooperative 'Sole Proprietorship'"`
}

type Company struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description,omitempty"`
	AmountOfEmployees int    `json:"amount_of_employees"`
	Registered        bool   `json:"registered"`
	Type              string `json:"type"`
}

func CompanyFromService(company services.Company) Company {
	return Company{
		ID:                company.ID,
		Name:              company.Name,
		Description:       company.Description,
		AmountOfEmployees: company.AmountOfEmployees,
		Registered:        company.Registered,
		Type:              company.Type,
	}
}

func CompanyUpdateToService(id string, req UpdateCompanyRequest) services.CompanyUpdate {
	return services.CompanyUpdate{
		ID:                id,
		Name:              req.Name,
		Description:       req.Description,
		AmountOfEmployees: req.AmountOfEmployees,
		Registered:        req.Registered,
		Type:              req.Type,
	}
}
