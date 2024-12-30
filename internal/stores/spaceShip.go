// Copyright (C) 2024 Lukas Rotermund
// See end of file for extended copyright information.

package stores

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/lrotermund/golactic-union/internal/logger"
	"github.com/lrotermund/golactic-union/internal/models"
	"go.uber.org/zap"
)

type (
	SpaceShipStore interface {
		Get() ([]models.SpaceShip, error)
		Create(tx *sql.Tx, spaceShip *models.SpaceShip) error
		DeleteByID(id uuid.UUID) error
	}

	spaceShipStore struct {
		*sql.DB
	}
)

func (s *spaceShipStore) Get() ([]models.SpaceShip, error) {
	spaceShips := []models.SpaceShip{}
	rows, err := s.Query("SELECT id, name, intact, created_at from spaceShips")

	if err != nil {
		logger.Error("failed to get spaceShip", zap.Error(err))

		return nil, err
	}

	for rows.Next() {
		var s models.SpaceShip
		err = rows.Scan(&s.ID, &s.Name, &s.Intact, &s.CreatedAt)
		spaceShips = append(spaceShips, s)
	}

	return spaceShips, nil
}

func (s *spaceShipStore) Create(tx *sql.Tx, spaceShip *models.SpaceShip) error {
	var query string
	query = "INSERT INTO spaceShips (id, name) VALUES ($1, $2)"
	var err error

	if tx != nil {
		_, err = tx.Exec(query, spaceShip.ID, spaceShip.Name)
	} else {
		_, err = s.Exec(query, spaceShip.ID, spaceShip.Name)
	}

	if err != nil {
		logger.Error("failed to create spaceShip", zap.Error(err))

		return err
	}

	return nil
}

func (s *spaceShipStore) DeleteByID(id uuid.UUID) error {
	row, err := s.Exec("DELETE FROM spaceShips WHERE spaceShips.id = $1 RETURNING spaceShips.id", id)
	if err != nil {
		logger.Error("failed to delete spaceShip by id", zap.Error(err))

		return err
	}

	if r, err := row.RowsAffected(); err != nil {
		return err
	} else if r == 0 {
		return sql.ErrNoRows
	}

	return nil
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
