// +build e2e

package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/peteraba/roadmapper/pkg/roadmap"
	"github.com/peteraba/roadmapper/pkg/testutils"
)

func TestE2E_CLI(t *testing.T) {
	var (
		dateFormat        = "2006-01-02"
		fw, lh     uint64 = 800, 30
		rw                = roadmap.NewIO()
	)

	type args struct {
		rw                  roadmap.IO
		content, output     string
		format              string
		dateFormat, baseUrl string
		fw, lh              uint64
		mt                  bool
	}

	tests := []struct {
		name string
		args args
	}{
		{
			"svg size",
			args{
				rw,
				e2eTxt,
				"test.svg",
				"svg",
				dateFormat,
				e2eBaseURL,
				fw,
				lh,
				false,
			},
		},
		{
			"png size",
			args{
				rw,
				e2eTxt,
				"test.png",
				"png",
				dateFormat,
				e2eBaseURL,
				fw,
				lh,
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zap.NewNop()

			expectedData := testutils.LoadFile(t, "golden_files", tt.args.output)

			err := Render(
				rw,
				logger,
				tt.args.content,
				tt.args.output,
				tt.args.format,
				tt.args.dateFormat,
				tt.args.baseUrl,
				tt.args.fw,
				tt.args.lh,
				tt.args.mt,
			)

			require.NoError(t, err)

			actualData, err := ioutil.ReadFile(tt.args.output)
			require.NoError(t, err)

			if testutils.ShouldUpdateGoldenFiles() {
				testutils.SaveFile(t, actualData, "golden_files", tt.args.output)
			}

			assert.Equal(t, expectedData, actualData)

			// ed0, ad0 := float64(len(expectedData)), float64(len(actualData))
			// ed1, ad1 := ed0*1.2, ad0*1.2
			//
			// assert.Greater(t, ed1, ad0, "generated and golden files differ a lot")
			// assert.Less(t, ed0, ad1, "generated and golden files differ a lot")

			if !t.Failed() {
				err = os.Remove(tt.args.output) // remove a single file
				if err != nil {
					t.Errorf("failed to delete file: %s", tt.args.output)
				}
			}
		})
	}
}
