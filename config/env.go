package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Config struct {
	HTTPAddr     string
	DatabaseURL  string
	AuthJWKSURL  string
	AuthIssuer   string
	AuthAudience string
	AdminUserIDs []uuid.UUID
}

func Load() (Config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid HTTP_ADDR %q: %w", addr, err)
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		return Config{}, fmt.Errorf("HTTP_ADDR port must be between 1 and 65535")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	adminIDs, err := parseAdminUserIDs(os.Getenv("ADMIN_USER_IDS"))
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		HTTPAddr:     addr,
		DatabaseURL:  dbURL,
		AuthJWKSURL:  os.Getenv("AUTH_JWKS_URL"),
		AuthIssuer:   os.Getenv("AUTH_ISSUER"),
		AuthAudience: os.Getenv("AUTH_AUDIENCE"),
		AdminUserIDs: adminIDs,
	}
	for _, entry := range []struct{ name, value string }{
		{"AUTH_JWKS_URL", cfg.AuthJWKSURL},
		{"AUTH_ISSUER", cfg.AuthIssuer},
		{"AUTH_AUDIENCE", cfg.AuthAudience},
	} {
		if strings.TrimSpace(entry.value) == "" {
			return Config{}, fmt.Errorf("%s is required", entry.name)
		}
	}
	return cfg, nil
}

func parseAdminUserIDs(raw string) ([]uuid.UUID, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	parts := strings.Split(raw, ",")
	ids := make([]uuid.UUID, 0, len(parts))
	seen := make(map[uuid.UUID]struct{}, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("ADMIN_USER_IDS contains an empty value")
		}
		id, err := uuid.Parse(part)
		if err != nil {
			return nil, fmt.Errorf("ADMIN_USER_IDS contains invalid uuid %q", part)
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}
