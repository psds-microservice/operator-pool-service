package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/psds-microservice/operator-pool-service/internal/errs"
	"github.com/psds-microservice/operator-pool-service/internal/model"
	"github.com/psds-microservice/operator-pool-service/internal/service"
	"github.com/psds-microservice/operator-pool-service/pkg/gen/operator_pool_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Deps — зависимости gRPC-сервера (D: зависимость от абстракций).
type Deps struct {
	Operator service.OperatorServicer
}

// Server implements operator_pool_service.OperatorPoolServiceServer
type Server struct {
	operator_pool_service.UnimplementedOperatorPoolServiceServer
	Deps
}

// NewServer создаёт gRPC-сервер с внедрёнными сервисами
func NewServer(deps Deps) *Server {
	return &Server{Deps: deps}
}

func (s *Server) mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, errs.ErrNoOperatorAvailable) {
		return status.Error(codes.NotFound, err.Error())
	}
	log.Printf("grpc: unhandled error: %v", err)
	return status.Error(codes.Internal, err.Error())
}

func toProtoOperatorStatus(op *model.OperatorStatus) *operator_pool_service.OperatorStatus {
	if op == nil {
		return nil
	}
	return &operator_pool_service.OperatorStatus{
		UserId:         op.UserID.String(),
		Available:      op.Available,
		ActiveSessions: int32(op.ActiveSessions),
		MaxSessions:    int32(op.MaxSessions),
	}
}

func (s *Server) SetStatus(ctx context.Context, req *operator_pool_service.SetStatusRequest) (*operator_pool_service.SetStatusResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	maxSessions := int(req.GetMaxSessions())
	if maxSessions <= 0 {
		maxSessions = 5 // default
	}
	if err := s.Operator.SetStatus(userID, req.GetAvailable(), maxSessions); err != nil {
		return nil, s.mapError(err)
	}
	return &operator_pool_service.SetStatusResponse{Ok: true}, nil
}

func (s *Server) GetNext(ctx context.Context, req *operator_pool_service.GetNextRequest) (*operator_pool_service.GetNextResponse, error) {
	operatorID, err := s.Operator.Next()
	if err != nil {
		return nil, s.mapError(err)
	}
	return &operator_pool_service.GetNextResponse{OperatorId: operatorID.String()}, nil
}

func (s *Server) GetStats(ctx context.Context, req *operator_pool_service.GetStatsRequest) (*operator_pool_service.GetStatsResponse, error) {
	available, total, err := s.Operator.Stats()
	if err != nil {
		return nil, s.mapError(err)
	}
	return &operator_pool_service.GetStatsResponse{
		Available: int32(available),
		Total:     int32(total),
	}, nil
}

func (s *Server) ListOperators(ctx context.Context, req *operator_pool_service.ListOperatorsRequest) (*operator_pool_service.ListOperatorsResponse, error) {
	list, err := s.Operator.ListAll()
	if err != nil {
		return nil, s.mapError(err)
	}
	operators := make([]*operator_pool_service.OperatorStatus, len(list))
	for i, op := range list {
		operators[i] = toProtoOperatorStatus(&op)
	}
	return &operator_pool_service.ListOperatorsResponse{
		Operators: operators,
	}, nil
}
