package routes

import (
	"example-server/dependencies"
	"example-server/models"
	"example-server/repos"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ITEMS API

func SetupItemsAPIRoutes(router *gin.Engine, deps *dependencies.Dependencies) {
	itemsRouterGroup := router.Group("/api/items")
	itemsRouterGroup.GET("/:id", HandleGetItem(deps))
	itemsRouterGroup.GET("", HandleGetItems(deps))
	itemsRouterGroup.POST("", HandleCreateItem(deps))
}

type GetItemResponse struct {
	Data *models.Item `json:"data"`
	Meta struct{}     `json:"meta"`
}

// GetItem godoc
// @Summary Get Item
// @Description Returns Item by id.
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} routes.GetItemResponse
// @Failure 400 {object} string
// @Router /api/items/{id} [get]
func HandleGetItem(deps *dependencies.Dependencies) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Parse Item ID
		itemId, err := strconv.Atoi(g.Param("id"))
		if err != nil {
			log.Println("Error parsing Item ID:", err)
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
			return
		}
		// Fetch Item by ID
		status, item := repos.FetchItemById(deps.DBPool, itemId)
		if !status {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Item"})
			return
		}
		// Return Item if found otherwise 404
		if item != nil {
			g.JSON(http.StatusOK, GetItemResponse{Data: item, Meta: struct{}{}})
		} else {
			g.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		}
	}
}

type GetItemsResponse struct {
	Data []*models.Item `json:"data"`
	Meta struct{}       `json:"meta"`
}

// GetItems godoc
// @Summary Get Items
// @Description Returns Items by ids. Only returns subset of Items found.
// @Tags items
// @Accept json
// @Produce json
// @Param item_ids query []int true "Item IDs" collectionFormat(multi)
// @Success 200 {array} routes.GetItemsResponse
// @Failure 400 {object} string
// @Router /api/items [get]
func HandleGetItems(deps *dependencies.Dependencies) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Parse Item IDs
		var itemIds []int
		var err error
		itemIdsStrArr, ok := g.GetQueryArray("item_ids")
		if ok {
			itemIds = make([]int, len(itemIdsStrArr))
			for i, itemIdStr := range itemIdsStrArr {
				itemIds[i], err = strconv.Atoi(itemIdStr)
				// Handle Item ID parse error
				if err != nil {
					log.Println("Error parsing Item ID:", err)
					g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
					return
				}
			}
		}
		// Fetch Items by IDs
		status, items := repos.FetchItemsByIds(deps.DBPool, itemIds)
		if !status {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Items"})
			return
		}
		// Return response
		g.JSON(http.StatusOK, GetItemsResponse{Data: items, Meta: struct{}{}})
	}
}

type CreateItemRequest struct {
	Data models.ItemIn `json:"data"`
}

type CreateItemResponseMeta struct {
	Created bool `json:"created"`
}

type CreateItemResponse struct {
	Data *models.Item           `json:"data"`
	Meta CreateItemResponseMeta `json:"meta"`
}

// CreateItem godoc
// @Summary Create Item
// @Description Creates Item.
// @Tags items
// @Accept json
// @Produce json
// @Param createItemRequest body routes.CreateItemRequest true "Create Item Request"
// @Success 200 {object} routes.CreateItemResponse
// @Failure 400 {object} string
// @Router /api/items [post]
func HandleCreateItem(deps *dependencies.Dependencies) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Deserialize request
		var createItemRequest CreateItemRequest
		if err := g.ShouldBindJSON(&createItemRequest); err != nil {
			log.Println("Error deserializing request:", err)
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}
		// Validate request ItemIn data
		if err := deps.Validator.Struct(createItemRequest.Data); err != nil {
			log.Println("Error validating request:", err)
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item data payload"})
			return
		}
		// Insert Item
		status, item, pgErr := repos.InsertItem(deps.DBPool, createItemRequest.Data)
		// Handle Item insert error
		if !status {
			if pgErr != nil {
				// Duplicate entry error handling
				if pgErr.Code == "23505" {
					log.Println("Duplicate Item entry error:", pgErr)
					g.JSON(
						http.StatusConflict,
						gin.H{"error": "Item already exists"},
					)
					return
				}
			}
			log.Println("Error inserting Item:", pgErr)
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Item"})
			return
		}
		// Return response
		g.JSON(
			http.StatusOK,
			CreateItemResponse{Data: item, Meta: CreateItemResponseMeta{Created: true}},
		)
	}
}
