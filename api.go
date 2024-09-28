package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	db         Storer
	listenAddr string
}

func NewServer(listenAddr string, db Storer) *Server {
	return &Server{
		db:         db,
		listenAddr: listenAddr,
	}
}

// CreateOrderHandler handles POST requests to create a new order
func (s *Server) CreateOrderHandler(c *gin.Context) {
	var request Coordinates
	// Bind JSON to request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate the coordinates
	if err := ValidateCoordinates(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert origin and destination to float64 for haversine
	originLat, _ := strconv.ParseFloat(request.Origin[0], 64)
	originLon, _ := strconv.ParseFloat(request.Origin[1], 64)
	destLat, _ := strconv.ParseFloat(request.Destination[0], 64)
	destLon, _ := strconv.ParseFloat(request.Destination[1], 64)

	// Compute distance using the haversine function
	distance := haversine(originLat, originLon, destLat, destLon)

	// Create the order based on validated coordinates and computed distance
	order := Order{
		Distance: distance,
		Status:   "UNASSIGNED", // Using the Coordinates struct directly
	}

	// Insert the order into the database and retrieve the generated ID
	ctx := context.Background()
	createdOrder, err := s.db.InsertOrder(ctx, order) // Get the order with ID
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Return the created order, including the generated ID
	c.JSON(http.StatusCreated, gin.H{
		"id":       createdOrder.ID,
		"distance": createdOrder.Distance,
		"status":   createdOrder.Status,
	})
}

func (s *Server) GetOrderHandler(c *gin.Context) {
	// Extract order ID from URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Get the order using the Storer
	ctx := context.Background()
	order, err := s.db.GetOrder(ctx, int(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Order with ID %d not found", id)})
		return
	}

	// Return the fetched order
	c.JSON(http.StatusOK, order)
}

// ModifyOrderHandler handles PUT requests to update an existing order by ID
// TakeOrderHandler handles PATCH requests to take an order
func (s *Server) TakeOrderHandler(c *gin.Context) {
	orderID := c.Param("id")
	orderIDInt, _ := strconv.Atoi(orderID)

	// Fetch the order to check its current status
	order, err := s.db.GetOrder(context.Background(), orderIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check if the order is already taken
	if order.Status == "TAKEN" {
		c.JSON(http.StatusConflict, gin.H{"error": "Order is already taken"})
		return
	}

	// Update the order status to "TAKEN"
	order.Status = "TAKEN"

	// Use optimistic locking to update the order
	if err := s.db.UpdateOrder(context.Background(), order); err != nil {
		// If the error is due to a version conflict, handle it
		if err.Error() == "record not found" {
			c.JSON(http.StatusConflict, gin.H{"error": "Order has already been taken by another request"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS"})
}

// GetOrdersHandler handles GET requests to retrieve a paginated list of orders

func (s *Server) GetOrdersHandler(c *gin.Context) {
	// Extract query parameters
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	// Validate and parse the page and limit parameters
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Fetch orders from the database
	orders, err := s.db.GetOrders(context.Background(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// Prepare response
	response := make([]gin.H, len(orders))
	for i, order := range orders {
		response[i] = gin.H{
			"id":       order.ID,
			"distance": order.Distance,
			"status":   order.Status,
		}
	}

	c.JSON(http.StatusOK, response)
	// Return the list of orders
}

func SetupRoutes(s *Server) *gin.Engine {

	router := gin.Default()

	router.GET("/orders", s.GetOrdersHandler)
	router.POST("/orders", s.CreateOrderHandler)
	router.GET("/orders/:id", s.GetOrderHandler)
	router.PATCH("/orders/:id", s.TakeOrderHandler)

	return router

}
