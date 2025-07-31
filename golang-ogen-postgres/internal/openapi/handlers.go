package openapi

import (
	"context"
	"errors"
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

func (s *ItemsService) Ping(
	ctx context.Context,
) (*ogen.PingResponse, error) {
	log.Info().Msg("Handling PingGet")
	return &ogen.PingResponse{
		Message: "pong",
	}, nil
}

func (s *ItemsService) CreateItem(
	ctx context.Context,
	req *ogen.ItemCreateRequest,
) (ogen.CreateItemRes, error) {
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
	return &ogen.ItemCreateResponse{
		Data: itemOut,
		Meta: ogen.ItemMeta{
			ItemStatus: ogen.OptItemMetaItemStatus{
				Value: ogen.ItemMetaItemStatusCreated,
				Set:   true,
			},
		},
	}, nil
}

func (s *ItemsService) GetItem(
	ctx context.Context,
	params ogen.GetItemParams,
) (ogen.GetItemRes, error) {
	// STUB
	return nil, errors.New("not implemented")
}

func (s *ItemsService) UpdateItem(
	ctx context.Context,
	req *ogen.ItemUpdateRequest,
	params ogen.UpdateItemParams,
) (ogen.UpdateItemRes, error) {
	// STUB
	return nil, errors.New("not implemented")
}

func (s *ItemsService) DeleteItem(
	ctx context.Context,
	params ogen.DeleteItemParams,
) (ogen.DeleteItemRes, error) {
	// STUB
	return nil, errors.New("not implemented")
}
