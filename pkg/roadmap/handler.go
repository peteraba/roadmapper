//go:generate mockery -name DbReadWriter -case snake -inpkg -output .

package roadmap

import (
	"bytes"
	"fmt"
	"image/gif"
	"image/jpeg"
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

	err = h.displayHTML(ctx, r, nil)

	_ = h.Logger.Sync()

	return err
}

func (h *Handler) displayHTML(ctx echo.Context, r *Roadmap, origErr error) error {
	output, err := r.viewHtml(h.appVersion, h.matomoDomain, h.docBaseURL, ctx.Request().RequestURI, h.selfHosted, origErr)
	if err != nil {
		h.Logger.Info("failed to create HTML response", zap.Error(err))

		return ctx.HTML(herr.ToHttpCode(err, http.StatusInternalServerError), err.Error())
	}

	return ctx.HTML(herr.ToHttpCode(origErr, http.StatusOK), output)
}

func (h *Handler) CreateRoadmapHTML(ctx echo.Context) error {
	prevID, err := h.getPrevID(ctx.Param("identifier"))
	if err != nil {
		h.Logger.Info("failed to parse the identifier parameter", zap.Error(err))

		return h.displayHTML(ctx, nil, err)
	}

	err = h.isValidRoadmapRequest(ctx)
	if err != nil {
		h.Logger.Info("not a valid roadmap request", zap.Error(err))

		return h.displayHTML(ctx, nil, err)
	}

	title := ctx.FormValue("title")
	content := ctx.FormValue("txt")
	dateFormat := ctx.FormValue("dateFormat")
	baseURL := ctx.FormValue("baseUrl")
	now := time.Now()

	roadmap := Content(content).ToRoadmap(code.NewCode64().ID(), prevID, title, dateFormat, baseURL, now)

	err = h.isValidRoadmap(roadmap, dateFormat)
	if err != nil {
		h.Logger.Info("not a valid roadmap", zap.Error(err))

		return h.displayHTML(ctx, nil, err)
	}

	err = h.rw.Write(roadmap)
	if err != nil {
		h.Logger.Info("failed to write the new roadmap", zap.Error(err))

		return h.displayHTML(ctx, &roadmap, err)
	}

	c, err := h.cb.NewFromID(roadmap.ID)
	if err != nil {
		h.Logger.Info("failed to generate the new  url", zap.Error(err))

		return h.displayHTML(ctx, &roadmap, err)
	}

	newURL := fmt.Sprintf("/%s", c.String())

	err = ctx.Redirect(http.StatusSeeOther, newURL)

	_ = h.Logger.Sync()

	return err
}

func (h *Handler) getPrevID(identifier string) (*uint64, error) {
	if identifier == "" {
		return nil, nil
	}

	c, err := h.cb.NewFromString(identifier)
	if err != nil {
		return nil, herr.NewHttpError(err, http.StatusNotFound)
	}

	n := c.ID()

	return &n, err
}

const iAmHuman = "Yes, I am indeed."
const onlyHumansAreAllowed = "only humans are allowed"

func (h *Handler) isValidRoadmapRequest(ctx echo.Context) error {
	areYouAHuman := ctx.FormValue("areYouAHuman")

	if areYouAHuman == iAmHuman {
		return nil
	}

	if areYouAHuman != "" {
		return fmt.Errorf(onlyHumansAreAllowed)
	}

	timeSpent := ctx.FormValue("ts")
	ts, err := strconv.ParseUint(timeSpent, 10, 64)
	if err != nil {
		return fmt.Errorf(onlyHumansAreAllowed)
	}

	if ts < 5 {
		return fmt.Errorf(onlyHumansAreAllowed)
	}

	return nil
}

func (h *Handler) isValidRoadmap(r Roadmap, dateFormat string) error {
	for _, r := range r.Projects {
		if r.Dates != nil && r.Dates.EndAt.Before(r.Dates.StartAt) {
			return fmt.Errorf(
				"end at before start at. start at: %s, end at: %s",
				r.Dates.StartAt.Format(dateFormat),
				r.Dates.EndAt.Format(dateFormat),
			)
		}
	}

	return nil
}

func (h *Handler) GetRoadmapImage(ctx echo.Context) error {
	format, err := NewFormatType(ctx.Param("format"))
	if err != nil {
		h.Logger.Info("format is not supported", zap.Error(err))

		return err
	}

	fw, err := strconv.ParseUint(ctx.QueryParam("width"), 10, 64)
	if err != nil {
		fw = defaultSvgWidth
	}

	lh, err := strconv.ParseUint(ctx.QueryParam("lineHeight"), 10, 64)
	if err != nil {
		lh = defaultSvgLineHeight
	}

	fw, lh = GetCanvasSizes(fw, lh)

	roadmap, err := load(h.rw, h.cb, ctx.Param("identifier"))
	if err != nil {
		h.Logger.Info("failed to load roadmap", zap.Error(err))

		return ctx.HTML(herr.ToHttpCode(err, http.StatusInternalServerError), fmt.Sprintf("%v", err))
	}

	cvs := roadmap.ToVisual().Draw(float64(fw), float64(lh))

	img := RenderImg(cvs, format)

	setHeaderContentType(ctx.Response().Header(), format)

	err = ctx.String(http.StatusOK, string(img))

	_ = h.Logger.Sync()

	return err
}

func load(rw DbReadWriter, b code.Builder, identifier string) (*Roadmap, error) {
	if identifier == "" {
		return nil, nil
	}

	c, err := b.NewFromString(identifier)
	if err != nil {
		return nil, herr.NewHttpError(err, http.StatusBadRequest)
	}

	roadmap, err := rw.Get(c)
	if err != nil {
		return nil, herr.NewHttpError(err, http.StatusInternalServerError)
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
	PdfFormat FileFormat = "pdf"
	JpgFormat FileFormat = "jpg"
	GifFormat FileFormat = "gif"
)

func NewFormatType(t string) (FileFormat, error) {
	switch t {
	case "svg":
		return SvgFormat, nil
	case "pdf":
		return PdfFormat, nil
	case "png":
		return PngFormat, nil
	case "jpg":
		return JpgFormat, nil
	case "jpeg":
		return JpgFormat, nil
	case "gif":
		return GifFormat, nil
	}

	return "", fmt.Errorf("unsupported image format: %s", t)
}

func setHeaderContentType(header http.Header, fileFormat FileFormat) {
	switch fileFormat {
	case SvgFormat:
		header.Set(echo.HeaderContentType, "image/svg+xml")
	case PdfFormat:
		header.Set(echo.HeaderContentType, "application/pdf")
	case PngFormat:
		header.Set(echo.HeaderContentType, "image/png")
	case GifFormat:
		header.Set(echo.HeaderContentType, "image/gif")
	case JpgFormat:
		header.Set(echo.HeaderContentType, "image/jpeg")
	}
}

func RenderImg(cvs *canvas.Canvas, fileFormat FileFormat) []byte {
	var buf bytes.Buffer

	switch fileFormat {
	case SvgFormat:
		img := canvas.NewSVG(&buf, cvs.W, cvs.H)

		cvs.Render(img)

		img.Close()
	case PdfFormat:
		img := canvas.NewPDF(&buf, cvs.W, cvs.H)

		cvs.Render(img)

		img.Close()
	case PngFormat:
		w := rasterizer.PNGWriter(3.2)

		err := w(&buf, cvs)
		if err != nil {
			return nil
		}
	case GifFormat:
		options := &gif.Options{
			NumColors: 256,
		}
		w := rasterizer.GIFWriter(3.2, options)

		err := w(&buf, cvs)
		if err != nil {
			return nil
		}
	case JpgFormat:
		options := &jpeg.Options{
			Quality: 90,
		}
		w := rasterizer.JPGWriter(3.2, options)

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
