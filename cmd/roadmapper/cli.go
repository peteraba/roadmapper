package roadmapper

import (
	"time"

	"github.com/peteraba/roadmapper/pkg/roadmap"
	"go.uber.org/zap"
)

// Render renders a roadmap
func Render(rw roadmap.FileReadWriter, l *zap.Logger, content, output string, fileFormat, dateFormat, baseUrl string, fw, lh uint64) error {
	format, err := roadmap.NewFormatType(fileFormat)
	if err != nil {
		l.Info("format is not supported", zap.Error(err))

		return err
	}

	fw, lh = roadmap.GetCanvasSizes(fw, lh)

	r := roadmap.Content(content).ToRoadmap(0, nil, "", dateFormat, baseUrl, time.Now())
	cvs := r.ToVisual().Draw(float64(fw), float64(lh))
	img := roadmap.RenderImg(cvs, format)

	err = rw.Write(output, string(img))

	return err
}
