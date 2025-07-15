package openapi

import (
	"context"
	"example-server/internal/openapi/ogen"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ItemsService struct{}

func (s *ItemsService) NewError(ctx context.Context, err error) *ogen.ErrorResponseStatusCode {
	return &ogen.ErrorResponseStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response:   ogen.ErrorResponse{Error: err.Error()},
	}
}

func (s *ItemsService) PingGet(
	ctx context.Context,
) (*ogen.PingGetOK, error) {
	log.Info().Msg("Handling PingGet...")
	return &ogen.PingGetOK{
		Message: "pong",
	}, nil
}

// func (s *ItemsService) ItemsAllGet(
// 	ctx context.Context,
// 	params ogen.ItemsAllGetParams,
// ) (*ogen.GetItemsResponse, error) {
// 	// MOCK DATA
// 	return &ogen.GetItemsResponse{
// 		Data: []ogen.Item{
// 			{
// 				ID:        1,
// 				UUID:      uuid.New(),
// 				CreatedAt: time.Now(),
// 				Name:      "Item 1",
// 				Price:     100,
// 			},
// 		},
// 		Meta: ogen.GetItemsResponseMeta{},
// 	}, nil
// }

// func (s *ItemsService) ItemsGet(
// 	ctx context.Context,
// 	params ogen.ItemsGetParams,
// ) (*ogen.GetItemsResponse, error) {
// 	// MOCK DATA
// 	return &ogen.GetItemsResponse{
// 		Data: []ogen.Item{
// 			{
// 				ID:        1,
// 				UUID:      uuid.New(),
// 				CreatedAt: time.Now(),
// 				Name:      "Item 1",
// 				Price:     100,
// 			},
// 		},
// 		Meta: ogen.GetItemsResponseMeta{},
// 	}, nil
// }

// func (s *ItemsService) ItemsIDGet(
// 	ctx context.Context,
// 	params ogen.ItemsIDGetParams,
// ) (*ogen.GetItemResponse, error) {
// 	// MOCK DATA
// 	return &ogen.GetItemResponse{
// 		Data: ogen.Item{
// 			ID:        1,
// 			UUID:      uuid.New(),
// 			CreatedAt: time.Now(),
// 			Name:      "Item 1",
// 			Price:     100,
// 		},
// 	}, nil
// }

// func (s *ItemsService) ItemsPost(
// 	ctx context.Context,
// 	req *ogen.CreateItemRequest,
// ) (*ogen.CreateItemResponse, error) {
// 	// MOCK DATA
// 	return &ogen.CreateItemResponse{
// 		Data: ogen.Item{
// 			ID:        1,
// 			UUID:      uuid.New(),
// 			CreatedAt: time.Now(),
// 			Name:      req.Data.Name,
// 			Price:     req.Data.Price,
// 		},
// 		Meta: ogen.CreateItemResponseMeta{
// 			Created: true,
// 		},
// 	}, nil
// }
