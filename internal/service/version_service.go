package service

import (
	"context"

	"github.com/google/uuid"

	"security/internal/dto"
)

// VersionService handles version submission and permissions.
type VersionService interface {
	Submit(ctx context.Context, agentID uuid.UUID, idempotencyKey string) (dto.SubmitVersionResponse, error)
	GetPermissions(ctx context.Context, agentID, versionID uuid.UUID) (dto.PermissionsResponse, error)
}
