package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/repositories"
	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/services"
	"github.com/stretchr/testify/require"
)

func TestCompaniesCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mockCompaniesRepository{
			t:               t,
			expectedCompany: createTestRepoCompany(),
			returnCompany:   createTestRepoCompany(),
		}

		service := services.NewCompaniesService(repo)
		created, err := service.Create(context.Background(), createTestCompany())
		require.NoError(t, err)
		require.Equal(t, createTestCompany(), created)
	})

	t.Run("error", func(t *testing.T) {
		repo := mockCompaniesRepository{
			t:               t,
			expectedCompany: createTestRepoCompany(),
			returnError:     errors.New("error"),
		}

		service := services.NewCompaniesService(repo)
		created, err := service.Create(context.Background(), createTestCompany())
		require.ErrorAs(t, err, &services.ErrDb{})
		require.Empty(t, created)
	})
}

func TestCompaniesGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mockCompaniesRepository{
			t:             t,
			expectedId:    "id",
			returnCompany: createTestRepoCompany(),
		}

		service := services.NewCompaniesService(repo)
		company, err := service.Get(context.Background(), "id")
		require.NoError(t, err)
		require.Equal(t, createTestCompany(), company)
	})

	t.Run("not found error", func(t *testing.T) {
		repo := mockCompaniesRepository{
			t:           t,
			expectedId:  "id",
			returnError: repositories.ErrNotFound{},
		}

		service := services.NewCompaniesService(repo)
		company, err := service.Get(context.Background(), "id")
		require.ErrorAs(t, err, &services.ErrNotFound{})
		require.Empty(t, company)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mockCompaniesRepository{
			t:           t,
			expectedId:  "id",
			returnError: errors.New("error"),
		}

		service := services.NewCompaniesService(repo)
		company, err := service.Get(context.Background(), "id")
		require.ErrorAs(t, err, &services.ErrDb{})
		require.Empty(t, company)
	})
}

func TestCompaniesUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mockCompaniesRepository{
			t:                     t,
			expectedCompanyUpdate: createTestRepoCompanyUpdate(),
		}

		service := services.NewCompaniesService(repo)
		err := service.Update(context.Background(), createTestCompanyUpdate())
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		repo := mockCompaniesRepository{
			t:                     t,
			expectedCompanyUpdate: createTestRepoCompanyUpdate(),
			returnError:           errors.New("error"),
		}

		service := services.NewCompaniesService(repo)
		err := service.Update(context.Background(), createTestCompanyUpdate())
		require.ErrorAs(t, err, &services.ErrDb{})
	})
}

func TestCompaniesDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mockCompaniesRepository{
			t:          t,
			expectedId: "id",
		}

		service := services.NewCompaniesService(repo)
		err := service.Delete(context.Background(), "id")
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		repo := mockCompaniesRepository{
			t:           t,
			expectedId:  "id",
			returnError: errors.New("error"),
		}

		service := services.NewCompaniesService(repo)
		err := service.Delete(context.Background(), "id")
		require.ErrorAs(t, err, &services.ErrDb{})
	})
}

func createTestCompany() services.Company {
	return services.Company{
		ID:                "id",
		Name:              "test",
		Description:       "test description",
		AmountOfEmployees: 1,
		Registered:        true,
		Type:              "test",
	}
}

func createTestRepoCompany() repositories.Company {
	return repositories.Company{
		ID:                "id",
		Name:              "test",
		Description:       "test description",
		AmountOfEmployees: 1,
		Registered:        true,
		Type:              "test",
	}
}

func createTestCompanyUpdate() services.CompanyUpdate {
	description := "test description"
	registered := true
	return services.CompanyUpdate{
		ID:                "id",
		Name:              "test",
		Description:       &description,
		AmountOfEmployees: 1,
		Registered:        &registered,
		Type:              "test",
	}
}

func createTestRepoCompanyUpdate() repositories.CompanyUpdate {
	description := "test description"
	registered := true
	return repositories.CompanyUpdate{
		ID:                "id",
		Name:              "test",
		Description:       &description,
		AmountOfEmployees: 1,
		Registered:        &registered,
		Type:              "test",
	}
}

type mockCompaniesRepository struct {
	t                     *testing.T
	returnError           error
	returnCompany         repositories.Company
	expectedCompany       repositories.Company
	expectedCompanyUpdate repositories.CompanyUpdate
	expectedId            string
}

func (m mockCompaniesRepository) Create(ctx context.Context, company repositories.Company) (repositories.Company, error) {
	m.t.Helper()

	require.Equal(m.t, m.expectedCompany, company)
	return m.returnCompany, m.returnError
}

func (m mockCompaniesRepository) Get(ctx context.Context, id string) (repositories.Company, error) {
	m.t.Helper()
	require.Equal(m.t, m.expectedId, id)
	return m.returnCompany, m.returnError
}

func (m mockCompaniesRepository) Update(ctx context.Context, company repositories.CompanyUpdate) error {
	m.t.Helper()
	require.Equal(m.t, m.expectedCompanyUpdate, company)
	return m.returnError
}

func (m mockCompaniesRepository) Delete(ctx context.Context, id string) error {
	m.t.Helper()

	require.Equal(m.t, m.expectedId, id)
	return m.returnError
}
