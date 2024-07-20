package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	pb "api-gateway/api-gateway"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

func getOrdergRPCClient() pb.OrderServiceClient {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient("localhost:50051", opts)
	if err != nil {
		fmt.Println(err, "error conecting on gRPC server")
	}

	ordergRPCClient := pb.NewOrderServiceClient(conn)
	return ordergRPCClient
}

func verifyProduct(item *pb.OrderItem) (bool, *pb.ProductResponse, error) {
	grpcClient := getProductgRPCClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := grpcClient.GetProducts(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("could not read the products: %v", err)
	}

	for _, product := range r.Products {
		if product.Id == item.ProductId {
			//Product founded and has quantity
			if product.Quantity >= item.Quantity {
				return true, product, nil
			} else {
				return false, nil, errors.New("Product does not have enough quantity to order")
			}
		} else {
			//Product not founded
			return false, nil, errors.New("Product not found")
		}
	}
	return false, nil, errors.New("No products found")
}

func PlaceOrder(c *gin.Context) {
	// Realizar pedido con gRPC
	var orderRequest OrderRequest
	grpcClient := getOrdergRPCClient()
	grpcPClient := getProductgRPCClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var items []*pb.OrderItem
	for _, item := range orderRequest.Items {
		orderItem := &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		}
		pass, product, err := verifyProduct(orderItem)
		if err != nil {
			log.Printf("could not add the product to the order: %v", err)
		}
		if pass {
			product.Quantity = product.Quantity - int32(item.Quantity)
			r, err := grpcPClient.UpdateProduct(ctx, &pb.UpdateProductRequest{Id: product.Id, Name: product.Name, Description: product.Description, Price: float32(product.Price), Quantity: int32(product.Quantity)})
			if err != nil {
				log.Printf("could not find the product: %v", err)
				return
			}
			log.Printf("product quantity disconunted: %v", r.Name)
			items = append(items, orderItem)
		}

	}

	if len(items) > 0 {
		r, err := grpcClient.PlaceOrder(ctx, &pb.OrderRequest{Items: items})
		if err != nil {
			log.Printf("could not create the order: %v", err)
		}
		log.Printf("Created order with id: %s", r.Id)
		c.JSON(http.StatusOK, gin.H{"message": "Order created successfully!", "order": r.Id})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "could not create the order"})
	}

}

func GetOrders(c *gin.Context) {
	// Obtener las ordenes creadas con gRPC
	grpcClient := getOrdergRPCClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := grpcClient.GetOrders(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("could not read the orders: %v", err)
	}
	//log.Printf("List of products: %v", r.Products)
	if len(r.Orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No orders found"})
	} else {
		c.JSON(http.StatusOK, r.Orders)
	}

}
