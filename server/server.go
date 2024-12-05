// main package to be able to run the studentsServer for now
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	gpb "github.com/BetterGR/students-microservice/students_protobuf"
)

type studentsServer struct {
	// throws unimplemented error
	gpb.UnimplementedStudentsServer
}

// GetStudent implementation
func (s *studentsServer) GetStudent(ctx context.Context, req *gpb.GetStudentRequest) (*gpb.GetStudentResponse, error) {
	log.Printf("Received GetStudent request for ID: %s", req.Id)
	return &gpb.GetStudentResponse{
		Student: &gpb.Student{
			FirstName:  "John",
			SecondName: "Doe",
			Id:         req.Id,
		},
	}, nil
}

// CreateStudent implementation
func (s *studentsServer) CreateStudent(ctx context.Context, req *gpb.CreateStudentRequest) (*gpb.CreateStudentResponse, error) {
	log.Printf("Received CreateStudent request for: %s %s", req.Student.FirstName, req.Student.SecondName)
	return &gpb.CreateStudentResponse{Student: req.Student}, nil
}

// UpdateStudent implementation
func (s *studentsServer) UpdateStudent(ctx context.Context, req *gpb.UpdateStudentRequest) (*gpb.UpdateStudentResponse, error) {
	log.Printf("Received UpdateStudent request for ID: %s", req.Id)
	return &gpb.UpdateStudentResponse{Student: req.Student}, nil
}

// GetStudentCourses implementation
func (s *studentsServer) GetStudentCourses(ctx context.Context, req *gpb.GetStudentCoursesRequest) (*gpb.GetStudentCoursesResponse, error) {
	log.Printf("Received GetStudentCourses request for ID: %s", req.Id)
	return &gpb.GetStudentCoursesResponse{
		Courses: []*gpb.Course{
			{Id: "C1", Name: "Mathematics", Semester: "Fall 2024"},
			{Id: "C2", Name: "Physics", Semester: "Spring 2024"},
		},
	}, nil
}

// GetStudentGrades implementation
func (s *studentsServer) GetStudentGrades(ctx context.Context, req *gpb.GetStudentGradesRequest) (*gpb.GetStudentGradesResponse, error) {
	log.Printf("Received GetStudentGrades request for ID: %s", req.Id)
	return &gpb.GetStudentGradesResponse{
		Grades: []*gpb.Grade{
			{CourseId: "C1", Grade: "A"},
			{CourseId: "C2", Grade: "B"},
		},
	}, nil
}

// main studentsServer function
func main() {
	// create a listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// create a grpc studentsServer
	grpcServer := grpc.NewServer()
	gpb.RegisterStudentsServer(grpcServer, &studentsServer{})

	// serve the grpc studentsServer
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}