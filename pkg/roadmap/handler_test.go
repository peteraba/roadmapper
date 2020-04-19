package roadmap

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/peteraba/roadmapper/pkg/herr"

	"github.com/stretchr/testify/mock"

	"github.com/labstack/echo"
	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_handler_getRoadmapHTML(t *testing.T) {
	t.Run("fail - not found", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/abc", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/:identifier")
		ctx.SetParamNames("identifier")
		ctx.SetParamValues("abc")

		h, drwMock := setupHandler()
		drwMock.
			On("Get", mock.AnythingOfType("code.Code64")).
			Return(nil, herr.NewHttpError(assert.AnError, http.StatusNotFound))

		// Run
		err := h.GetRoadmapHTML(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "</html>")
	})

	t.Run("fail - loading", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/abc", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/:identifier")
		ctx.SetParamNames("identifier")
		ctx.SetParamValues("abc")

		h, drwMock := setupHandler()
		drwMock.
			On("Get", mock.AnythingOfType("code.Code64")).
			Return(nil, assert.AnError)

		// Run
		err := h.GetRoadmapHTML(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "</html>")
	})

	t.Run("success - new", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		h, _ := setupHandler()

		// Run
		err := h.GetRoadmapHTML(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "</html>")
	})

	t.Run("success - existing", func(t *testing.T) {
		// Setup
		rdmp := createStubRoadmap()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/abc", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/:identifier")
		ctx.SetParamNames("identifier")
		ctx.SetParamValues("abc")

		h, drwMock := setupHandler()
		drwMock.
			On("Get", mock.AnythingOfType("code.Code64")).
			Return(rdmp, nil)

		// Run
		err := h.GetRoadmapHTML(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "</html>")
		drwMock.AssertExpectations(t)
	})
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
		rw           Repository
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

func Test_handler_getRoadmapImage(t *testing.T) {
	t.Run("success - empty roadmap", func(t *testing.T) {
		// Setup
		rdmp := &Roadmap{}
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/abc/svg", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/:identifier/:format")
		ctx.SetParamNames("identifier", "format")
		ctx.SetParamValues("abc", "svg")

		h, drwMock := setupHandler()
		drwMock.
			On("Get", mock.AnythingOfType("code.Code64")).
			Return(rdmp, nil)

		// Run
		err := h.GetRoadmapImage(ctx)

		// Assertions
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "</svg>")
	})

	t.Run("success - non-empty roadmap", func(t *testing.T) {
		// Setup
		rdmp := createStubRoadmap()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/abc/svg", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/:identifier/:format")
		ctx.SetParamNames("identifier", "format")
		ctx.SetParamValues("abc", "svg")

		h, drwMock := setupHandler()
		drwMock.
			On("Get", mock.AnythingOfType("code.Code64")).
			Return(rdmp, nil)

		// Run
		err := h.GetRoadmapImage(ctx)

		// Assertions
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "</svg>")
	})
}

func createStubRoadmap() *Roadmap {
	createdAt := time.Date(2020, 4, 19, 0, 0, 0, 0, time.UTC)
	startAt0 := time.Date(2020, 4, 20, 0, 0, 0, 0, time.UTC)
	startAt1 := time.Date(2020, 4, 22, 0, 0, 0, 0, time.UTC)
	endAt0 := time.Date(2020, 4, 30, 0, 0, 0, 0, time.UTC)
	endAt1 := time.Date(2020, 5, 5, 0, 0, 0, 0, time.UTC)

	return &Roadmap{
		ID:         123,
		Title:      "abc",
		DateFormat: "2006-01-02",
		BaseURL:    "https://example.com/",
		Projects: []Project{
			{Indentation: 0, Title: "foo", Dates: &Dates{StartAt: startAt0, EndAt: endAt0}, Percentage: 50},
			{Indentation: 1, Title: "bar", Dates: &Dates{StartAt: startAt1, EndAt: endAt1}, Percentage: 40},
			{Indentation: 0, Title: "baz", Dates: &Dates{StartAt: endAt0, EndAt: endAt1}, Percentage: 20},
		},
		Milestones: []Milestone{
			{Title: "quix", DeadlineAt: &endAt1},
		},
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
		AccessedAt: time.Now(),
	}
}

func setupHandler() (*Handler, *MockDbReadWriter) {
	rw := &MockDbReadWriter{}
	h := NewHandler(zap.NewNop(), rw, code.Builder{}, "", "", "", false)

	return h, rw
}
