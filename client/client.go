package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "github.com/BetterGR/students-microservice/students_protobuf"
)

func main() {
	// Establish connection with the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewStudentsServiceClient(conn)

	// Example: Call GetStudent
	getStudent(client)

	// Example: Call CreateStudent
	createStudent(client)
}

func getStudent(client pb.StudentsServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.GetStudent(ctx, &pb.GetStudentRequest{
		Token: "example-token",
		Id:    "123",
	})
	if err != nil {
		log.Fatalf("Error calling GetStudent: %v", err)
	}

	log.Printf("GetStudent Response: %+v", response.Student)
}

func createStudent(client pb.StudentsServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newStudent := &pb.Student{
		FirstName:  "Jane",
		SecondName: "Doe",
		Id:         "124",
	}

	response, err := client.CreateStudent(ctx, &pb.CreateStudentRequest{
		Token:  "example-token",
		Student: newStudent,
	})
	if err != nil {
		log.Fatalf("Error calling CreateStudent: %v", err)
	}

	log.Printf("CreateStudent Response: %+v", response.Student)
}
