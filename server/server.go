// main package to be able to run the StudentsServiceServer for now
package main

import (
	"context"
	"log"
	"net"

	gpb "github.com/BetterGR/students-microservice/students_protobuf"
	"google.golang.org/grpc"
)

const (
	// define port
	address = "localhost:50052"
)

// StudentsServiceServer is an implementation of GRPC Students microservice.
type StudentsServiceServer struct {
	// throws unimplemented error
	gpb.UnimplementedStudentsServiceServer
}

// GetStudent search for the Student that corresponds to the given id and returns them.
func (s *StudentsServiceServer) GetStudent(ctx context.Context, req *gpb.GetStudentRequest) (*gpb.GetStudentResponse, error) {
	log.Printf("Received GetStudent request for ID: %s", req.Id)
	return &gpb.GetStudentResponse{
		Student: &gpb.Student{
			FirstName:  "Rick",
			SecondName: "Roll",
			Id:         req.Id,
		},
	}, nil
}

// CreateStudent creates a new Student with the given details and returns them.
func (s *StudentsServiceServer) CreateStudent(ctx context.Context, req *gpb.CreateStudentRequest) (*gpb.CreateStudentResponse, error) {
	log.Printf("Received CreateStudent request for: %s %s", req.Student.FirstName, req.Student.SecondName)
	return &gpb.CreateStudentResponse{Student: req.Student}, nil
}

// UpdateStudent updates the given Student and returns them after the update.
func (s *StudentsServiceServer) UpdateStudent(ctx context.Context, req *gpb.UpdateStudentRequest) (*gpb.UpdateStudentResponse, error) {
	log.Printf("Received UpdateStudent request for: %s %s", req.Student.FirstName, req.Student.SecondName)
	return &gpb.UpdateStudentResponse{Student: req.Student}, nil
}

// GetStudentCourses searches the courses that the Student is enrolled in during the given semester and returns them.
func (s *StudentsServiceServer) GetStudentCourses(ctx context.Context, req *gpb.GetStudentCoursesRequest) (*gpb.GetStudentCoursesResponse, error) {
	log.Printf("Received GetStudentCourses request for %s %s", req.Student.FirstName, req.Student.SecondName,
		"in semester %s", req.Semester)
	return &gpb.GetStudentCoursesResponse{
		Courses: []*gpb.Course{
			{Id: "C1", Name: "Mathematics", Semester: "Spring 2024"},
			{Id: "C2", Name: "Physics", Semester: "Spring 2024"},
		},
	}, nil
}

// GetStudentGrades searches the course that corresponds to the given course_id in the given semester
// and returns the students grades in this course.
func (s *StudentsServiceServer) GetStudentGrades(ctx context.Context, req *gpb.GetStudentGradesRequest) (*gpb.GetStudentGradesResponse, error) {
	log.Printf("Received GetStudentGrades request for  %s %s", req.Student.FirstName, req.Student.SecondName,
		"in %s, semester %s", req.GetCourseId(), req.Semester)
	return &gpb.GetStudentGradesResponse{
		Grades: []*gpb.Grade{
			{Semester: "S24", Id: "C1", Grade: "A"},
			{Semester: "W30", Id: "C2", Grade: "B"},
		},
	}, nil
}

// DeleteStudent deletes the Student from the system.
func (s *StudentsServiceServer) DeleteStudent(ctx context.Context, req *gpb.DeleteStudentRequest) (*gpb.DeleteStudentResponse, error) {
	log.Printf("Received DeleteStudent request for ID: %s", req.Student.GetId())
	// Delete the Student
	log.Printf("Student with ID: %s has been deleted", req.Student.GetId())
	return &gpb.DeleteStudentResponse{}, nil
}

// main StudentsServiceServer function
func main() {
	// create a listener on port 'address'
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// create a grpc StudentsServiceServer
	grpcServer := grpc.NewServer()
	gpb.RegisterStudentsServiceServer(grpcServer, &StudentsServiceServer{})

	// serve the grpc StudentsServiceServer
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
