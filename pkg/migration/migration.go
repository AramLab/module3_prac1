package migration

import (
	"fmt"
	"github.com/AramLab/module3_prac1/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
)

func RunMigrations(cfg config.PostgreSQL) error {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode)
	m, err := migrate.New("file://./pkg/migration/scripts", dbURL)
	if err != nil {
		return errors.Wrap(err, "failed to create migration instance")
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "failed to run migrations")
	}
	return nil
}
