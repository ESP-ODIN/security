package implementation

import (
	"context"

	"github.com/google/uuid"

	"security/internal/dto"
	"security/internal/model"
	"security/internal/repository"
	"security/internal/service"
)

type versionService struct {
	repo repository.VersionRepository
}

// NewVersionService returns a VersionService.
func NewVersionService(repo repository.VersionRepository) service.VersionService {
	return &versionService{repo: repo}
}

func (s *versionService) Submit(ctx context.Context, agentID uuid.UUID, idempotencyKey string) (dto.SubmitVersionResponse, error) {
	version := model.Version{AgentID: agentID, Status: "accepted"}
	if _, err := s.repo.Create(ctx, version, idempotencyKey); err != nil {
		return dto.SubmitVersionResponse{}, translate(err)
	}
	return dto.SubmitVersionResponse{
		AgentID:        agentID,
		IdempotencyKey: idempotencyKey,
		Status:         "accepted",
		Message:        "version submitted; security pipeline will start",
	}, nil
}

func (s *versionService) GetPermissions(ctx context.Context, agentID, versionID uuid.UUID) (dto.PermissionsResponse, error) {
	perms, err := s.repo.GetPermissions(ctx, agentID, versionID)
	if err != nil {
		return dto.PermissionsResponse{}, translate(err)
	}
	if perms == nil {
		perms = []model.Permission{}
	}
	return dto.PermissionsResponse{
		AgentID:     agentID,
		VersionID:   versionID,
		Permissions: perms,
	}, nil
}
