package main

import (
	"context"
	"userm-grpc/proto"

	"fmt"
	"log"
	"net"
	"time"
	"userm-grpc/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedUserServiceServer
}

var c *cache.Cache
var db *gorm.DB

const dbPath = "host=localhost port=5432 user=postgres dbname=userdb password=password sslmode=disable"

func init() {
	var e error
	db, e = gorm.Open("postgres", dbPath)
	defer db.Close()

	if e != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&models.Users{})
}

func (*server) GetList(ctx context.Context, r *proto.UserRequest) (*proto.UserResponse, error) {
	id := r.GetId()

	fmt.Println("Checking cache...")
	userVal, found := c.Get(id)
	if found {
		user := userVal.(*proto.User)
		return &proto.UserResponse{
			Result: user,
		}, nil
	} else {
		fmt.Println("Checking database...")
		db, e := gorm.Open("postgres", dbPath)
		defer db.Close()
		if e != nil {
			panic(fmt.Sprintf("failed to connect to database: %v", e))
		}

		var p models.Users
		if e = db.First(&p, id).Error; e != nil {
			return nil, fmt.Errorf("could not find player with id: %v", id)
		} else {
			gRPCResult := p.GetgRPCModel()
			c.Set(id, &gRPCResult, cache.DefaultExpiration)

			return &proto.UserResponse{
				Result: &gRPCResult,
			}, nil
		}
	}

}
func main() {

	fmt.Println("Starting gRPC micro-service...")
	c = cache.New(60*time.Minute, 70*time.Minute)
	c.Set("2", &proto.User{
		Id:      2,
		Fname:   "John",
		City:    "Mirjapur",
		Phone:   9873798097,
		Height:  13,
		Married: true,
	}, cache.DefaultExpiration)

	l, e := net.Listen("tcp", ":50051")
	if e != nil {
		log.Fatalf("Failed to start listener %v", e)
	}

	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, &server{})
	reflection.Register(s)
	if e := s.Serve(l); e != nil {
		log.Fatalf("failed to serve %v", e)
	}
}
