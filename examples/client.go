package main

import (
	"context"
	"log"
	"time"

	gpb "github.com/BetterGR/students-microservice/students_protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// define port
	address = "localhost:50052"
)

func main() {
	// Establish connection with the gRPC server
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := gpb.NewStudentsServiceClient(conn)

	req := &gpb.GetStudentRequest{Token: "I am admin", Id: "123"}
	response, err := client.GetStudent(ctx, req)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("FirstName: %s, SecondName: %s, ID: %s", response.Student.GetFirstName(), response.Student.GetSecondName(), response.Student.GetId())
}
