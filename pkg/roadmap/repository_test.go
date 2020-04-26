// +build integration

package roadmap

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-pg/pg"
	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/testutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

const (
	iDbHost = "localhost"
	iDbName = "rdmp"
	iDbUser = "rdmp"
	iDbPass = "secret"
)

func TestIntegration_pgRepository_InTx(t *testing.T) {
	type args struct {
		operation func(tx *pg.Tx) error
	}
	tests := []struct {
		name    string
		args    args
		lines   []string
		wantErr bool
	}{
		{
			"empty",
			args{
				operation: func(tx *pg.Tx) error {
					return nil
				},
			},
			[]string{"BEGIN", "COMMIT"},
			false,
		},
		{
			"simple select",
			args{
				operation: func(tx *pg.Tx) error {
					_, err := tx.ExecOne("SELECT 2")

					return err
				},
			},
			[]string{"BEGIN", "SELECT 2", "COMMIT"},
			false,
		},
		{
			"error",
			args{
				operation: func(tx *pg.Tx) error {
					_, err := tx.ExecOne("FOO")

					return err
				},
			},
			[]string{"BEGIN", "FOO", "ROLLBACK"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a new database
			dbPool, dbResource, dbPort := testutils.SetupDb(t, iDbUser, iDbPass, iDbName)
			defer testutils.TeardownDb(t, dbPool, dbResource)

			logger, buf := testutils.SetupLogger(t)

			repo := NewRepository("TestIntegration_pgRepository_InTx", iDbHost, dbPort, iDbName, iDbUser, iDbPass, logger)

			if err := repo.InTx(tt.args.operation); (err != nil) != tt.wantErr {
				t.Errorf("InTx() error = %v, wantErr %v", err, tt.wantErr)
			}

			logs := buf.Lines()
			if assert.Equal(t, len(tt.lines), len(logs)) {
				for i, l := range tt.lines {
					assert.Contains(t, logs[i], l)
				}
			}
		})
	}
}

func TestIntegration_Repository_Get(t *testing.T) {
	type fields struct {
		pgRepository pgRepository
		pgOptions    *pg.Options
		logger       *zap.Logger
	}
	type args struct {
		code code.Code
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Roadmap
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a new database
			dbPool, dbResource, dbPort := testutils.SetupDb(t, iDbUser, iDbPass, iDbName)
			defer testutils.TeardownDb(t, dbPool, dbResource)

			logger := zaptest.NewLogger(t)

			repo := NewRepository("TestIntegration_pgRepository_InTx", iDbHost, dbPort, iDbName, iDbUser, iDbPass, logger)

			got, err := repo.Get(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegration_Repository_Upsert(t *testing.T) {
	type fields struct {
		pgRepository pgRepository
		pgOptions    *pg.Options
		logger       *zap.Logger
	}
	type args struct {
		roadmap Roadmap
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a new database
			dbPool, dbResource, dbPort := testutils.SetupDb(t, iDbUser, iDbPass, iDbName)
			defer testutils.TeardownDb(t, dbPool, dbResource)

			logger := zaptest.NewLogger(t)

			repo := NewRepository("TestIntegration_pgRepository_InTx", iDbHost, dbPort, iDbName, iDbUser, iDbPass, logger)

			if err := repo.Upsert(tt.args.roadmap); (err != nil) != tt.wantErr {
				t.Errorf("Upsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
