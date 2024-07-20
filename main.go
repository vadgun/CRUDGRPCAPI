package main

import (
	auth "api-gateway/auth"
	restrict "api-gateway/restrict"
	services "api-gateway/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Rutas para el servicio de productos
	r.POST("/products", auth.AuthMiddleware(), services.CreateProduct)
	r.GET("/products", auth.AuthMiddleware(), services.GetProducts)
	r.GET("/products/:id", auth.AuthMiddleware(), services.GetProductByID)
	r.PATCH("/products/:id", auth.AuthMiddleware(), services.UpdateProduct)
	r.DELETE("/products/:id", auth.AuthMiddleware(), services.DeleteProduct)

	// Rutas para el servicio de pedidos
	r.POST("/orders", auth.AuthMiddleware(), services.PlaceOrder)
	r.GET("/orders", auth.AuthMiddleware(), services.GetOrders)

	// Middleware de limitación de tasa
	r.Use(restrict.RateLimitMiddleware())

	// Endpoint para iniciar sesión y obtener un token JWT
	r.POST("/login", services.Login)

	// gRPC Server on Default tcp:50051
	go services.GrpcServer()

	// HTTP Server on port:8080
	r.Run(":8080")
}
