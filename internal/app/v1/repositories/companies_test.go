package repositories_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/repositories"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var testCompaniesCollection *mongo.Collection
var brokenMongoCollection *mongo.Collection

func TestCreate(t *testing.T) {
	t.Run("create company successfully", func(t *testing.T) {
		repo := repositories.NewCompaniesRepository(testCompaniesCollection)

		company := createTestCompany("TestCreate")
		createdCompany, err := repo.Create(context.Background(), company)
		require.NoError(t, err)
		require.NotEmpty(t, createdCompany.ID)

		var cmp repositories.Company
		testCompaniesCollection.FindOne(context.Background(), bson.M{"_id": createdCompany.ID}).Decode(&cmp)
		require.Equal(t, createdCompany, cmp)
	})

	t.Run("create company failed", func(t *testing.T) {
		repo := repositories.NewCompaniesRepository(brokenMongoCollection)
		company := createTestCompany("name")
		createdCompany, err := repo.Create(context.Background(), company)
		require.Error(t, err)
		require.Empty(t, createdCompany)
	})
}

func TestGet(t *testing.T) {
	t.Run("getting company successfully", func(t *testing.T) {
		company := createTestCompany("TestGet")

		_, err := testCompaniesCollection.InsertOne(context.Background(), company)
		require.NoError(t, err)

		repo := repositories.NewCompaniesRepository(testCompaniesCollection)

		resCompany, err := repo.Get(context.Background(), company.ID)
		require.NoError(t, err)
		require.Equal(t, company, resCompany)
	})

	t.Run("company not found", func(t *testing.T) {
		repo := repositories.NewCompaniesRepository(testCompaniesCollection)

		resCompany, err := repo.Get(context.Background(), "id")
		require.ErrorAs(t, err, &repositories.ErrNotFound{})
		require.Empty(t, resCompany)
	})

	t.Run("getting company failed", func(t *testing.T) {
		repo := repositories.NewCompaniesRepository(brokenMongoCollection)

		resCompany, err := repo.Get(context.Background(), "id")
		require.NotErrorIs(t, err, repositories.ErrNotFound{})
		require.Empty(t, resCompany)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("update company successfully, full update", func(t *testing.T) {
		company := createTestCompany("TestUpdate")
		_, err := testCompaniesCollection.InsertOne(context.Background(), company)
		require.NoError(t, err)

		repo := repositories.NewCompaniesRepository(testCompaniesCollection)

		newEmptyDescription := ""
		falseRegistered := false
		updateCompany := repositories.CompanyUpdate{
			ID:                company.ID,
			Name:              "new name",
			Description:       &newEmptyDescription,
			AmountOfEmployees: 200,
			Registered:        &falseRegistered,
			Type:              "CompanyTypeNonProfit",
		}

		err = repo.Update(context.Background(), updateCompany)
		require.NoError(t, err)

		var updatedCompany repositories.Company
		testCompaniesCollection.FindOne(context.Background(), bson.M{"_id": company.ID}).Decode(&updatedCompany)
		require.Equal(t, "new name", updatedCompany.Name)
		require.Equal(t, "", updatedCompany.Description)
		require.Equal(t, 200, updatedCompany.AmountOfEmployees)
		require.Equal(t, false, updatedCompany.Registered)
		require.Equal(t, "CompanyTypeNonProfit", updatedCompany.Type)
	})

	t.Run("update company successfully, partial update", func(t *testing.T) {
		company := createTestCompany("TestUpdate2")
		_, err := testCompaniesCollection.InsertOne(context.Background(), company)
		require.NoError(t, err)

		repo := repositories.NewCompaniesRepository(testCompaniesCollection)

		updateCompany := repositories.CompanyUpdate{
			ID:                company.ID,
			Name:              "new name",
			Description:       nil,
			AmountOfEmployees: 0,
			Registered:        nil,
			Type:              "CompanyTypeNonProfit",
		}

		err = repo.Update(context.Background(), updateCompany)
		require.NoError(t, err)

		var updatedCompany repositories.Company
		testCompaniesCollection.FindOne(context.Background(), bson.M{"_id": company.ID}).Decode(&updatedCompany)
		require.Equal(t, "new name", updatedCompany.Name)
		require.Equal(t, "test description", updatedCompany.Description)
		require.Equal(t, company.AmountOfEmployees, updatedCompany.AmountOfEmployees)
		require.Equal(t, true, updatedCompany.Registered)
		require.Equal(t, "CompanyTypeNonProfit", updatedCompany.Type)
	})

	t.Run("update company failed", func(t *testing.T) {
		repo := repositories.NewCompaniesRepository(brokenMongoCollection)

		updateCompany := repositories.CompanyUpdate{
			ID:                "id",
			Name:              "new name",
			Description:       nil,
			AmountOfEmployees: 200,
			Registered:        nil,
			Type:              "CompanyTypeNonProfit",
		}

		err := repo.Update(context.Background(), updateCompany)
		require.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("delete company successfully", func(t *testing.T) {
		company := createTestCompany("TestDelete")
		_, err := testCompaniesCollection.InsertOne(context.Background(), company)
		require.NoError(t, err)

		repo := repositories.NewCompaniesRepository(testCompaniesCollection)

		err = repo.Delete(context.Background(), company.ID)
		require.NoError(t, err)

		var deletedCompany repositories.Company
		err = testCompaniesCollection.FindOne(context.Background(), bson.M{"_id": company.ID}).Decode(&deletedCompany)
		require.ErrorIs(t, err, mongo.ErrNoDocuments)
	})

	t.Run("delete company failed", func(t *testing.T) {
		repo := repositories.NewCompaniesRepository(brokenMongoCollection)

		err := repo.Delete(context.Background(), "id")
		require.Error(t, err)
	})
}

func createTestCompany(name string) repositories.Company {
	return repositories.Company{
		ID:                bson.NewObjectId().Hex(),
		Name:              name,
		Description:       "test description",
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              "CompanyTypeCorporations",
	}
}

func TestMain(m *testing.M) {
	pool, resource := startMongoContainer()
	defer func() {
		resource.Expire(10) //nolint errcheck
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("could not purge resource: %s", err)
		}
	}()

	testCompaniesCollection = connectToMongoCollection()
	brokenMongoCollection = createBrokenMongoCollection()

	m.Run()
}

func startMongoContainer() (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not init docker pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("could not connect to Docker: %s", err)
	}

	runOpts := &dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "latest",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017/tcp": {{HostIP: "", HostPort: "37017"}},
		},
	}
	resource, err := pool.RunWithOptions(runOpts)
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	return pool, resource
}

func connectToMongoCollection() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:37017")
	clientOptions.SetConnectTimeout(time.Duration(2) * time.Second)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("failed to ping mongo: %v", err)
	}

	return client.Database("test_companies").Collection("companies")
}

func createBrokenMongoCollection() *mongo.Collection {
	clientOptions := options.Client()
	clientOptions.SetServerSelectionTimeout(time.Second)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	return client.Database("test_companies").Collection("companies")
}
