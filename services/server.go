package services

import (
	pb "api-gateway/api-gateway"
	"context"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type productServer struct {
	pb.UnimplementedProductServiceServer
	mu       sync.Mutex
	products []*pb.ProductResponse
}

type orderServer struct {
	pb.UnimplementedOrderServiceServer
	mu     sync.Mutex
	orders []*pb.OrderResponse
}

func (s *productServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product := &pb.ProductResponse{
		Id:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}

	s.products = append(s.products, product)
	return product, nil
}

func (s *productServer) GetProducts(context.Context, *pb.Empty) (*pb.ProductsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &pb.ProductsResponse{Products: s.products}, nil
}

func (s *productServer) GetProductByID(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, product := range s.products {
		if product.Id == req.Id {
			return product, nil
		}
	}
	return nil, errors.New("Product not found")
}

func (s *productServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for index, product := range s.products {
		if product.Id == req.Id {
			iproduct := &pb.ProductResponse{
				Id:          product.Id,
				Name:        req.Name,
				Description: req.Description,
				Price:       req.Price,
				Quantity:    req.Quantity,
			}
			s.products[index] = iproduct
			return iproduct, nil
		}
	}

	return nil, errors.New("Product not found")
}

func (s *productServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for index, product := range s.products {
		if product.Id == req.Id {
			s.products = append(s.products[:index], s.products[index+1:]...)
			return &pb.Empty{}, nil
		}
	}

	return &pb.Empty{}, errors.New("Product not found")
}

func (s *orderServer) PlaceOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order := &pb.OrderResponse{
		Id:    uuid.New().String(),
		Items: req.Items,
		Total: 0,
	}

	s.orders = append(s.orders, order)
	return order, nil
}

func (s *orderServer) GetOrders(ctx context.Context, req *pb.Empty) (*pb.OrdersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &pb.OrdersResponse{Orders: s.orders}, nil
}

func GrpcServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, &productServer{})
	pb.RegisterOrderServiceServer(grpcServer, &orderServer{})

	log.Println("gRPCServer is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
