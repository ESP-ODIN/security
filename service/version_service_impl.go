package service

import (
	"context"

	"security/dto"
	"security/model"
	"security/repository"
)

type versionService struct {
	repo repository.VersionRepository
}

// NewVersionService returns a VersionService.
func NewVersionService(repo repository.VersionRepository) VersionService {
	return &versionService{repo: repo}
}

func (s *versionService) Submit(ctx context.Context, agentID, idempotencyKey string) (dto.SubmitVersionResponse, error) {
	version := model.Version{AgentID: agentID, Status: "accepted"}
	if _, err := s.repo.Create(ctx, version, idempotencyKey); err != nil {
		return dto.SubmitVersionResponse{}, err
	}
	return dto.SubmitVersionResponse{
		AgentID:        agentID,
		IdempotencyKey: idempotencyKey,
		Status:         "accepted",
		Message:        "version submitted; security pipeline will start",
	}, nil
}

func (s *versionService) GetPermissions(ctx context.Context, agentID, versionID string) (dto.PermissionsResponse, error) {
	perms, err := s.repo.GetPermissions(ctx, agentID, versionID)
	if err != nil {
		return dto.PermissionsResponse{}, err
	}
	items := make([]any, 0, len(perms))
	for _, p := range perms {
		items = append(items, p)
	}
	return dto.PermissionsResponse{
		AgentID:     agentID,
		VersionID:   versionID,
		Permissions: items,
	}, nil
}
