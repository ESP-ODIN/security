package service

import (
	"context"

	"security/dto"
)

// VersionService handles version submission and permissions.
type VersionService interface {
	Submit(ctx context.Context, agentID, idempotencyKey string) (dto.SubmitVersionResponse, error)
	GetPermissions(ctx context.Context, agentID, versionID string) (dto.PermissionsResponse, error)
}
