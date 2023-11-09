package repos

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"example-server/database"
	"example-server/models"
)

func FetchAllItems(dbPool database.PgxPoolIface) (bool, []*models.Item) {
	// Fetch all Items
	rows, err := dbPool.Query(
		context.Background(),
		"SELECT id, uuid, created_at, name, price FROM item",
	)
	// Handle Items fetch error
	if err != nil {
		log.Println("Error querying Items:", err)
		return false, nil
	}
	defer rows.Close()
	// Iterate over rows and append Items
	var items []*models.Item
	for rows.Next() {
		var item models.Item
		// Scan Item and append to Items unless error
		if err := rows.Scan(&item.ID, &item.UUID, &item.CreatedAt, &item.Name, &item.Price); err != nil {
			log.Println("Error scanning Item:", err)
			return false, nil
		}
		items = append(items, &item)
	}
	return true, items
}

func FetchItemById(dbPool database.PgxPoolIface, itemId int) (bool, *models.Item) {
	// Fetch Item by ID
	var item models.Item
	fetchErr := dbPool.QueryRow(
		context.Background(),
		"SELECT id, uuid, created_at, name, price FROM item WHERE id = $1",
		itemId,
	).Scan(&item.ID, &item.UUID, &item.CreatedAt, &item.Name, &item.Price)
	// Handle Item fetch error
	if fetchErr != nil {
		// Handle Item not found error
		if fetchErr == pgx.ErrNoRows {
			log.Println("Item not found")
			return true, nil
		}
		log.Println("Error querying Item:", fetchErr)
		return false, nil
	}
	return true, &item
}

func FetchItemsByIds(dbPool database.PgxPoolIface, itemIds []int) (bool, []*models.Item) {
	// Fetch Items by IDs
	var err error
	var rows pgx.Rows
	if len(itemIds) > 0 {
		rows, err = dbPool.Query(
			context.Background(),
			"SELECT id, uuid, created_at, name, price FROM item WHERE id = ANY($1)",
			itemIds,
		)
	}
	// Handle Items fetch error
	if err != nil {
		log.Println("Error querying Items:", err)
		return false, nil
	}
	defer rows.Close()
	// Iterate over rows and append Items
	var items []*models.Item
	for rows.Next() {
		var item models.Item
		// Scan Item and append to Items unless error
		if err := rows.Scan(&item.ID, &item.UUID, &item.CreatedAt, &item.Name, &item.Price); err != nil {
			log.Println("Error scanning Item:", err)
			return false, nil
		}
		items = append(items, &item)
	}
	return true, items
}

func InsertItem(dbPool database.PgxPoolIface, itemIn models.ItemIn) (bool, *models.Item, *pgconn.PgError) {
	// Insert Item
	var itemId int
	insertErr := dbPool.QueryRow(
		context.Background(),
		"INSERT INTO item (name, price) VALUES ($1, $2) RETURNING id",
		itemIn.Name,
		itemIn.Price,
	).Scan(&itemId)
	// Handle Item insert error
	if insertErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(insertErr, &pgErr) {
			// Duplicate entry error handling
			if pgErr.Code == "23505" {
				log.Println("Duplicate Item entry error:", pgErr)
			}
		} else {
			log.Println("Error inserting Item:", insertErr)
		}
		return false, nil, pgErr
	}
	log.Printf("Inserted itemId: %+v\n", itemId)
	// Fetch Item by ID
	var item *models.Item
	status, item := FetchItemById(dbPool, itemId)
	if !status {
		return false, nil, nil
	}
	return true, item, nil
}
