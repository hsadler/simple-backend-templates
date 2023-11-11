package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"example-server/dependencies"
	"example-server/models"
	"example-server/repos"
)

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
// @Success 200 {object} models.GetItemsResponse
// @Router /api/items/all [get]
func HandleGetAllItems(deps *dependencies.Dependencies) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Fetch Items
		status, items := repos.FetchAllItems(deps.DBPool)
		if !status {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Items"})
			return
		}
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
			log.Println("Error parsing Item ID:", err)
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
			return
		}
		log.Println("Fetching item by id:", itemId)
		// Fetch Item by ID
		status, item := repos.FetchItemById(deps.DBPool, itemId)
		if !status {
			log.Println("Error fetching Item by id:", itemId)
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Item"})
			return
		}
		// Return Item if found otherwise 404
		if item != nil {
			g.JSON(http.StatusOK, models.GetItemResponse{Data: item, Meta: struct{}{}})
		} else {
			g.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		}
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
		log.Println("itemIdsStrArr:", itemIdsStrArr)
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
		} else {
			g.JSON(http.StatusBadRequest, gin.H{"error": "Missing Item IDs"})
			return
		}
		// Fetch Items by IDs
		status, items := repos.FetchItemsByIds(deps.DBPool, itemIds)
		if !status {
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
			http.StatusCreated,
			models.CreateItemResponse{Data: item, Meta: models.CreateItemResponseMeta{Created: true}},
		)
	}
}
