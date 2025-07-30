package openapi

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"example-server/internal/dependencies"
	"example-server/internal/models"
	"example-server/internal/openapi/ogen"
	itemsrepo "example-server/internal/repos"
)

type ItemsService struct {
	Deps *dependencies.Dependencies
}

func (s *ItemsService) NewError(ctx context.Context, err error) *ogen.ErrorResponseStatusCode {
	return &ogen.ErrorResponseStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response:   ogen.ErrorResponse{Error: err.Error()},
	}
}

func (s *ItemsService) PingGet(
	ctx context.Context,
) (*ogen.PingGetOK, error) {
	log.Info().Msg("Handling PingGet")
	return &ogen.PingGetOK{
		Message: "pong",
	}, nil
}

func (s *ItemsService) ItemsPost(
	ctx context.Context,
	req *ogen.CreateItemRequest,
) (*ogen.CreateItemResponse, error) {
	log.Info().
		Str("name", req.Data.Name).
		Float32("price", req.Data.Price).
		Msg("Handling ItemsPost")
	// Insert item into DB
	item, err := itemsrepo.InsertItem(s.Deps.DBPool, models.ItemIn{
		Name:  req.Data.Name,
		Price: req.Data.Price,
	})
	if err != nil {
		return nil, s.NewError(ctx, err)
	}
	// Convert models.Item to ogen.Item
	itemOut := ogen.Item{
		ID:        int64(item.ID),
		UUID:      uuid.MustParse(item.UUID),
		CreatedAt: item.CreatedAt,
		Name:      item.Name,
		Price:     item.Price,
	}
	// Compose and return response
	return &ogen.CreateItemResponse{
		Data: itemOut,
		Meta: ogen.CreateItemResponseMeta{
			Created: true,
		},
	}, nil
}

func (s *ItemsService) ItemsGet(
	ctx context.Context,
	params ogen.ItemsGetParams,
) (*ogen.GetItemsResponse, error) {
	log.Info().
		Ints("item_ids", params.ItemIds).
		Msg("Handling ItemsGet")
	// Fetch items from DB by IDs
	items, err := itemsrepo.FetchItemsByIds(s.Deps.DBPool, params.ItemIds)
	if err != nil {
		return nil, s.NewError(ctx, err)
	}
	// Convert models.Items to ogen.Items
	itemsOut := make([]ogen.Item, len(items))
	for i, item := range items {
		itemsOut[i] = ogen.Item{
			ID:        int64(item.ID),
			UUID:      uuid.MustParse(item.UUID),
			CreatedAt: item.CreatedAt,
			Name:      item.Name,
			Price:     item.Price,
		}
	}
	// Compose and return response
	return &ogen.GetItemsResponse{
		Data: itemsOut,
		Meta: ogen.GetItemsResponseMeta{},
	}, nil
}

func (s *ItemsService) ItemsIDGet(
	ctx context.Context,
	params ogen.ItemsIDGetParams,
) (*ogen.GetItemResponse, error) {
	log.Info().
		Int("item_id", params.ID).
		Msg("Handling ItemsIDGet")
	// Fetch item from DB by ID
	item, err := itemsrepo.FetchItemById(s.Deps.DBPool, params.ID)
	if err != nil {
		return nil, s.NewError(ctx, err)
	}
	// Convert models.Item to ogen.Item
	itemOut := ogen.Item{
		ID:        int64(item.ID),
		UUID:      uuid.MustParse(item.UUID),
		CreatedAt: item.CreatedAt,
		Name:      item.Name,
		Price:     item.Price,
	}
	// Compose and return response
	return &ogen.GetItemResponse{
		Data: itemOut,
	}, nil
}

func (s *ItemsService) ItemsAllGet(
	ctx context.Context,
	params ogen.ItemsAllGetParams,
) (*ogen.GetItemsResponse, error) {
	log.Info().Msg("Handling ItemsAllGet")
	// Fetch all items from DB
	// Note: Large limit to get all items
	items, err := itemsrepo.FetchPaginatedItems(s.Deps.DBPool, 0, 1000)
	if err != nil {
		return nil, s.NewError(ctx, err)
	}
	// Convert models.Items to ogen.Items
	itemsOut := make([]ogen.Item, len(items))
	for i, item := range items {
		itemsOut[i] = ogen.Item{
			ID:        int64(item.ID),
			UUID:      uuid.MustParse(item.UUID),
			CreatedAt: item.CreatedAt,
			Name:      item.Name,
			Price:     item.Price,
		}
	}
	// Compose and return response
	return &ogen.GetItemsResponse{
		Data: itemsOut,
		Meta: ogen.GetItemsResponseMeta{},
	}, nil
}
