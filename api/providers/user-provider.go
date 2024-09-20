package providers

import (
	"api/models"
	ddbproxy "api/proxy/ddb"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

var usersTableName = "Users"

func CreateUser() models.User {
	id := uuid.NewString()

	hashKeyAttributeValues := map[string]types.AttributeValue{
		":uId": &types.AttributeValueMemberS{Value: id},
	}

	results := ddbproxy.QueryTable[models.User](usersTableName, "id = :uId", hashKeyAttributeValues)

	if len(results) == 0 {
		user := models.User{
			Id:           id,
			HouseholdIds: []string{},
		}

		ddbproxy.CreateItem(usersTableName, user)

		return user
	}

	return CreateUser()
}

func UpdateUser(user models.User) error {
	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: user.Id},
	}
	ignoreKeys := []string{"id"}

	return ddbproxy.UpdateItem(usersTableName, key, user, ignoreKeys)
}

func GetUsers(id string) []models.User {
	hashKeyAttributeValues := map[string]types.AttributeValue{
		":uId": &types.AttributeValueMemberS{Value: id},
	}

	return ddbproxy.QueryTable[models.User](usersTableName, "id = :uId", hashKeyAttributeValues)
}

func GetOrCreateUser(id string) models.User {
	users := GetUsers(id)

	if len(users) > 0 {
		return users[0]
	}

	return CreateUser()
}
