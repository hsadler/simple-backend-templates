package itemsrepo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"

	"example-server/internal/database"
	"example-server/internal/logger"
	"example-server/internal/models"
)

var (
	ErrorItemNotFound = errors.New("Item not found")
	ErrorItemExists   = errors.New("Item already exists")
	ErrorItemInsert   = errors.New("Error inserting Item")
	ErrorItemsQuery   = errors.New("Error querying Items")
)

func FetchPaginatedItems(dbPool database.PgxPoolIface, offset, chunkSize int) ([]*models.Item, error) {
	// Fetch paginated Items
	rows, err := dbPool.Query(
		context.Background(),
		"SELECT id, uuid, created_at, name, price FROM item ORDER BY id OFFSET $1 LIMIT $2",
		offset, chunkSize,
	)
	// Handle Items fetch error
	if err != nil {
		logger.LogErrorWithStacktrace(err, "Error querying Items")
		return nil, ErrorItemsQuery
	}
	defer rows.Close()
	// Iterate over rows and append Items
	var items []*models.Item
	for rows.Next() {
		var item models.Item
		// Scan Item and append to Items unless error
		if err := rows.Scan(&item.ID, &item.UUID, &item.CreatedAt, &item.Name, &item.Price); err != nil {
			logger.LogErrorWithStacktrace(err, "Error scanning Item")
			return nil, ErrorItemsQuery
		}
		items = append(items, &item)
	}
	// Handle row iteration error
	if err := rows.Err(); err != nil {
		logger.LogErrorWithStacktrace(err, "Error iterating over paginated Items")
		return nil, ErrorItemsQuery
	}
	// Check if the slice is nil and replace it with an empty slice
	if items == nil {
		items = []*models.Item{}
	}
	return items, nil
}

func FetchItemById(dbPool database.PgxPoolIface, itemId int) (*models.Item, error) {
	// Fetch Item by ID
	var item models.Item
	err := dbPool.QueryRow(
		context.Background(),
		"SELECT id, uuid, created_at, name, price FROM item WHERE id = $1",
		itemId,
	).Scan(&item.ID, &item.UUID, &item.CreatedAt, &item.Name, &item.Price)
	// Handle Item fetch error
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrorItemNotFound
		}
		logger.LogErrorWithStacktrace(err, "Error querying Item")
		return nil, ErrorItemsQuery
	}
	return &item, nil
}

func FetchItemsByIds(dbPool database.PgxPoolIface, itemIds []int) ([]*models.Item, error) {
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
		logger.LogErrorWithStacktrace(err, "Error querying Items")
		return nil, ErrorItemsQuery
	}
	defer rows.Close()
	// Iterate over rows and append Items
	var items []*models.Item
	for rows.Next() {
		var item models.Item
		// Scan Item and append to Items unless error
		if err := rows.Scan(&item.ID, &item.UUID, &item.CreatedAt, &item.Name, &item.Price); err != nil {
			logger.LogErrorWithStacktrace(err, "Error scanning Item")
			return nil, ErrorItemsQuery
		}
		items = append(items, &item)
	}
	// Handle row iteration error
	if err := rows.Err(); err != nil {
		logger.LogErrorWithStacktrace(err, "Error iterating over Items")
		return nil, ErrorItemsQuery
	}
	// Check if the slice is nil and replace it with an empty slice
	if items == nil {
		items = []*models.Item{}
	}
	return items, nil
}

func InsertItem(dbPool database.PgxPoolIface, itemIn models.ItemIn) (*models.Item, error) {
	// Insert Item
	var itemId int
	err := dbPool.QueryRow(
		context.Background(),
		"INSERT INTO item (name, price) VALUES ($1, $2) RETURNING id",
		itemIn.Name,
		itemIn.Price,
	).Scan(&itemId)
	// Handle Item insert error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Duplicate entry error handling
			if pgErr.Code == "23505" {
				return nil, ErrorItemExists
			}
		}
		logger.LogErrorWithStacktrace(err, "Error inserting Item")
		return nil, ErrorItemInsert
	}
	// Fetch Item by ID
	item, err := FetchItemById(dbPool, itemId)
	if err != nil {
		return nil, err
	}
	return item, nil
}
