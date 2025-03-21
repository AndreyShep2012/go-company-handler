package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Companies struct {
	collection *mongo.Collection
}

func NewCompaniesRepository(collection *mongo.Collection) *Companies {
	return &Companies{collection: collection}
}

func (r Companies) Create(ctx context.Context, company Company) (Company, error) {
	company.ID = bson.NewObjectId().Hex()
	_, err := r.collection.InsertOne(ctx, company)
	if err != nil {
		return Company{}, handleError(err)
	}
	return company, nil
}

func (m Companies) Get(ctx context.Context, id string) (Company, error) {
	var company Company
	err := m.collection.FindOne(ctx, getIdFilter(id)).Decode(&company)
	return company, handleError(err)
}

func (m Companies) Update(ctx context.Context, company CompanyUpdate) error {
	set := bson.M{}
	if company.Name != "" {
		set["name"] = company.Name
	}
	if company.Description != nil {
		set["description"] = *company.Description
	}
	if company.AmountOfEmployees != 0 {
		set["amount_of_employees"] = company.AmountOfEmployees
	}
	if company.Registered != nil {
		set["registered"] = *company.Registered
	}
	if company.Type != "" {
		set["type"] = company.Type
	}

	_, err := m.collection.UpdateOne(ctx, getIdFilter(company.ID), bson.M{"$set": set})
	return handleError(err)
}

func (m Companies) Delete(ctx context.Context, id string) error {
	_, err := m.collection.DeleteOne(ctx, getIdFilter(id))
	return err
}

func getIdFilter(id string) bson.M {
	return bson.M{"_id": id}
}
