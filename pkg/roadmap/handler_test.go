package roadmap

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/herr"
	"github.com/peteraba/roadmapper/pkg/testutils"
)

func Test_handler_getRoadmapJSON(t *testing.T) {
	t.Run("fail - not found", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/abc", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/api/:identifier")
		ctx.SetParamNames("identifier")
		ctx.SetParamValues("abc")

		h, drwMock := setupHandler()
		drwMock.
			On("Get", mock.AnythingOfType("code.Code64")).
			Return(nil, herr.NewFromError(assert.AnError, http.StatusNotFound))

		// Run
		err := h.GetRoadmapJSON(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "Roadmap Payload Loading Error")
	})

	t.Run("fail - loading", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/abc", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/api/:identifier")
		ctx.SetParamNames("identifier")
		ctx.SetParamValues("abc")

		h, drwMock := setupHandler()
		drwMock.
			On("Get", mock.AnythingOfType("code.Code64")).
			Return(nil, assert.AnError)

		// Run
		err := h.GetRoadmapJSON(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Roadmap Payload Loading Error")
	})

	t.Run("fail - no identifier = not implemented", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h, _ := setupHandler()

		// Run
		err := h.GetRoadmapJSON(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotImplemented, rec.Code)
		assert.Contains(t, rec.Body.String(), "Not Implemented")
	})

	t.Run("fail - no roadmap + no error = internal server error", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/api/:identifier")
		ctx.SetParamNames("identifier")
		ctx.SetParamValues("abc")

		h, drwMock := setupHandler()
		drwMock.
			On("Get", mock.AnythingOfType("code.Code64")).
			Return(nil, nil)

		// Run
		err := h.GetRoadmapJSON(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Roadmap Payload Loading Error")
	})

	t.Run("success - existing", func(t *testing.T) {
		// Setup
		rdmp := createStubRoadmap()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/abc", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/api/:identifier")
		ctx.SetParamNames("identifier")
		ctx.SetParamValues("abc")

		h, drwMock := setupHandler()
		drwMock.
			On("Get", mock.AnythingOfType("code.Code64")).
			Return(rdmp, nil)

		// Run
		err := h.GetRoadmapJSON(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), rdmp.Title)
		drwMock.AssertExpectations(t)
	})
}

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
			Return(nil, herr.NewFromError(assert.AnError, http.StatusNotFound))

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

func Test_handler_createRoadmapHTML(t *testing.T) {
	t.Run("error - getPrevID failure", func(t *testing.T) {
		// Setup
		rdmp := createStubRoadmap()

		f := make(url.Values)
		f.Set("title", rdmp.Title)
		f.Set("txt", rdmp.String())
		f.Set("dateFormat", rdmp.DateFormat)
		f.Set("baseUrl", rdmp.BaseURL)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/üüü", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/:identifier")
		ctx.SetParamNames("identifier")
		ctx.SetParamValues("üüü")

		h, _ := setupHandler()

		// Run
		err := h.CreateRoadmapHTML(ctx)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})

	t.Run("error - only humans are allowed", func(t *testing.T) {
		// Setup
		rdmp := createStubRoadmap()

		f := make(url.Values)
		f.Set("title", rdmp.Title)
		f.Set("txt", rdmp.String())
		f.Set("dateFormat", rdmp.DateFormat)
		f.Set("baseUrl", rdmp.BaseURL)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/")

		h, _ := setupHandler()

		// Run
		err := h.CreateRoadmapHTML(ctx)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})

	t.Run("error - missing title", func(t *testing.T) {
		// Setup
		rdmp := createStubRoadmap()

		f := make(url.Values)
		f.Set("title", "")
		f.Set("txt", rdmp.String())
		f.Set("dateFormat", rdmp.DateFormat)
		f.Set("baseUrl", rdmp.BaseURL)
		f.Set("ts", "20")

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/")

		h, _ := setupHandler()

		// Run
		err := h.CreateRoadmapHTML(ctx)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})

	t.Run("error - writing database", func(t *testing.T) {
		// Setup
		rdmp := createStubRoadmap()

		f := make(url.Values)
		f.Set("title", rdmp.Title)
		f.Set("txt", rdmp.String())
		f.Set("dateFormat", rdmp.DateFormat)
		f.Set("baseUrl", rdmp.BaseURL)
		f.Set("ts", "20")

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/")

		h, drwMock := setupHandler()
		drwMock.
			On("Create", mock.AnythingOfType("Roadmap")).
			Return(assert.AnError)

		// Run
		err := h.CreateRoadmapHTML(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
		drwMock.AssertExpectations(t)
	})

	t.Run("success - create from scratch", func(t *testing.T) {
		// Setup
		rdmp := createStubRoadmap()

		f := make(url.Values)
		f.Set("title", rdmp.Title)
		f.Set("txt", rdmp.String())
		f.Set("dateFormat", rdmp.DateFormat)
		f.Set("baseUrl", rdmp.BaseURL)
		f.Set("ts", "20")

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/")

		h, drwMock := setupHandler()
		drwMock.
			On("Create", mock.AnythingOfType("Roadmap")).
			Return(nil)

		// Run
		err := h.CreateRoadmapHTML(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Empty(t, rec.Body.String())
		drwMock.AssertExpectations(t)
	})
}

func Test_handler_getPrevID(t *testing.T) {
	var ui64 uint64 = 63*64*64 + 63*64 + 36

	type args struct {
		identifier string
	}
	tests := []struct {
		name    string
		args    args
		want    *uint64
		wantErr bool
	}{
		{
			"foo",
			args{"~~A"},
			&ui64,
			false,
		},
		{
			"invalid identifier",
			args{"ü"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}

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

func TestHandler_isValidRoadmapRequest(t *testing.T) {
	type args struct {
		areYouAHuman string
		timeSpent    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"error - no areYouAHuman, no timeSpent",
			args{},
			true,
		},
		{
			"error - invalid areYouAHuman, no timeSpent",
			args{"foo", ""},
			true,
		},
		{
			"success - valid areYouAHuman, no timeSpent",
			args{iAmHuman, ""},
			false,
		},
		{
			"success - no areYouAHuman, valid timeSpent",
			args{"", "10"},
			false,
		},
		{
			"error - invalid areYouAHuman, valid timeSpent",
			args{"foo", "10"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			if err := h.isValidRoadmapRequest(tt.args.areYouAHuman, tt.args.timeSpent); (err != nil) != tt.wantErr {
				t.Errorf("isValidRoadmapRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_isValidRoadmap(t *testing.T) {
	startAt0 := time.Date(2020, 1, 20, 0, 0, 0, 0, time.UTC)
	endAt0 := time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC)

	type args struct {
		r Roadmap
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"error - no title",
			args{
				r: Roadmap{
					Title:      "",
					DateFormat: "2006-01-02",
					Projects: []Project{
						{},
					},
				},
			},
			true,
		},
		{
			"error - no dateFormat",
			args{
				r: Roadmap{
					Title:      "foo",
					DateFormat: "",
					Projects: []Project{
						{},
					},
				},
			},
			true,
		},
		{
			"error - no txt",
			args{
				r: Roadmap{
					Title:      "foo",
					DateFormat: "2006-01-02",
					Projects:   []Project{},
				},
			},
			true,
		},
		{
			"error - end date before start date",
			args{
				r: Roadmap{
					Title:      "foo",
					DateFormat: "2006-01-02",
					Projects: []Project{
						{Dates: &Dates{StartAt: endAt0, EndAt: startAt0}},
					},
				},
			},
			true,
		},
		{
			"success - end date before start date",
			args{
				r: Roadmap{
					Title:      "foo",
					DateFormat: "2006-01-02",
					Projects: []Project{
						{Dates: &Dates{StartAt: startAt0, EndAt: endAt0}},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			if err := h.isValidRoadmap(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("isValidRoadmap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_handler_getRoadmapImage(t *testing.T) {
	t.Run("error - format not supported", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/abc/ico", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/:identifier/:format")
		ctx.SetParamNames("identifier", "format")
		ctx.SetParamValues("abc", "ico")

		h, _ := setupHandler()

		// Run
		err := h.GetRoadmapImage(ctx)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})

	t.Run("error - bad request", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/abc/ü", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/:identifier/:format")
		ctx.SetParamNames("identifier", "format")
		ctx.SetParamValues("ü", "svg")

		h, _ := setupHandler()

		// Run
		err := h.GetRoadmapImage(ctx)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})

	t.Run("error - roadmap not found", func(t *testing.T) {
		// Setup
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
			Return(nil, herr.NewFromError(assert.AnError, http.StatusNotFound))

		// Run
		err := h.GetRoadmapImage(ctx)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})

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
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "</svg>")
	})

	t.Run("success - non-empty roadmap SVG", func(t *testing.T) {
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
		require.NoError(t, err)

		// Update golden files
		if testutils.ShouldUpdateGoldenFiles() {
			testutils.SaveFile(t, rec.Body.Bytes(), "golden_files", "nonempty.svg")
		}

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, rec.Body.Bytes(), testutils.LoadFile(t, "golden_files", "nonempty.svg"))
	})

	t.Run("success - non-empty roadmap PNG", func(t *testing.T) {
		// Setup
		rdmp := createStubRoadmap()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/abc/png", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/:identifier/:format")
		ctx.SetParamNames("identifier", "format")
		ctx.SetParamValues("abc", "png")

		h, drwMock := setupHandler()
		drwMock.
			On("Get", mock.AnythingOfType("code.Code64")).
			Return(rdmp, nil)

		// Run
		err := h.GetRoadmapImage(ctx)

		// Update golden files
		if testutils.ShouldUpdateGoldenFiles() {
			testutils.SaveFile(t, rec.Body.Bytes(), "golden_files", "nonempty.png")
		}

		// Assertions
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, rec.Body.Bytes(), testutils.LoadFile(t, "golden_files", "nonempty.png"))
	})
}

func createStubRoadmap() *Roadmap {
	createdAt := time.Date(2020, 1, 19, 0, 0, 0, 0, time.UTC)
	startAt0 := time.Date(2020, 1, 20, 0, 0, 0, 0, time.UTC)
	startAt1 := time.Date(2020, 1, 22, 0, 0, 0, 0, time.UTC)
	endAt0 := time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC)
	endAt1 := time.Date(2020, 2, 5, 0, 0, 0, 0, time.UTC)

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
