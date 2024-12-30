// Copyright (C) 2024 Lukas Rotermund
// See end of file for extended copyright information.

package db

import (
	"database/sql"
	"embed"
	"fmt"
	"math/rand/v2"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	"github.com/lrotermund/golactic-union/internal/logger"
	"github.com/lrotermund/golactic-union/internal/models"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var (
	spaceShipTypes = []string{
		"SCX",
		"Nimbus",
		"Arrow",
		"Fighter",
		"Freighter",
		"Frodo",
		"Shadow",
	}

	generations = []string{
		"Alpha",
		"Beta",
		"Gamma",
		"Version",
		"Gen",
	}
)

func New(development bool) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./union.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate() error {
	db, err := sql.Open("sqlite", "./union.db")
	if err != nil {
		return err
	}

	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationsFS,
		Root:       "migrations",
	}

	if _, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up); err != nil {
		return err
	}

	return nil
}

func Seed() error {
	db, err := sql.Open("sqlite", "./union.db")
	if err != nil {
		return err
	}

	sqlCode := "INSERT OR IGNORE INTO space_ships (id, name) VALUES"
	spaceShips := make([]models.SpaceShip, 100)

	for i := range 100 {
		spaceShips[i] = models.SpaceShip{
			ID:   uuid.Must(uuid.NewV7()),
			Name: spaceShipName(),
		}

		if i+i < 100 {
			sqlCode += " (?, ?),"
		} else {
			sqlCode += " (?, ?)"
		}
	}

	_, err = db.Exec(sqlCode)
	if err != nil {
		logger.Error("failed to seed database", zap.Error(err))

		return err
	}

	logger.Info("seeded database")

	return nil
}

func spaceShipName() string {
	return fmt.Sprintf(
		"%s %s-%d",
		spaceShipTypes[rand.IntN(len(spaceShipTypes))],
		generations[rand.IntN(len(generations))],
		rand.IntN(2000),
	)
}

func Mock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		logger.Fatal("failed to create mock db", zap.Error(err))
	}

	return db, mock
}

// golactic-union is the golang equivalent of a golang (echo) vs PHP (Symfony)
// benchmark project.
// Copyright (C) 2024 Lukas Rotermund
//
// This file is part of golactic-union.
//
// golactic-union is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// golactic-union is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with golactic-union. If not, see <https://www.gnu.org/licenses/>.
