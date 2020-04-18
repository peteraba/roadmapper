package roadmap

//
// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"reflect"
// 	"testing"
//
// 	"go.uber.org/zap"
//
// 	"github.com/labstack/echo"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )
//
// func Test_handler_getRoadmapHTML(t *testing.T) {
//
// 	// Setup
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	h := &handler{}
//
// 	// Assertions
// 	require.NoError(t, h.getRoadmapHTML(c))
//
// 	assert.Equal(t, http.StatusCreated, rec.Code)
// 	assert.Equal(t, userJSON, rec.Body.String())
// }
//
// func Test_handler_createRoadmapHTML(t *testing.T) {
//
// 	// Setup
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	h := &handler{mockDB}
//
// 	// Assertions
// 	require.NoError(t, h.createRoadmapHTML(c))
//
// 	assert.Equal(t, http.StatusCreated, rec.Code)
// 	assert.Equal(t, userJSON, rec.Body.String())
// }
//
// func Test_handler_getPrevID(t *testing.T) {
// 	type fields struct {
// 		rw           PqReadWriter
// 		cb           CodeBuilder
// 		matomoDomain string
// 		docBaseURL   string
// 		selfHosted   bool
// 		logger       *zap.Logger
// 	}
// 	type args struct {
// 		identifier string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *uint64
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &handler{
// 				rw:           tt.fields.rw,
// 				cb:           tt.fields.cb,
// 				matomoDomain: tt.fields.matomoDomain,
// 				docBaseURL:   tt.fields.docBaseURL,
// 				selfHosted:   tt.fields.selfHosted,
// 				logger:       tt.fields.logger,
// 			}
// 			got, err := h.getPrevID(tt.args.identifier)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("getPrevID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("getPrevID() got = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func Test_handler_isValidRoadmapRequest(t *testing.T) {
// 	type fields struct {
// 		rw           PqReadWriter
// 		cb           CodeBuilder
// 		matomoDomain string
// 		docBaseURL   string
// 		selfHosted   bool
// 		logger       *zap.Logger
// 	}
// 	type args struct {
// 		ctx echo.Context
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &handler{
// 				rw:           tt.fields.rw,
// 				cb:           tt.fields.cb,
// 				matomoDomain: tt.fields.matomoDomain,
// 				docBaseURL:   tt.fields.docBaseURL,
// 				selfHosted:   tt.fields.selfHosted,
// 				logger:       tt.fields.logger,
// 			}
// 			if err := h.isValidRoadmapRequest(tt.args.ctx); (err != nil) != tt.wantErr {
// 				t.Errorf("isValidRoadmapRequest() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
//
// func Test_handler_getRoadmapImage(t *testing.T) {
// 	// Setup
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	h := &handler{mockDB}
//
// 	// Assertions
// 	require.NoError(t, h.createUser(c))
//
// 	assert.Equal(t, http.StatusCreated, rec.Code)
// 	assert.Equal(t, userJSON, rec.Body.String())
// }
