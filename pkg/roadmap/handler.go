package roadmap

import (
	"bytes"
	"fmt"
	"image/gif"
	"image/jpeg"
	"net/http"
	"net/url"
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

	roadmap, c, err := load(h.rw, h.cb, identifier)
	if err != nil {
		return ctx.HTML(c, err.Error())
	}

	output, err := bootstrapRoadmap(roadmap, h.appVersion, h.matomoDomain, h.docBaseURL, ctx.Request().RequestURI, h.selfHosted)
	if err != nil {
		h.Logger.Info("failed to create HTML response", zap.Error(err))

		return ctx.HTML(herr.ToHttpCode(err, http.StatusInternalServerError), err.Error())
	}

	return ctx.HTML(http.StatusOK, output)
}

func (h *Handler) CreateRoadmapHTML(ctx echo.Context) error {
	prevID, err := h.getPrevID(ctx.Param("identifier"))
	if err != nil {
		h.Logger.Info("failed to parse the identifier parameter", zap.Error(err))

		return ctx.Redirect(http.StatusSeeOther, "/?error="+url.QueryEscape(err.Error()))
	}

	err = h.isValidRoadmapRequest(ctx)
	if err != nil {
		h.Logger.Info("not a valid roadmap request", zap.Error(err))

		return ctx.Redirect(http.StatusSeeOther, "/?error="+url.QueryEscape(err.Error()))
	}

	title := ctx.FormValue("title")
	content := ctx.FormValue("txt")
	dateFormat := ctx.FormValue("dateFormat")
	baseURL := ctx.FormValue("baseUrl")
	now := time.Now()

	roadmap := Content(content).ToRoadmap(code.NewCode64().ID(), prevID, title, dateFormat, baseURL, now)

	err = h.rw.Write(h.cb, roadmap)
	if err != nil {
		h.Logger.Info("failed to write the new roadmap", zap.Error(err))

		return ctx.HTML(http.StatusMethodNotAllowed, err.Error())
	}

	code, err := h.cb.NewFromID(roadmap.ID)
	if err != nil {
		h.Logger.Info("failed to generate the new  url", zap.Error(err))

		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}

	newURL := fmt.Sprintf("/%s", code.String())

	return ctx.Redirect(http.StatusSeeOther, newURL)
}

func (h *Handler) getPrevID(identifier string) (*uint64, error) {
	if identifier == "" {
		return nil, nil
	}

	code, err := h.cb.NewFromString(identifier)
	if err != nil {
		return nil, herr.NewHttpError(err, http.StatusNotFound)
	}

	n := code.ID()

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

	roadmap, c, err := load(h.rw, h.cb, ctx.Param("identifier"))
	if err != nil {
		h.Logger.Info("failed to load roadmap", zap.Error(err))

		return ctx.HTML(c, fmt.Sprintf("%v", err))
	}

	cvs := roadmap.ToVisual().Draw(float64(fw), float64(lh))

	img := RenderImg(cvs, format)

	setHeaderContentType(ctx.Response().Header(), format)

	return ctx.String(http.StatusOK, string(img))
}

func load(rw DbReadWriter, b code.Builder, identifier string) (*Roadmap, int, error) {
	if identifier == "" {
		return nil, 0, nil
	}

	c, err := b.NewFromString(identifier)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	roadmap, err := rw.Read(c)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return roadmap, 0, nil
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
		defer img.Close()

		cvs.Render(img)
	case PdfFormat:
		img := canvas.NewPDF(&buf, cvs.W, cvs.H)
		defer img.Close()

		cvs.Render(img)
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
