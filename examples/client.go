package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	gpb "github.com/BetterGR/students-microservice/students_protobuf"
)

const (
	// define port
	address = "localhost:50051"
)

func main() {
	// Establish connection with the gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := gpb.NewStudentsServiceClient(conn)

	// Example: Call GetStudent
	getStudent(client)

	// Example: Call CreateStudent
	createStudent(client)
}

func getStudent(client gpb.StudentsServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.GetStudent(ctx, &gpb.GetStudentRequest{
		Token: "example-token",
		Id:    "123",
	})
	if err != nil {
		log.Fatalf("Error calling GetStudent: %v", err)
	}

	log.Printf("GetStudent Response: %+v", response.Student)
}

func createStudent(client gpb.StudentsServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newStudent := &gpb.Student{
		FirstName:  "Jane",
		SecondName: "Doe",
		Id:         "124",
	}

	response, err := client.CreateStudent(ctx, &gpb.CreateStudentRequest{
		Token:  "example-token",
		Student: newStudent,
	})
	if err != nil {
		log.Fatalf("Error calling CreateStudent: %v", err)
	}

	log.Printf("CreateStudent Response: %+v", response.Student)
}
