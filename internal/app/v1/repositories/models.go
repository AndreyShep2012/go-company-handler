package repositories

type Company struct {
	ID                string `bson:"_id"`
	Name              string `bson:"name"`
	Description       string `bson:"description"`
	AmountOfEmployees int    `bson:"amount_of_employees"`
	Registered        bool   `bson:"registered"`
	Type              string `bson:"type"`
}

type CompanyUpdate struct {
	ID                string
	Name              string
	Description       *string
	AmountOfEmployees int
	Registered        *bool
	Type              string
}
