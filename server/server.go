// main package to be able to run the studentsServer for now
package main

import (
	"context"
	"log"
	"net"

	gpb "github.com/BetterGR/students-microservice/students_protobuf"
	"google.golang.org/grpc"
)

// studentsServer is an implementation of GRPC Students microservice.
type studentsServer struct {
	// throws unimplemented error
	gpb.UnimplementedStudentsServer
}

// GetStudent search for the Student that corresponds to the given id and returns them.
func (s *studentsServer) GetStudent(ctx context.Context, req *gpb.GetStudentRequest) (*gpb.GetStudentResponse, error) {
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
func (s *studentsServer) CreateStudent(ctx context.Context, req *gpb.CreateStudentRequest) (*gpb.CreateStudentResponse, error) {
	log.Printf("Received CreateStudent request for: %s %s", req.student.firstName, req.student.secondName)
	return &gpb.CreateStudentResponse{Student: req.student}, nil
}

// UpdateStudent updates the given Student and returns them after the update.
func (s *studentsServer) UpdateStudent(ctx context.Context, req *gpb.UpdateStudentRequest) (*gpb.UpdateStudentResponse, error) {
	log.Printf("Received UpdateStudent request for: %s %s", req.student.firstName, req.student.secondName)
	return &gpb.UpdateStudentResponse{Student: req.student}, nil
}

// GetStudentCourses searches the courses that the Student is enrolled in during the given semester and returns them.
func (s *studentsServer) GetStudentCourses(ctx context.Context, req *gpb.GetStudentCoursesRequest) (*gpb.GetStudentCoursesResponse, error) {
	log.Printf("Received GetStudentCourses request for %s %s", req.student.firstName, req.student.secondName,
		"in semester %s", req.semester)
	return &gpb.GetStudentCoursesResponse{
		Courses: []*gpb.Course{
			{Id: "C1", Name: "Mathematics", Semester: "Spring 2024"},
			{Id: "C2", Name: "Physics", Semester: "Spring 2024"},
		},
	}, nil
}

// GetStudentGrades searches the course that corresponds to the given course_id in the given semester
// and returns the students grades in this course.
func (s *studentsServer) GetStudentGrades(ctx context.Context, req *gpb.GetStudentGradesRequest) (*gpb.GetStudentGradesResponse, error) {
	log.Printf("Received GetStudentGrades request for  %s %s", req.student.firstName, req.student.secondName,
		"in %s, semester %s", req.course_id, req.semester)
	return &gpb.GetStudentGradesResponse{
		Grades: []*gpb.Grade{
			{CourseId: "C1", Grade: "A"},
			{CourseId: "C2", Grade: "B"},
		},
	}, nil
}

// DeleteStudent deletes the Student from the system.
func (s *studentsServer) DeleteStudent(ctx context.Context, req *gpb.DeleteStudentRequest) (*gpb.DeleteStudentResponse, error) {
	log.Printf("Received DeleteStudent request for ID: %s", req.Id)
	// Delete the Student
	log.Printf("Student with ID: %s has been deleted", req.Id)
	return &gpb.DeleteStudentResponse{}, nil
}

const (
	// define port
	address = "localhost:50051"
)

// main studentsServer function
func main() {
	// create a listener on port 'address'
	lis, err := net.Listen("tcp", address)
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
