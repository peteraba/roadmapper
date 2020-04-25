//go:generate mockery -name DbReadWriter -case snake -inpkg -output .

package roadmap

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/herr"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/rasterizer"
	"go.uber.org/zap"
)

type (
	Handler struct {
		Logger       *zap.Logger
		rw           DbReadWriter
		cb           code.Builder
		appVersion   string
		matomoDomain string
		docBaseURL   string
		selfHosted   bool
	}
)

func NewHandler(logger *zap.Logger, rw DbReadWriter, cb code.Builder, appVersion, matomoDomain, docBaseURL string, selfHosted bool) *Handler {
	return &Handler{
		Logger:       logger,
		rw:           rw,
		cb:           cb,
		appVersion:   appVersion,
		matomoDomain: matomoDomain,
		docBaseURL:   docBaseURL,
		selfHosted:   selfHosted,
	}
}

func (h *Handler) GetRoadmapHTML(ctx echo.Context) error {
	identifier := ctx.Param("identifier")

	r, err := load(h.rw, h.cb, identifier)
	if err != nil {
		return h.displayHTML(ctx, r, err)
	}

	return h.displayHTML(ctx, r, nil)
}

func (h *Handler) displayHTML(ctx echo.Context, r *Roadmap, origErr error) error {
	output, err := r.viewHtml(h.appVersion, h.matomoDomain, h.docBaseURL, ctx.Request().RequestURI, h.selfHosted, origErr)
	if origErr == nil && err == nil {
		return ctx.HTML(http.StatusOK, output)
	}

	if err != nil {
		h.Logger.Info("failed to create HTML response", zap.Error(err))

		return ctx.String(herr.ToHttpCode(err, http.StatusInternalServerError), err.Error())
	}

	return ctx.HTML(herr.ToHttpCode(origErr, http.StatusInternalServerError), output)
}

func (h *Handler) CreateRoadmapHTML(ctx echo.Context) error {
	prevID, err := h.getPrevID(ctx.Param("identifier"))
	if err != nil {
		h.Logger.Info("failed to parse the identifier parameter", zap.Error(err))

		return h.displayHTML(ctx, nil, err)
	}

	areYouAHuman := ctx.FormValue("areYouAHuman")
	timeSpent := ctx.FormValue("ts")
	err = h.isValidRoadmapRequest(areYouAHuman, timeSpent)
	if err != nil {
		h.Logger.Info("not timeSpent valid roadmap request", zap.Error(err))

		return h.displayHTML(ctx, nil, herr.NewFromError(err, http.StatusBadRequest))
	}

	title := ctx.FormValue("title")
	content := ctx.FormValue("txt")
	dateFormat := ctx.FormValue("dateFormat")
	baseURL := ctx.FormValue("baseUrl")
	now := time.Now()

	roadmap := Content(content).ToRoadmap(code.NewCode64().ID(), prevID, title, dateFormat, baseURL, now)

	err = h.isValidRoadmap(roadmap)
	if err != nil {
		h.Logger.Info("not timeSpent valid roadmap", zap.Error(err))

		return h.displayHTML(ctx, nil, herr.NewFromError(err, http.StatusBadRequest))
	}

	err = h.rw.Write(roadmap)
	if err != nil {
		h.Logger.Info("failed to write the new roadmap", zap.Error(err))

		return h.displayHTML(ctx, &roadmap, err)
	}

	c, _ := h.cb.NewFromID(roadmap.ID)

	newURL := fmt.Sprintf("/%s", c.String())

	return ctx.Redirect(http.StatusSeeOther, newURL)
}

func (h *Handler) getPrevID(identifier string) (*uint64, error) {
	if identifier == "" {
		return nil, nil
	}

	c, err := h.cb.NewFromString(identifier)
	if err != nil {
		return nil, herr.NewFromError(err, http.StatusNotFound)
	}

	n := c.ID()

	return &n, err
}

const iAmHuman = "Yes, I am indeed."
const onlyHumansAreAllowed = "only humans are allowed"

func (h *Handler) isValidRoadmapRequest(areYouAHuman, timeSpent string) error {
	ts, _ := strconv.ParseUint(timeSpent, 10, 64)

	if areYouAHuman == iAmHuman && ts == 0 {
		return nil
	}

	if areYouAHuman == "" && ts > 0 {
		return nil
	}

	return fmt.Errorf(onlyHumansAreAllowed)
}

func (h *Handler) isValidRoadmap(r Roadmap) error {
	if len(r.Title) == 0 || len(r.DateFormat) == 0 || len(r.Projects) == 0 {
		return fmt.Errorf("title, dateFormat and txt are mandatory fields")
	}

	for _, p := range r.Projects {
		if p.Dates != nil && p.Dates.EndAt.Before(p.Dates.StartAt) {
			return fmt.Errorf(
				"end at before start at. start at: %s, end at: %s",
				p.Dates.StartAt.Format(r.DateFormat),
				p.Dates.EndAt.Format(r.DateFormat),
			)
		}
	}

	return nil
}

func (h *Handler) GetRoadmapImage(ctx echo.Context) error {
	format, err := NewFormatType(ctx.Param("format"))
	if err != nil {
		h.Logger.Info("format is not supported", zap.Error(err))

		return ctx.String(herr.ToHttpCode(err, http.StatusBadRequest), "format is not supported")
	}

	fw, _ := strconv.ParseUint(ctx.QueryParam("width"), 10, 64)

	lh, _ := strconv.ParseUint(ctx.QueryParam("lineHeight"), 10, 64)

	fw, lh = GetCanvasSizes(fw, lh)

	roadmap, err := load(h.rw, h.cb, ctx.Param("identifier"))
	if err != nil {
		h.Logger.Info("roadmap not found", zap.Error(err))

		return ctx.String(herr.ToHttpCode(err, http.StatusNotFound), "roadmap not found")
	}

	cvs := roadmap.ToVisual().Draw(float64(fw), float64(lh))

	img := RenderImg(cvs, format)

	setHeaderContentType(ctx.Response().Header(), format)

	err = ctx.String(http.StatusOK, string(img))

	return err
}

func load(rw DbReadWriter, b code.Builder, identifier string) (*Roadmap, error) {
	if identifier == "" {
		return nil, nil
	}

	c, err := b.NewFromString(identifier)
	if err != nil {
		return nil, herr.NewFromError(err, http.StatusBadRequest)
	}

	roadmap, err := rw.Get(c)
	if err != nil {
		return nil, herr.NewFromError(err, http.StatusInternalServerError)
	}

	return roadmap, nil
}

const (
	defaultSvgWidth      = 800
	defaultSvgLineHeight = 40
)

type FileFormat string

const (
	SvgFormat FileFormat = "svg"
	PngFormat FileFormat = "png"
)

func NewFormatType(t string) (FileFormat, error) {
	switch t {
	case "svg":
		return SvgFormat, nil
	case "png":
		return PngFormat, nil
	}

	return "", fmt.Errorf("unsupported image format: %s", t)
}

func setHeaderContentType(header http.Header, fileFormat FileFormat) {
	switch fileFormat {
	case SvgFormat:
		header.Set(echo.HeaderContentType, "image/svg+xml")
	case PngFormat:
		header.Set(echo.HeaderContentType, "image/png")
	}
}

func RenderImg(cvs *canvas.Canvas, fileFormat FileFormat) []byte {
	var buf bytes.Buffer

	switch fileFormat {
	case SvgFormat:
		img := canvas.NewSVG(&buf, cvs.W, cvs.H)

		cvs.Render(img)

		_ = img.Close()
	case PngFormat:
		w := rasterizer.PNGWriter(3.2)

		err := w(&buf, cvs)
		if err != nil {
			return nil
		}
	}

	return buf.Bytes()
}

func GetCanvasSizes(fw, lh uint64) (uint64, uint64) {
	if fw < defaultSvgWidth {
		fw = defaultSvgWidth
	}
	if lh < defaultSvgLineHeight {
		lh = defaultSvgLineHeight
	}

	return fw, lh
}
