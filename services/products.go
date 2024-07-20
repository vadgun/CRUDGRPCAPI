package services

import (
	pb "api-gateway/api-gateway"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

func getProductgRPCClient() pb.ProductServiceClient {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient("localhost:50051", opts)
	if err != nil {
		fmt.Println(err, "error conecting on gRPC server")
	}

	productgRPCClient := pb.NewProductServiceClient(conn)
	return productgRPCClient
}

func convertToProduct(pbProduct *pb.ProductResponse) Product {
	return Product{
		Id:          pbProduct.Id,
		Name:        pbProduct.Name,
		Description: pbProduct.Description,
		Price:       float64(pbProduct.Price),
		Quantity:    int(pbProduct.Quantity),
	}
}

func CreateProduct(c *gin.Context) {
	// Crear un producto usando gRPC
	var product Product
	grpcClient := getProductgRPCClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := grpcClient.CreateProduct(ctx, &pb.CreateProductRequest{Name: product.Name, Description: product.Description, Price: float32(product.Price), Quantity: int32(product.Quantity)})
	if err != nil {
		log.Printf("could not create the product: %v", err)
	}
	log.Printf("Created product with id: %s", r.GetId())
	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully!", "product": r.Id})

}

func GetProducts(c *gin.Context) {
	// Listar productos usando gRPC
	grpcClient := getProductgRPCClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := grpcClient.GetProducts(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("could not read the products: %v", err)
	}
	//log.Printf("List of products: %v", r.Products)
	if len(r.Products) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No product found"})
	} else {
		c.JSON(http.StatusOK, r.Products)
	}

}

func GetProductByID(c *gin.Context) {
	// Obtener un solo producto por ID usando gRPC
	grpcClient := getProductgRPCClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	id := c.Param("id")
	r, err := grpcClient.GetProductByID(ctx, &pb.ProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		log.Printf("could not find the product: %v", err)
		return
	}
	log.Printf("product found: %v", r.Name)
	product := convertToProduct(r)
	c.JSON(http.StatusOK, product)
}

func UpdateProduct(c *gin.Context) {
	// Actualiza un solo producto usando gRPC
	var product Product
	grpcClient := getProductgRPCClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	id := c.Param("id")
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := grpcClient.UpdateProduct(ctx, &pb.UpdateProductRequest{Id: id, Name: product.Name, Description: product.Description, Price: float32(product.Price), Quantity: int32(product.Quantity)})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		log.Printf("could not find the product: %v", err)
		return
	}
	log.Printf("product found: %v", r.Name)
	product = convertToProduct(r)
	c.JSON(http.StatusOK, product)

}

func DeleteProduct(c *gin.Context) {
	// Eliminar un solo producto usando gRPC
	grpcClient := getProductgRPCClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	id := c.Param("id")

	_, err := grpcClient.DeleteProduct(ctx, &pb.DeleteProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		log.Printf("could not find the product: %v", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
