// main package to be able to run the StudentsServer for now
package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	spb "github.com/BetterGR/students-microservice/protos"
	ms "github.com/TekClinic/MicroService-Lib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

const (
	// define address.
	address            = "localhost:50052"
	connectionProtocol = "tcp"

	// Debugging logs.
	logLevelDebug = 5
)

// StudentsServer is an implementation of GRPC Students microservice.
type StudentsServer struct {
	ms.BaseServiceServer
	// throws unimplemented error
	spb.UnimplementedStudentsServiceServer
}

func initStudentsMicroserviceServer() (*StudentsServer, error) {
	base, err := ms.CreateBaseServiceServer()
	if err != nil {
		return nil, fmt.Errorf("failed to create base service: %w", err)
	}

	return &StudentsServer{
		BaseServiceServer:                  base,
		UnimplementedStudentsServiceServer: spb.UnimplementedStudentsServiceServer{},
	}, nil
}

// GetStudent search for the Student that corresponds to the given id and returns them.
func (s *StudentsServer) GetStudent(ctx context.Context,
	req *spb.GetStudentRequest,
) (*spb.GetStudentResponse, error) {
	_, err := s.VerifyToken(ctx, req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received GetStudent request", "studentId", req.GetId())

	courses := []*spb.Course{
		{Id: "C1", Name: "Mathematics", Semester: "S24"},
		{Id: "C2", Name: "Physics", Semester: "S24"},
	}

	student := &spb.Student{
		FirstName:  "Rick",
		SecondName: "Roll",
		Id:         req.GetId(),
		Courses:    courses,
	}

	return &spb.GetStudentResponse{
		Student: student,
	}, nil
}

// CreateStudent creates a new Student with the given details and returns them.
func (s *StudentsServer) CreateStudent(ctx context.Context,
	req *spb.CreateStudentRequest,
) (*spb.CreateStudentResponse, error) {
	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received CreateStudent request",
		"firstName", req.GetStudent().GetFirstName(), "secondName", req.GetStudent().GetSecondName())

	return &spb.CreateStudentResponse{Student: req.GetStudent()}, nil
}

// UpdateStudent updates the given Student and returns them after the update.
func (s *StudentsServer) UpdateStudent(ctx context.Context,
	req *spb.UpdateStudentRequest,
) (*spb.UpdateStudentResponse, error) {
	_, err := s.VerifyToken(ctx, req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received UpdateStudent request",
		"firstName", req.GetStudent().GetFirstName(), "secondName", req.GetStudent().GetSecondName())

	return &spb.UpdateStudentResponse{Student: req.GetStudent()}, nil
}

// GetStudentCourses searches the courses that the Student is enrolled in during the given semester and returns them.
func (s *StudentsServer) GetStudentCourses(ctx context.Context,
	req *spb.GetStudentCoursesRequest,
) (*spb.GetStudentCoursesResponse, error) {
	_, err := s.VerifyToken(ctx, req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received GetStudentCourses request",
		"firstName", req.GetStudent().GetFirstName(), "secondName", req.GetStudent().GetSecondName(),
		"semester", req.GetSemester())

	courses := []*spb.Course{
		{Id: "C1", Name: "Mathematics", Semester: "S24"},
		{Id: "C2", Name: "Physics", Semester: "S24"},
	}

	return &spb.GetStudentCoursesResponse{
		Courses: courses,
	}, nil
}

// GetStudentGrades searches the course that corresponds to the given course_id in the given semester
// and returns the students grades in this course.
func (s *StudentsServer) GetStudentGrades(ctx context.Context,
	req *spb.GetStudentGradesRequest,
) (*spb.GetStudentGradesResponse, error) {
	_, err := s.VerifyToken(ctx, req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received GetStudentGrades request",
		"firstName", req.GetStudent().GetFirstName(), "secondName", req.GetStudent().GetSecondName(),
		"courseId", req.GetCourseId(), "semester", req.GetSemester())

	grades := []*spb.Grade{
		{Semester: "S24", CourseId: "C1", Grade: "100"},
		{Semester: "S24", CourseId: "C2", Grade: "98"},
	}

	return &spb.GetStudentGradesResponse{
		Grades: grades,
	}, nil
}

// DeleteStudent deletes the Student from the system.
func (s *StudentsServer) DeleteStudent(ctx context.Context,
	req *spb.DeleteStudentRequest,
) (*spb.DeleteStudentResponse, error) {
	_, err := s.VerifyToken(ctx, req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received DeleteStudent request", "studentId", req.GetStudent().GetId())

	logger.Info("Deleted", "studentId", req.GetStudent().GetId())

	return &spb.DeleteStudentResponse{}, nil
}

// main StudentsServer function.
func main() {
	// init klog
	klog.InitFlags(nil)
	flag.Parse()

	// init the StudentsServer
	server, err := initStudentsMicroserviceServer()
	if err != nil {
		klog.Error("Failed to init StudentsServer", err)
	}

	// create a listener on port 'address'
	lis, err := net.Listen(connectionProtocol, address)
	if err != nil {
		klog.Error("Failed to listen:", err)
	}

	// create a grpc StudentsServer
	grpcServer := grpc.NewServer()
	spb.RegisterStudentsServiceServer(grpcServer, server)

	// serve the grpc StudentsServer
	if err := grpcServer.Serve(lis); err != nil {
		klog.Error("Failed to serve:", err)
	}
}
