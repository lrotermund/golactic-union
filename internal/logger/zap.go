// Copyright (C) 2024 Lukas Rotermund
// See end of file for extended copyright information.

package logger

import "go.uber.org/zap"

var zapLogger *zap.Logger

func New() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	zapLogger = logger

	return nil
}

func Debug(msg string, fields ...zap.Field) {
	zapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	zapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	zapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	zapLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	zapLogger.Fatal(msg, fields...)
}

func Sync() {
	_ = zapLogger.Sync()
}

func Delete() {
	zapLogger = nil
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
