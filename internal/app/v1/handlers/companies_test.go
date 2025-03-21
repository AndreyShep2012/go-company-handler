package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/handlers"
	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestCreateCompany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedResponse := handlers.Company{
			ID:                "605c72efb1e2c3d1f8a1b2c3",
			Name:              "name",
			Description:       "description",
			AmountOfEmployees: 10,
			Registered:        true,
			Type:              "Sole Proprietorship",
		}

		fiberApp := initFiberApp()

		ch := make(chan any)
		handlers.SetupCompaniesRoutes(fiberApp, mockCompaniesService{
			t: t,
			expectedCompany: services.Company{
				Name:              "name",
				Description:       "description",
				AmountOfEmployees: 100,
				Registered:        true,
				Type:              "Sole Proprietorship",
			},
			returnCompany: services.Company{
				ID:                "605c72efb1e2c3d1f8a1b2c3",
				Name:              "name",
				Description:       "description",
				AmountOfEmployees: 10,
				Registered:        true,
				Type:              "Sole Proprietorship",
			},
		}, newMockPublisher(ch))
		defer close(ch)

		body := `{
			"name":"name",
			"description":"description",
			"amount_of_employees":100,
			"registered":true,
			"type":"Sole Proprietorship"
		}`

		req := httptest.NewRequest("POST", "/companies/create", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		response, err := fiberApp.Test(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, fiber.StatusOK, response.StatusCode)

		defer response.Body.Close()
		bodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		var res handlers.Company
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)

		require.Equal(t, expectedResponse, res)

		select {
		case e := <-ch:
			require.Equal(t, expectedResponse, e)
		case <-time.After(time.Millisecond * 500):
			require.Fail(t, "timeout")
		}
	})

	t.Run("bad request error", func(t *testing.T) {
		fiberApp := initFiberApp()

		handlers.SetupCompaniesRoutes(fiberApp, nil, nil)

		doTest := func(body string) {
			req := httptest.NewRequest("POST", "/companies/create", bytes.NewReader([]byte(body)))
			req.Header.Set("Content-Type", "application/json")

			response, err := fiberApp.Test(req)
			require.NoError(t, err)
			require.NotNil(t, response)
			defer response.Body.Close()
			require.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		}

		doTest(`{}`)
		doTest(`not a json`)
	})

	t.Run("internal server error", func(t *testing.T) {
		fiberApp := initFiberApp()

		handlers.SetupCompaniesRoutes(fiberApp, mockCompaniesService{
			t: t,
			expectedCompany: services.Company{
				Name:              "name",
				Description:       "description",
				AmountOfEmployees: 100,
				Registered:        true,
				Type:              "Sole Proprietorship",
			},
			returnError: services.ErrDb{},
		}, nil)
		body := `{
			"name":"name",
			"description":"description",
			"amount_of_employees":100,
			"registered":true,
			"type":"Sole Proprietorship"
		}`

		req := httptest.NewRequest("POST", "/companies/create", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		response, err := fiberApp.Test(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		defer response.Body.Close()
		require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	})
}

func TestGetCompany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedResponse := handlers.Company{
			ID:                "605c72efb1e2c3d1f8a1b2c3",
			Name:              "name",
			Description:       "description",
			AmountOfEmployees: 10,
			Registered:        true,
			Type:              "Sole Proprietorship",
		}

		fiberApp := initFiberApp()

		handlers.SetupCompaniesRoutes(fiberApp, mockCompaniesService{
			t: t,
			returnCompany: services.Company{
				ID:                "605c72efb1e2c3d1f8a1b2c3",
				Name:              "name",
				Description:       "description",
				AmountOfEmployees: 10,
				Registered:        true,
				Type:              "Sole Proprietorship",
			},
			expectedId: "605c72efb1e2c3d1f8a1b2c3",
		}, nil)

		req := httptest.NewRequest("GET", "/companies/605c72efb1e2c3d1f8a1b2c3", nil)

		response, err := fiberApp.Test(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, fiber.StatusOK, response.StatusCode)

		defer response.Body.Close()
		bodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		var res handlers.Company
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)

		require.Equal(t, expectedResponse, res)
	})

	t.Run("bad request error", func(t *testing.T) {
		fiberApp := initFiberApp()

		handlers.SetupCompaniesRoutes(fiberApp, nil, nil)

		req := httptest.NewRequest("GET", "/companies/wrong_id", nil)

		response, err := fiberApp.Test(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		defer response.Body.Close()
		require.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	})

	t.Run("internal server error", func(t *testing.T) {
		fiberApp := initFiberApp()

		handlers.SetupCompaniesRoutes(fiberApp, mockCompaniesService{
			t:           t,
			returnError: services.ErrDb{},
			expectedId:  "605c72efb1e2c3d1f8a1b2c3",
		}, nil)

		req := httptest.NewRequest("GET", "/companies/605c72efb1e2c3d1f8a1b2c3", nil)

		response, err := fiberApp.Test(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		defer response.Body.Close()
		require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	})
}

func TestUpdateCompany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fiberApp := initFiberApp()

		description := "description"
		registered := true
		expectedCompanyUpdate := services.CompanyUpdate{
			ID:                "605c72efb1e2c3d1f8a1b2c3",
			Name:              "name",
			Description:       &description,
			AmountOfEmployees: 100,
			Registered:        &registered,
			Type:              "Sole Proprietorship",
		}
		ch := make(chan any)
		handlers.SetupCompaniesRoutes(fiberApp, mockCompaniesService{
			t:                     t,
			expectedCompanyUpdate: expectedCompanyUpdate,
		}, newMockPublisher(ch))

		body := `{
			"name":"name",
			"description":"description",
			"amount_of_employees":100,
			"registered":true,
			"type":"Sole Proprietorship"
		}`

		req := httptest.NewRequest("PATCH", "/companies/605c72efb1e2c3d1f8a1b2c3", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		response, err := fiberApp.Test(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		defer response.Body.Close()
		require.Equal(t, fiber.StatusNoContent, response.StatusCode)

		select {
		case e := <-ch:
			require.Equal(t, expectedCompanyUpdate, e)
		case <-time.After(time.Millisecond * 500):
			require.Fail(t, "timeout")
		}
	})

	t.Run("bad request error", func(t *testing.T) {
		fiberApp := initFiberApp()

		handlers.SetupCompaniesRoutes(fiberApp, nil, nil)

		doTest := func(id, body string) {
			req := httptest.NewRequest("PATCH", "/companies/"+id, bytes.NewReader([]byte(body)))
			req.Header.Set("Content-Type", "application/json")

			response, err := fiberApp.Test(req)
			require.NoError(t, err)
			require.NotNil(t, response)
			defer response.Body.Close()
			require.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		}

		doTest("605c72efb1e2c3d1f8a1b2c3", `{}`)
		doTest("wrong_id", `{
			"name":"name",
			"description":"description",
			"amount_of_employees":100,
			"registered":true,
			"type":"Sole Proprietorship"
		}`)
		doTest("605c72efb1e2c3d1f8a1b2c3", `not a json`)
	})

	t.Run("internal server error", func(t *testing.T) {
		fiberApp := initFiberApp()

		description := "description"
		registered := true
		handlers.SetupCompaniesRoutes(fiberApp, mockCompaniesService{
			t:           t,
			returnError: services.ErrDb{},
			expectedCompanyUpdate: services.CompanyUpdate{
				ID:                "605c72efb1e2c3d1f8a1b2c3",
				Name:              "name",
				Description:       &description,
				AmountOfEmployees: 100,
				Registered:        &registered,
				Type:              "Sole Proprietorship",
			},
		}, nil)

		body := `{
			"name":"name",
			"description":"description",
			"amount_of_employees":100,
			"registered":true,
			"type":"Sole Proprietorship"
		}`
		req := httptest.NewRequest("PATCH", "/companies/605c72efb1e2c3d1f8a1b2c3", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		response, err := fiberApp.Test(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		defer response.Body.Close()
		require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	})
}

func TestDeleteCompany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fiberApp := initFiberApp()

		ch := make(chan any)
		handlers.SetupCompaniesRoutes(fiberApp, mockCompaniesService{
			t:          t,
			expectedId: "605c72efb1e2c3d1f8a1b2c3",
		}, newMockPublisher(ch))

		req := httptest.NewRequest("DELETE", "/companies/605c72efb1e2c3d1f8a1b2c3", nil)

		response, err := fiberApp.Test(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		defer response.Body.Close()
		require.Equal(t, fiber.StatusNoContent, response.StatusCode)

		select {
		case e := <-ch:
			require.Equal(t, "605c72efb1e2c3d1f8a1b2c3", e)
		case <-time.After(time.Millisecond * 500):
			require.Fail(t, "timeout")
		}
	})

	t.Run("bad request error", func(t *testing.T) {
		fiberApp := initFiberApp()

		handlers.SetupCompaniesRoutes(fiberApp, nil, nil)

		req := httptest.NewRequest("DELETE", "/companies/wrong_id", nil)

		response, err := fiberApp.Test(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		defer response.Body.Close()
		require.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	})

	t.Run("internal server error", func(t *testing.T) {
		fiberApp := initFiberApp()

		handlers.SetupCompaniesRoutes(fiberApp, mockCompaniesService{
			t:           t,
			returnError: services.ErrDb{},
			expectedId:  "605c72efb1e2c3d1f8a1b2c3",
		}, nil)

		req := httptest.NewRequest("DELETE", "/companies/605c72efb1e2c3d1f8a1b2c3", nil)

		response, err := fiberApp.Test(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		defer response.Body.Close()
		require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	})
}

func initFiberApp() *fiber.App {
	return fiber.New(fiber.Config{})
}

type mockCompaniesService struct {
	t                     *testing.T
	expectedCompany       services.Company
	expectedCompanyUpdate services.CompanyUpdate
	expectedId            string
	returnCompany         services.Company
	returnError           error
}

func (m mockCompaniesService) Create(ctx context.Context, company services.Company) (services.Company, error) {
	m.t.Helper()

	require.Equal(m.t, m.expectedCompany, company)
	return m.returnCompany, m.returnError
}

func (m mockCompaniesService) Get(ctx context.Context, id string) (services.Company, error) {
	m.t.Helper()

	require.Equal(m.t, m.expectedId, id)
	return m.returnCompany, m.returnError
}

func (m mockCompaniesService) Update(ctx context.Context, update services.CompanyUpdate) error {
	m.t.Helper()

	require.Equal(m.t, m.expectedCompanyUpdate, update)
	return m.returnError
}

func (m mockCompaniesService) Delete(ctx context.Context, id string) error {
	m.t.Helper()

	require.Equal(m.t, m.expectedId, id)
	return m.returnError
}

type mockPublisher struct {
	ch chan<- any
}

func newMockPublisher(c chan<- any) *mockPublisher {
	return &mockPublisher{ch: c}
}

func (m *mockPublisher) OnCreateCompany(e any) {
	m.ch <- e
}

func (m *mockPublisher) OnPatchCompany(e any) {
	m.ch <- e
}

func (m *mockPublisher) OnDeleteCompany(e any) {
	m.ch <- e
}
