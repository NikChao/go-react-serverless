package models

type User struct {
	Id           string   `json:"id" dynamodbav:"id"`
	HouseholdIds []string `json:"householdIds" dynamodbav:"householdIds"`
}
