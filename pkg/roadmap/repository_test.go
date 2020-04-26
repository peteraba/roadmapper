// +build integration

package roadmap

import (
	"reflect"
	"testing"

	"github.com/go-pg/pg"

	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/testutils"
)

const (
	iDbName = "rdmp"
	iDbUser = "rdmp"
	iDbPass = "secret"
)

func TestIntegration_Repository_InTx(t *testing.T) {
	var res interface{}

	type args struct {
		operation func(tx *pg.Tx) error
	}
	tests := []struct {
		name    string
		args    args
		queries []string
		want    interface{}
		wantErr bool
	}{
		{
			"empty",
			args{
				operation: func(tx *pg.Tx) error {
					res = nil

					return nil
				},
			},
			[]string{"BEGIN", "COMMIT"},
			nil,
			false,
		},
		{
			"simple select",
			args{
				operation: func(tx *pg.Tx) error {
					_, err := tx.Model(res).ExecOne("SELECT 2")

					return err
				},
			},
			[]string{"BEGIN", "SELECT 2", "COMMIT"},
			nil,
			false,
		},
		{
			"error",
			args{
				operation: func(tx *pg.Tx) error {
					_, err := tx.Model(res).ExecOne("FOO")

					return err
				},
			},
			[]string{"BEGIN", "FOO", "ROLLBACK"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, buf := testutils.SetupLogger()
			baseRepo, teardown := testutils.SetupRepository(t, "TestIntegration_pgRepository_InTx", iDbUser, iDbPass, iDbName, logger)
			defer teardown()

			repo := Repository{baseRepo}

			if err := repo.InTx(tt.args.operation); (err != nil) != tt.wantErr {
				t.Errorf("InTx() error = %v, wantErr %v", err, tt.wantErr)
			}

			testutils.AssertQueries(t, buf, tt.queries)

			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("InTx() got = %v, want %v", res, tt.want)
			}
		})
	}
}

func TestIntegration_Repository_Get(t *testing.T) {
	r := Roadmap{ID: 123, Title: "Foo", DateFormat: "2006-01-02"}

	type args struct {
		code code.Code
	}
	tests := []struct {
		name    string
		fixture []interface{}
		args    args
		queries []string
		want    *Roadmap
		wantErr bool
	}{
		{
			"not found",
			[]interface{}{},
			args{code.Code64(123)},
			[]string{
				`SELECT .* FROM \\\"roadmaps\\\" AS \\\"roadmap\\\" WHERE .*\\\"id\\\" \= 123`,
			},
			nil,
			true,
		},
		{
			"success",
			[]interface{}{
				&r,
			},
			args{code.Code64(123)},
			[]string{
				`SELECT .* FROM \\\"roadmaps\\\" AS \\\"roadmap\\\" WHERE .*\\\"id\\\" \= 123`,
				`UPDATE roadmaps SET accessed_at = NOW\(\) WHERE id = 123`,
			},
			&r,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, buf := testutils.SetupLogger()
			baseRepo, teardown := testutils.SetupRepository(t, "TestIntegration_Repository_Get", iDbUser, iDbPass, iDbName, logger, tt.fixture...)
			defer teardown()

			repo := Repository{baseRepo}

			got, err := repo.Get(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				got.AccessedAt = tt.want.AccessedAt
				got.CreatedAt = tt.want.CreatedAt
				got.UpdatedAt = tt.want.UpdatedAt
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}

			testutils.AssertQueriesRegexp(t, buf, tt.queries)
		})
	}
}

func TestIntegration_Repository_Create(t *testing.T) {
	r0 := Roadmap{ID: 123, Title: "Foo", DateFormat: "2006-01-02"}
	r1 := Roadmap{ID: 123, Title: "Foo", DateFormat: "2006-01-02"}
	r2 := Roadmap{ID: 123, Title: "Foo", DateFormat: "2006-01-02"}

	type args struct {
		roadmap Roadmap
	}
	tests := []struct {
		name    string
		fixture []interface{}
		args    args
		queries []string
		wantErr bool
	}{
		{
			"conflict",
			[]interface{}{
				&r0,
			},
			args{r1},
			[]string{
				`123, DEFAULT, 'Foo', '2006-01-02', DEFAULT, DEFAULT, DEFAULT`,
			},
			true,
		},
		{
			"success",
			[]interface{}{},
			args{r2},
			[]string{
				`123, DEFAULT, 'Foo', '2006-01-02', DEFAULT, DEFAULT, DEFAULT`,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, buf := testutils.SetupLogger()
			baseRepo, teardown := testutils.SetupRepository(t, "TestIntegration_Repository_Get", iDbUser, iDbPass, iDbName, logger, tt.fixture...)
			defer teardown()

			repo := Repository{baseRepo}

			if err := repo.Create(tt.args.roadmap); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			testutils.AssertQueriesRegexp(t, buf, tt.queries)
		})
	}
}
