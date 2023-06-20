package repos

import (
	"example-server/models"
)

func GetItemsByIds(ids []int) ([]models.Item, error) {
	// stub
	return []models.Item{}, nil
}

func InsertItem(itemIn models.ItemIn) (models.Item, error) {
	// stub
	return models.Item{}, nil
}
