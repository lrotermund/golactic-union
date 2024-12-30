// Copyright (C) 2024 Lukas Rotermund
// See end of file for extended copyright information.

package stores_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lrotermund/golactic-union/internal/db"
	"github.com/lrotermund/golactic-union/internal/logger"
	"github.com/lrotermund/golactic-union/internal/models"
	"github.com/lrotermund/golactic-union/internal/stores"
	"github.com/stretchr/testify/assert"
)

func TestSpaceShipStore_GetSuccessCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	createdAt, _ := time.Parse(time.RFC3339, "2024-12-27T16:44:15Z00:00")
	spaceShips := []models.SpaceShip{
		{
			ID:        uuid.MustParse("5e6cb154-fd72-438a-a63a-bd783774b479"),
			Name:      "test-ship-1",
			Intact:    true,
			CreatedAt: createdAt,
		},
		{
			ID:        uuid.MustParse("3b7dbc34-f22b-44a7-8174-6628e80e408f"),
			Name:      "test-ship-2",
			Intact:    false,
			CreatedAt: createdAt,
		},
	}

	rows := mock.NewRows([]string{"id", "name", "intact", "created_at"})
	for _, s := range spaceShips {
		rows.AddRow(s.ID, s.Name, s.Intact, s.CreatedAt)
	}

	mock.ExpectQuery("SELECT id, name, intact, created_at from spaceShips").
		WillReturnRows(rows)

	s := stores.New(mockDB)

	r, err := s.SpaceShip.Get()

	assert.NoError(t, err)
	assert.Equal(t, spaceShips, r)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSpaceShipStore_GetErrorCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	expectedErr := fmt.Errorf("Something went wrong")

	mock.ExpectQuery("SELECT id, name, intact, created_at from spaceShips").
		WillReturnError(expectedErr)

	s := stores.New(mockDB)

	r, err := s.SpaceShip.Get()

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, r)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSpaceShipStore_CreateSuccessCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	spaceShip := models.SpaceShip{
		ID:   uuid.MustParse("5e6cb154-fd72-438a-a63a-bd783774b479"),
		Name: "test-ship-1",
	}

	mock.ExpectExec("INSERT INTO spaceShips (id, name) VALUES ($1, $2)").
		WithArgs(spaceShip.ID, spaceShip.Name).
		WillReturnResult(sqlmock.NewResult(0, 0))

	s := stores.New(mockDB)

	err := s.SpaceShip.Create(nil, &spaceShip)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSpaceShipStore_CreateOnTxSuccessCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	spaceShip := models.SpaceShip{
		ID:   uuid.MustParse("5e6cb154-fd72-438a-a63a-bd783774b479"),
		Name: "test-ship-1",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO spaceShips (id, name) VALUES ($1, $2)").
		WithArgs(spaceShip.ID, spaceShip.Name).
		WillReturnResult(sqlmock.NewResult(0, 0))

	s := stores.New(mockDB)

	tx, _ := mockDB.Begin()
	err := s.SpaceShip.Create(tx, &spaceShip)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSpaceShipStore_CreateErrorCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	spaceShip := models.SpaceShip{
		ID:   uuid.MustParse("5e6cb154-fd72-438a-a63a-bd783774b479"),
		Name: "test-ship-1",
	}

	expectedErr := fmt.Errorf("Something went wrong")

	mock.ExpectExec("INSERT INTO spaceShips (id, name) VALUES ($1, $2)").
		WithArgs(spaceShip.ID, spaceShip.Name).
		WillReturnResult(nil).
		WillReturnError(expectedErr)

	s := stores.New(mockDB)

	err = s.SpaceShip.Create(nil, &spaceShip)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSpaceShipStore_CreateOnTxErrorCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	spaceShip := models.SpaceShip{
		ID:   uuid.MustParse("5e6cb154-fd72-438a-a63a-bd783774b479"),
		Name: "test-ship-1",
	}

	expectedErr := fmt.Errorf("Something went wrong")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO spaceShips (id, name) VALUES ($1, $2)").
		WithArgs(spaceShip.ID, spaceShip.Name).
		WillReturnResult(nil).
		WillReturnError(expectedErr)

	s := stores.New(mockDB)

	tx, _ := mockDB.Begin()
	err = s.SpaceShip.Create(tx, &spaceShip)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSpaceShipStore_DeleteByIDSuccessCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	spaceShip := models.SpaceShip{
		ID: uuid.MustParse("5e6cb154-fd72-438a-a63a-bd783774b479"),
	}

	mock.ExpectExec("DELETE FROM spaceShips WHERE spaceShips.id = $1 RETURNING spaceShips.id").
		WithArgs(spaceShip.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s := stores.New(mockDB)

	err := s.SpaceShip.DeleteByID(spaceShip.ID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSpaceShipStore_DeleteByIDErrorCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	spaceShip := models.SpaceShip{
		ID: uuid.MustParse("5e6cb154-fd72-438a-a63a-bd783774b479"),
	}

	expectedErr := fmt.Errorf("Something went wrong")

	mock.ExpectExec("DELETE FROM spaceShips WHERE spaceShips.id = $1 RETURNING spaceShips.id").
		WithArgs(spaceShip.ID).
		WillReturnResult(nil).
		WillReturnError(expectedErr)

	s := stores.New(mockDB)

	err = s.SpaceShip.DeleteByID(spaceShip.ID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSpaceShipStore_DeleteByIDErrorResultCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	spaceShip := models.SpaceShip{
		ID: uuid.MustParse("5e6cb154-fd72-438a-a63a-bd783774b479"),
	}

	expectedErr := fmt.Errorf("Something went wrong")

	mock.ExpectExec("DELETE FROM spaceShips WHERE spaceShips.id = $1 RETURNING spaceShips.id").
		WithArgs(spaceShip.ID).
		WillReturnResult(sqlmock.NewErrorResult(expectedErr))

	s := stores.New(mockDB)

	err = s.SpaceShip.DeleteByID(spaceShip.ID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSpaceShipStore_DeleteByIDNoEffecedRowsCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	spaceShip := models.SpaceShip{
		ID: uuid.MustParse("5e6cb154-fd72-438a-a63a-bd783774b479"),
	}

	mock.ExpectExec("DELETE FROM spaceShips WHERE spaceShips.id = $1 RETURNING spaceShips.id").
		WithArgs(spaceShip.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	s := stores.New(mockDB)

	err = s.SpaceShip.DeleteByID(spaceShip.ID)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
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
