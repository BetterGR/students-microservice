// main package to be able to run the StudentsServer for now
package main

import (
	"context"
	"net"

	spb "github.com/BetterGR/students-microservice/students_protobuf"
	"google.golang.org/grpc"

	"k8s.io/klog/v2"
)

const (
	// define port
	address            = "localhost:50052"
	connectionProtocol = "tcp"
)

// StudentsServer is an implementation of GRPC Students microservice.
type StudentsServer struct {
	// throws unimplemented error
	spb.UnimplementedStudentsServiceServer
}

// GetStudent search for the Student that corresponds to the given id and returns them.
func (s *StudentsServer) GetStudent(ctx context.Context,
	req *spb.GetStudentRequest) (*spb.GetStudentResponse, error) {
	logger := klog.FromContext(ctx)
	logger.Info("Received GetStudent request for:", "student_id", req.Id)

	student := &spb.Student{
		FirstName:  "Rick",
		SecondName: "Roll",
		Id:         req.Id,
	}

	return &spb.GetStudentResponse{
		Student: student,
	}, nil
}

// CreateStudent creates a new Student with the given details and returns them.
func (s *StudentsServer) CreateStudent(ctx context.Context,
	req *spb.CreateStudentRequest) (*spb.CreateStudentResponse, error) {
	logger := klog.FromContext(ctx)
	logger.Info("Received CreateStudent request for:",
		"student_firstName", req.Student.FirstName, "student_secondName", req.Student.SecondName)

	return &spb.CreateStudentResponse{Student: req.Student}, nil
}

// UpdateStudent updates the given Student and returns them after the update.
func (s *StudentsServer) UpdateStudent(ctx context.Context,
	req *spb.UpdateStudentRequest) (*spb.UpdateStudentResponse, error) {
	logger := klog.FromContext(ctx)
	logger.Info("Received UpdateStudent request for:",
		"student_firstName", req.Student.FirstName, "student_secondName", req.Student.SecondName)

	return &spb.UpdateStudentResponse{Student: req.Student}, nil
}

// GetStudentCourses searches the courses that the Student is enrolled in during the given semester and returns them.
func (s *StudentsServer) GetStudentCourses(ctx context.Context,
	req *spb.GetStudentCoursesRequest) (*spb.GetStudentCoursesResponse, error) {
	logger := klog.FromContext(ctx)
	logger.Info("Received GetStudentCourses request for:",
		"student_firstName", req.Student.FirstName, "student_secondName", req.Student.SecondName,
		"semester", req.Semester)

	courses := []*spb.Course{
		{Id: "C1", Name: "Mathematics", Semester: "Spring 2024"},
		{Id: "C2", Name: "Physics", Semester: "Spring 2024"},
	}

	return &spb.GetStudentCoursesResponse{
		Courses: courses,
	}, nil
}

// GetStudentGrades searches the course that corresponds to the given course_id in the given semester
// and returns the students grades in this course.
func (s *StudentsServer) GetStudentGrades(ctx context.Context,
	req *spb.GetStudentGradesRequest) (*spb.GetStudentGradesResponse, error) {
	logger := klog.FromContext(ctx)
	logger.Info("Received GetStudentGrades request for:",
		"student_firstName", req.Student.FirstName, "student_secondName", req.Student.SecondName,
		"course", req.CourseId, "semester", req.Semester)

	grades := []*spb.Grade{
		{Semester: "S24", Id: "C1", Grade: "A"},
		{Semester: "W30", Id: "C2", Grade: "B"},
	}

	return &spb.GetStudentGradesResponse{
		Grades: grades,
	}, nil
}

// DeleteStudent deletes the Student from the system.
func (s *StudentsServer) DeleteStudent(ctx context.Context,
	req *spb.DeleteStudentRequest) (*spb.DeleteStudentResponse, error) {
	logger := klog.FromContext(ctx)
	logger.Info("Received DeleteStudent request for ID:", "student_id",req.Student.GetId())
	
	logger.Info("student was deleted")
	return &spb.DeleteStudentResponse{}, nil
}

// main StudentsServer function
func main() {
	// init klog
	klog.InitFlags(nil)
	// create a listener on port 'address'
	lis, err := net.Listen(connectionProtocol, address)
	if err != nil {
		klog.Error("Failed to listen:", err)
	}

	// create a grpc StudentsServer
	grpcServer := grpc.NewServer()
	spb.RegisterStudentsServiceServer(grpcServer, &StudentsServer{})

	// serve the grpc StudentsServer
	if err := grpcServer.Serve(lis); err != nil {
		klog.Error("Failed to serve:", err)
	}
}
