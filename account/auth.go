package account

import (
	"context"

	accountpb "github.com/ihtkas/farm/account/v1"

	"google.golang.org/grpc"
)

// This file has simple prototype to mimic user validation. Will be replaced by some open source IAM implementation

// ValidateUser checks the existence of the user in the database
func (m *Manager) ValidateUser(ctx context.Context, in *accountpb.ValidateUserRequest,
	opts ...grpc.CallOption) (*accountpb.ValidateUserResponse, error) {
	return &accountpb.ValidateUserResponse{}, m.store.ValidateUser(in)
}
