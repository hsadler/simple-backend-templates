package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"example-server/dependencies"
	"example-server/models"
	"example-server/repos"
)

// HELPERS

func parseQueryParam(g *gin.Context, key string, defaultValue interface{}) interface{} {
	valueStr := g.Query(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// ITEMS API

func SetupItemsAPIRoutes(router *gin.Engine, deps *dependencies.Dependencies) {
	itemsRouterGroup := router.Group("/api/items")
	itemsRouterGroup.GET("/all", HandleGetAllItems(deps))
	itemsRouterGroup.GET("/:id", HandleGetItem(deps))
	itemsRouterGroup.GET("", HandleGetItems(deps))
	itemsRouterGroup.POST("", HandleCreateItem(deps))
}

// GetAllItems godoc
// @Summary Get All Items
// @Description Returns all Items.
// @Tags items
// @Produce json
// @Param offset query int true "Offset" minimum(0)
// @Param chunkSize query int true "Chunk size" minimum(1) maximum(20)
// @Success 200 {object} models.GetItemsResponse
// @Router /api/items/all [get]
func HandleGetAllItems(deps *dependencies.Dependencies) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Parse query params and validate
		offsetParam := parseQueryParam(g, "offset", -1)
		chunkSizeParam := parseQueryParam(g, "chunkSize", -1)
		offset, offsetOk := offsetParam.(int)
		chunkSize, chunkSizeOk := chunkSizeParam.(int)
		if !offsetOk || !chunkSizeOk || offset < 0 || chunkSize < 1 || chunkSize > 20 {
			log.Warn().
				Msg("Invalid query parameters received on /api/items/all")
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
			return
		}
		log.Info().
			Int("offset", offset).
			Int("chunkSize", chunkSize).
			Msg("Fetching all items")
		// Fetch Items
		items, err := repos.FetchPaginatedItems(deps.DBPool, offset, chunkSize)
		if err != nil {
			log.Error().
				Err(err).
				Int("offset", offset).
				Int("chunkSize", chunkSize).
				Msg("Problem fetching paginated items")
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Items"})
			return
		}
		log.Info().
			Int("numItems", len(items)).
			Msg("Fetched items")
		// Return response
		g.JSON(http.StatusOK, models.GetItemsResponse{Data: items, Meta: struct{}{}})
	}
}

// GetItem godoc
// @Summary Get Item
// @Description Returns Item by id.
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} models.GetItemResponse
// @Failure 404 {object} string "Item not found"
// @Router /api/items/{id} [get]
func HandleGetItem(deps *dependencies.Dependencies) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Parse Item ID
		itemId, err := strconv.Atoi(g.Param("id"))
		if err != nil {
			log.Warn().
				Msg("Invalid Item ID received on /api/items/:id")
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
			return
		}
		log.Info().
			Int("itemId", itemId).
			Msg("Fetching item by id")
		// Fetch Item by ID
		item, err := repos.FetchItemById(deps.DBPool, itemId)
		if err != nil {
			if errors.Is(err, repos.ErrorItemNotFound) {
				log.Warn().
					Int("itemId", itemId).
					Msg("Item not found")
				g.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
				return
			}
			log.Error().
				Err(err).
				Int("itemId", itemId).
				Msg("Problem fetching item by id")
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Item"})
			return
		}
		// Return Item if found otherwise 404
		if item == nil {
			log.Warn().
				Int("itemId", itemId).
				Msg("Item not found")
			g.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		}
		g.JSON(http.StatusOK, models.GetItemResponse{Data: item, Meta: struct{}{}})
	}
}

// GetItems godoc
// @Summary Get Items
// @Description Returns Items by ids. Only returns subset of Items found.
// @Tags items
// @Accept json
// @Produce json
// @Param item_ids query []int true "Item IDs" collectionFormat(multi)
// @Success 200 {array} models.GetItemsResponse
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
					log.Warn().
						Msg("Invalid Item ID received on /api/items")
					g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
					return
				}
			}
		} else {
			log.Warn().
				Msg("Missing Item IDs received on /api/items")
			g.JSON(http.StatusBadRequest, gin.H{"error": "Missing Item IDs"})
			return
		}
		// Fetch Items by IDs
		items, err := repos.FetchItemsByIds(deps.DBPool, itemIds)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Problem fetching items by ids")
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Items"})
			return
		}
		// Return response

		g.JSON(http.StatusOK, models.GetItemsResponse{Data: items, Meta: struct{}{}})
	}
}

// CreateItem godoc
// @Summary Create Item
// @Description Creates Item.
// @Tags items
// @Accept json
// @Produce json
// @Param createItemRequest body models.CreateItemRequest true "Create Item Request"
// @Success 201 {object} models.CreateItemResponse
// @Failure 409 {object} string "Item already exists"
// @Router /api/items [post]
func HandleCreateItem(deps *dependencies.Dependencies) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Deserialize request
		var createItemRequest models.CreateItemRequest
		if err := g.ShouldBindJSON(&createItemRequest); err != nil {
			log.Warn().
				Msg("Invalid JSON payload received on /api/items")
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}
		// Validate request ItemIn data
		if err := deps.Validator.Struct(createItemRequest.Data); err != nil {
			// log.Println("Error validating request:", err)
			log.Warn().
				Msg("Invalid Item data payload received on /api/items")
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item data payload"})
			return
		}
		// Insert Item
		item, err := repos.InsertItem(deps.DBPool, createItemRequest.Data)
		// Handle Item insert error
		if err != nil {
			if errors.Is(err, repos.ErrorItemExists) {
				g.JSON(
					http.StatusConflict,
					gin.H{"error": "Item already exists"},
				)
				return
			}
			log.Error().
				Err(err).
				Msg("Problem inserting item")
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Item"})
			return
		}
		// Return response
		log.Info().
			Int("itemId", item.ID).
			Msg("Created item")
		g.JSON(
			http.StatusCreated,
			models.CreateItemResponse{Data: item, Meta: models.CreateItemResponseMeta{Created: true}},
		)
	}
}
