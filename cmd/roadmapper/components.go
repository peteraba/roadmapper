package main

import (
	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/migrations"
	"github.com/peteraba/roadmapper/pkg/repository"
	"github.com/peteraba/roadmapper/pkg/roadmap"
	"go.uber.org/zap"
)

// newLogger DON'T FORGET TO CALL logger.Sync() !!!!
func newLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger
}

func newRoadmapRepo(dbHost, dbPort, dbName, dbUser, dbPass string, logger *zap.Logger) roadmap.Repository {
	baseRepo := repository.NewPgRepository(AppName, dbHost, dbPort, dbName, dbUser, dbPass, logger)
	return roadmap.Repository{PgRepository: baseRepo}
}

func newCodeBuilder() code.Builder {
	return code.Builder{}
}

func newMigrations(dbHost, dbPort, dbName, dbUser, dbPass string) *migrations.Migrations {
	return migrations.New(dbUser, dbPass, dbHost, dbPort, dbName)
}
