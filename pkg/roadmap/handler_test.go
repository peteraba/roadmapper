package roadmap

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo"
	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_handler_getRoadmapHTML(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	h := &Handler{}

	// Assertions
	require.NoError(t, h.GetRoadmapHTML(ctx))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "</html>")
}

// func Test_handler_createRoadmapHTML(t *testing.T) {
// 	logger := zap.NewNop()
//
// 	// Setup
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	h := &Handler{Logger: logger}
//
// 	// Assertions
// 	require.NoError(t, h.CreateRoadmapHTML(c))
//
// 	assert.Equal(t, http.StatusCreated, rec.Code)
// 	assert.Equal(t, userJSON, rec.Body.String())
// }

func Test_handler_getPrevID(t *testing.T) {
	var ui64 uint64 = 63*64*64 + 63*64 + 36

	type fields struct {
		rw           DbReadWriter
		cb           code.Builder
		appVersion   string
		matomoDomain string
		docBaseURL   string
		selfHosted   bool
		logger       *zap.Logger
	}
	type args struct {
		identifier string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *uint64
		wantErr bool
	}{
		{
			"foo",
			fields{},
			args{"~~A"},
			&ui64,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(
				tt.fields.logger,
				tt.fields.rw,
				tt.fields.cb,
				tt.fields.appVersion,
				tt.fields.matomoDomain,
				tt.fields.docBaseURL,
				tt.fields.selfHosted,
			)

			got, err := h.getPrevID(tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPrevID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPrevID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_isValidRoadmapRequest(t *testing.T) {
	type fields struct {
		rw           PqReadWriter
		cb           code.Builder
		appVersion   string
		matomoDomain string
		docBaseURL   string
		selfHosted   bool
		logger       *zap.Logger
	}
	type args struct {
		ctx echo.Context
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
			h := NewHandler(
				tt.fields.logger,
				tt.fields.rw,
				tt.fields.cb,
				tt.fields.appVersion,
				tt.fields.matomoDomain,
				tt.fields.docBaseURL,
				tt.fields.selfHosted,
			)

			if err := h.isValidRoadmapRequest(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("isValidRoadmapRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func Test_handler_getRoadmapImage(t *testing.T) {
// 	// Setup
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "/", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
// 	rec := httptest.NewRecorder()
// 	ctx := e.NewContext(req, rec)
// 	h := &Handler{}
//
// 	// Assertions
// 	require.NoError(t, h.GetRoadmapImage(ctx))
//
// 	assert.Equal(t, http.StatusCreated, rec.Code)
// 	// assert.Equal(t, userJSON, rec.Body.String())
// }
