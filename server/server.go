// main package to be able to run the StudentsServer for now
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	spb "github.com/BetterGR/students-microservice/protos"
	ms "github.com/TekClinic/MicroService-Lib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

const (
	// define address.
	connectionProtocol = "tcp"
	// Debugging logs.
	logLevelDebug = 5
)

// StudentsServer is an implementation of GRPC Students microservice.
type StudentsServer struct {
	ms.BaseServiceServer
	db *Database
	spb.UnimplementedStudentsServiceServer
	Claims ms.Claims
}

// VerifyToken returns the injected Claims instead of the default.
func (s *StudentsServer) VerifyToken(ctx context.Context, token string) error {
	if s.Claims != nil {
		return nil
	}

	// Default behavior.
	if _, err := s.BaseServiceServer.VerifyToken(ctx, token); err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}

	return nil
}

// initStudentsMicroserviceServer initializes the StudentsServer.
func initStudentsMicroserviceServer() (*StudentsServer, error) {
	base, err := ms.CreateBaseServiceServer()
	if err != nil {
		return nil, fmt.Errorf("failed to create base service: %w", err)
	}

	database, err := InitializeDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &StudentsServer{
		BaseServiceServer:                  base,
		db:                                 database,
		UnimplementedStudentsServiceServer: spb.UnimplementedStudentsServiceServer{},
	}, nil
}

// GetStudent search for the Student that corresponds to the given id and returns them.
func (s *StudentsServer) GetStudent(ctx context.Context,
	req *spb.GetStudentRequest,
) (*spb.GetStudentResponse, error) {
	if err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received GetStudent request", "studentId", req.GetStudentID())

	student, err := s.db.GetStudent(ctx, req.GetStudentID())
	if err != nil {
		return nil, fmt.Errorf("student not found: %w",
			status.Error(codes.NotFound, err.Error()))
	}

	spbStudent := &spb.Student{
		StudentID:   student.StudentID,
		FirstName:   student.FirstName,
		LastName:    student.LastName,
		Email:       student.Email,
		PhoneNumber: student.PhoneNumber,
	}

	return &spb.GetStudentResponse{
		Student: spbStudent,
	}, nil
}

// processStudent processes the student creation or update request.
func (s *StudentsServer) processStudent(ctx context.Context, token string, student *spb.Student,
	action func(context.Context, *spb.Student) (*Student, error),
) (*spb.Student, error) {
	if err := s.VerifyToken(ctx, token); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Processing student request",
		"firstName", student.GetFirstName(), "lastName", student.GetLastName())

	processedStudent, err := action(ctx, student)
	if err != nil {
		return nil, fmt.Errorf("failed to process student: %w", status.Error(codes.Internal, err.Error()))
	}

	return &spb.Student{
		StudentID:   processedStudent.StudentID,
		FirstName:   processedStudent.FirstName,
		LastName:    processedStudent.LastName,
		Email:       processedStudent.Email,
		PhoneNumber: processedStudent.PhoneNumber,
	}, nil
}

// CreateStudent creates a new Student with the given details and returns them.
func (s *StudentsServer) CreateStudent(ctx context.Context,
	req *spb.CreateStudentRequest,
) (*spb.CreateStudentResponse, error) {
	spbStudent, err := s.processStudent(ctx, req.GetToken(), req.GetStudent(), s.db.AddStudent)
	if err != nil {
		return nil, err
	}

	return &spb.CreateStudentResponse{Student: spbStudent}, nil
}

// UpdateStudent updates the given Student and returns them after the update.
func (s *StudentsServer) UpdateStudent(ctx context.Context,
	req *spb.UpdateStudentRequest,
) (*spb.UpdateStudentResponse, error) {
	spbStudent, err := s.processStudent(ctx, req.GetToken(), req.GetStudent(), s.db.UpdateStudent)
	if err != nil {
		return nil, err
	}

	return &spb.UpdateStudentResponse{Student: spbStudent}, nil
}

// DeleteStudent deletes the Student from the system.
func (s *StudentsServer) DeleteStudent(ctx context.Context,
	req *spb.DeleteStudentRequest,
) (*spb.DeleteStudentResponse, error) {
	if err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received DeleteStudent request", "studentId", req.GetStudentID())

	if err := s.db.DeleteStudent(ctx, req.GetStudentID()); err != nil {
		return nil, fmt.Errorf("failed to delete student: %w", status.Error(codes.Internal, err.Error()))
	}

	logger.V(logLevelDebug).Info("Deleted", "studentId", req.GetStudentID())

	return &spb.DeleteStudentResponse{}, nil
}

// main StudentsServer function.
func main() {
	// init klog
	klog.InitFlags(nil)
	flag.Parse()

	if err := godotenv.Load(); err != nil {
	    klog.Warning("Warning: No .env file loaded, proceeding with environment variables only")
	}

	// init the StudentsServer
	server, err := initStudentsMicroserviceServer()
	if err != nil {
		klog.Fatalf("Failed to init StudentsServer: %v", err)
	}

	// create a listener on port 'address'
	address := "localhost:" + os.Getenv("GRPC_PORT")

	lis, err := net.Listen(connectionProtocol, address)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	klog.V(logLevelDebug).Info("Starting StudentsServer on port: ", address)
	// create a grpc StudentsServer
	grpcServer := grpc.NewServer()
	spb.RegisterStudentsServiceServer(grpcServer, server)

	// serve the grpc StudentsServer
	if err := grpcServer.Serve(lis); err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
