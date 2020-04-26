package main

import (
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/peteraba/roadmapper/pkg/roadmap"
)

const (
	txt = `Monocle ipsum dolor sit amet
Ettinger punctual izakaya concierge [2020-02-02, 2020-02-20, 60%]
	Zürich Baggu bureaux [/issues/1]
		Toto Comme des Garçons liveable [2020-02-04, 2020-02-25, 100%, https://example.com/abc, /issues/2]
		Winkreative boutique St Moritz [2020-02-06, 2020-02-22, 55%, /issues/3]
	Toto joy perfect Porter [2020-02-25, 2020-03-01, 100%, |1]
Craftsmanship artisanal
	Marylebone exclusive [2020-03-03, 2020-03-10, 100%]
	Beams elegant destination [2020-03-08, 2020-03-12, 100%, |1]
	Winkreative ryokan hand-crafted [2020-03-13, 2020-03-31, 20%]
Nordic Toto first-class Singap
	Concierge cutting-edge Zürich global bureaux
		Sunspel sophisticated lovely uniforms [2020-03-17, 2020-03-31]
		Share blog post on social media [2020-03-17, 2020-03-31, 80%]
	Talk about the tool in relevant meetups [2020-04-01, 2020-06-15, 20%]
Melbourne handsome boutique
	Boutique magna iconic
		Carefully curated laborum destination [2020-03-28, 2020-05-01, 60%]
	Qui incididunt sleepy
		Scandinavian occaecat culpa [2020-03-26, 2020-04-01, 90%]
Hand-crafted K-pop boulevard
	Charming sed quality [2020-03-18, 2020-05-31, 20%]
	Sunspel alluring ut dolore [2020-04-15, 2020-04-30, 30%]
Business class Shinkansen [2020-04-01, 2020-05-31, 45%]
	Nisi excepteur hand-crafted hub
	Ettinger Airbus A380
Essential conversation bespoke
Muji enim

|Laboris ullamco
|Muji enim finest [2020-02-12, https://example.com/abc, bcdef]`
	txtBaseURL = "https://example.com/foo"
)

func TestIntegration_TextToRoadmap(t *testing.T) {
	now := time.Now()
	content := roadmap.Content(txt)

	roadmap := content.ToRoadmap(123, nil, "", "2006-01-02", txtBaseURL, now)

	actual := roadmap.String()

	assert.Equal(t, txt, actual)
}

func TestIntegration_TextToVisual(t *testing.T) {
	now := time.Now()
	content := roadmap.Content(txt)
	expectedProjectLength := 28
	expectedMilestoneLength := 2
	expectedDeadline1 := time.Date(2020, 3, 12, 0, 0, 0, 0, time.UTC)

	roadmap := content.ToRoadmap(123, nil, "", "2006-01-02", txtBaseURL, now)
	visualRoadmap := roadmap.ToVisual()

	assert.Len(t, visualRoadmap.Projects, expectedProjectLength)
	assert.Len(t, visualRoadmap.Milestones, expectedMilestoneLength)
	assert.Equal(t, expectedDeadline1.String(), visualRoadmap.Milestones[0].DeadlineAt.String())
}
