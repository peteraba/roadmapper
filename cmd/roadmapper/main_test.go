package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/peteraba/roadmapper/pkg/code"
)

func Test_createApp(t *testing.T) {
	var (
		logger = zap.NewNop()
		b      = code.Builder{}
	)

	type args struct {
		logger *zap.Logger
		b      code.Builder
	}
	tests := []struct {
		name string
		args args
		want *cli.App
	}{
		{
			"default",
			args{
				logger: logger,
				b:      b,
			},
			&cli.App{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createApp(tt.args.logger, tt.args.b)

			assert.Len(t, got.Commands, 5)
		})
	}
}
