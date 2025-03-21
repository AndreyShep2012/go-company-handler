package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/handlers"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestCreateCompany(t *testing.T) {
	client := &http.Client{}

	createBody := `{
		"name":"%s",
		"description":"description",
		"amount_of_employees":100,
		"registered":true,
		"type":"Sole Proprietorship"
	}`

	var name string
	require.NoError(t, faker.FakeData(&name, options.WithRandomStringLength(10)))
	createBody = fmt.Sprintf(createBody, name)

	req, err := http.NewRequest("POST", createRequestUrl(testConf.ListenAddr, "/companies/create"), bytes.NewReader([]byte(createBody)))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", createToken(t, "test", []byte(testConf.JWTSecretKey)))
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var res handlers.Company
	err = json.Unmarshal(bodyBytes, &res)
	require.NoError(t, err)
	require.NotEmpty(t, res.ID)

	expectedResponse := handlers.Company{
		ID:                res.ID,
		Name:              name,
		Description:       "description",
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              "Sole Proprietorship",
	}

	require.Equal(t, expectedResponse, res)

	// try one more time to check unique of company name
	req, err = http.NewRequest("POST", createRequestUrl(testConf.ListenAddr, "/companies/create"), bytes.NewReader([]byte(createBody)))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", createToken(t, "test", []byte(testConf.JWTSecretKey)))
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestGetCompany(t *testing.T) {
	client := &http.Client{}

	id := createCompany(t)

	req, err := http.NewRequest("GET", createRequestUrl(testConf.ListenAddr, "/companies/"+id), nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var res handlers.Company
	err = json.Unmarshal(bodyBytes, &res)
	require.NoError(t, err)

	expectedResponse := handlers.Company{
		ID:                id,
		Name:              res.Name,
		Description:       "description",
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              "Sole Proprietorship",
	}

	require.Equal(t, expectedResponse, res)
}

func TestDeleteCompany(t *testing.T) {
	client := &http.Client{}

	id := createCompany(t)

	req, err := http.NewRequest("DELETE", createRequestUrl(testConf.ListenAddr, "/companies/"+id), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", createToken(t, "test", []byte(testConf.JWTSecretKey)))

	resp, err := client.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// get company
	req, err = http.NewRequest("GET", createRequestUrl(testConf.ListenAddr, "/companies/"+id), nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUpdateCompany(t *testing.T) {
	client := &http.Client{}

	id := createCompany(t)

	patchBody := `{
		"name":"%s",
		"description":"",
		"amount_of_employees":10,
		"registered":false,
		"type":"Corporations"
	}`

	var newName string
	require.NoError(t, faker.FakeData(&newName, options.WithRandomStringLength(10)))

	patchBody = fmt.Sprintf(patchBody, newName)

	req, err := http.NewRequest("PATCH", createRequestUrl(testConf.ListenAddr, "/companies/"+id), bytes.NewReader([]byte(patchBody)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", createToken(t, "test", []byte(testConf.JWTSecretKey)))

	resp, err := client.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// get company
	req, err = http.NewRequest("GET", createRequestUrl(testConf.ListenAddr, "/companies/"+id), nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var res handlers.Company
	err = json.Unmarshal(bodyBytes, &res)
	require.NoError(t, err)

	expectedResponse := handlers.Company{
		ID:                id,
		Name:              newName,
		Description:       "",
		AmountOfEmployees: 10,
		Registered:        false,
		Type:              "Corporations",
	}

	require.Equal(t, expectedResponse, res)
}

func createCompany(t *testing.T) string {
	t.Helper()

	client := &http.Client{}

	createBody := `{
		"name":"%s",
		"description":"description",
		"amount_of_employees":100,
		"registered":true,
		"type":"Sole Proprietorship"
	}`
	var name string
	require.NoError(t, faker.FakeData(&name, options.WithRandomStringLength(10)))

	createBody = fmt.Sprintf(createBody, name)

	req, err := http.NewRequest("POST", createRequestUrl(testConf.ListenAddr, "/companies/create"), bytes.NewReader([]byte(createBody)))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", createToken(t, "test", []byte(testConf.JWTSecretKey)))
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var res handlers.Company
	err = json.Unmarshal(bodyBytes, &res)
	require.NoError(t, err)
	require.NotEmpty(t, res.ID)

	return res.ID
}

func createRequestUrl(listenAddr, suffix string) string {
	return fmt.Sprintf("http://%s/api/v1%s", listenAddr, suffix)
}

func createToken(t *testing.T, username string, secretKey []byte) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	require.NoError(t, err)

	return tokenString
}
