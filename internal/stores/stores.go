// Copyright (C) 2024 Lukas Rotermund
// See end of file for extended copyright information.

package stores

import (
	"database/sql"
)

type Stores struct {
	DB        *sql.DB
	SpaceShip SpaceShipStore
}

func New(db *sql.DB) *Stores {
	return &Stores{
		DB:        db,
		SpaceShip: &spaceShipStore{db},
	}
}

func (s *Stores) Begin() (*sql.Tx, error) {
	return s.DB.Begin()
}

func (s *Stores) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (s *Stores) RollBack(tx *sql.Tx) error {
	return tx.Rollback()
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
