package providers

import (
	"api/models"
	ddbproxy "api/proxy/ddb"
	"fmt"

	"github.com/google/uuid"
)

var householdTableName = "Households"

func CreateHousehold() models.Household {
	householdId := uuid.NewString()
	household := models.Household{
		Id: householdId,
	}
	ddbproxy.CreateItem(householdTableName, household)

	return household
}

func JoinHousehold(userId string, householdId string) error {
	user := GetOrCreateUser(userId)

	// TODO: support joining multiple households,
	// for now this makes life simpler. If you join one you'll leave others
	user.HouseholdIds = []string{householdId}
	return UpdateUser(user)
}

func LeaveHousehold(userId string, householdIdToRemove string) error {
	users := GetUsers(userId)

	if len(users) == 0 {
		return fmt.Errorf("could not find user [%s]", userId)
	}

	user := users[0]

	newHouseholdIds := make([]string, 0)

	for _, householdId := range user.HouseholdIds {
		if householdId != householdIdToRemove {
			newHouseholdIds = append(newHouseholdIds, householdId)
		}
	}

	user.HouseholdIds = newHouseholdIds
	return UpdateUser(user)
}
