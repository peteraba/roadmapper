package main

import (
	"time"

	"go.uber.org/zap"

	"github.com/peteraba/roadmapper/pkg/roadmap"
)

// Render renders a roadmap
func Render(io roadmap.IO, l *zap.Logger, content, output string, fileFormat, dateFormat, baseUrl string, fw, lh uint64, mt bool) error {
	format, err := roadmap.NewFormatType(fileFormat)
	if err != nil {
		l.Info("format is not supported", zap.Error(err))

		return err
	}

	fw, lh = roadmap.GetCanvasSizes(fw, lh)

	r := roadmap.Content(content).ToRoadmap(0, nil, "", dateFormat, baseUrl, time.Now())

	cvs := r.ToVisual().Draw(float64(fw), float64(lh), mt)

	img := roadmap.RenderImg(cvs, format)

	err = io.Write(output, string(img))

	return err
}
