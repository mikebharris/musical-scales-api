package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/music/music"
	"github.com/mikebharris/music/scala"
)

const (
	defaultEqualTemperamentDivisions = 12
	defaultJustLimit                 = 5
	defaultDiatonicMode              = music.IonianMode
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

var scalaHeaders = map[string]string{
	"Content-Type": "text/plain",
}

type Handler struct {
}

type Scale struct {
	System      string   `json:"system"`
	Description string   `json:"description"`
	Intervals   []string `json:"intervals"`
}

func newScale(s music.Scale) Scale {
	return Scale{
		System:      s.System(),
		Description: s.Description(),
		Intervals:   s.IntervalStrings(),
	}
}

func (h Handler) HandleRequest(_ context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	q := request.QueryStringParameters

	var scale music.Scale

	switch q["tuningSystem"] {
	case "bachWellTemperament":
		scale = music.NewBachWohltemperierteKlavierScale()
	case "edo":
		scale = music.NewEqualTemperamentScale(uint(parseIntegerQueryParameter(q, "divisions", defaultEqualTemperamentDivisions)))
	case "extendedMeantone":
		scale = music.NewExtendedQuarterCommaMeantoneScale()
	case "harmonic":
		scale = music.NewHarmonicSeriesScale(uint(parseIntegerQueryParameter(q, "partials", 1)))
	case "justFromRatios":
		scale = music.NewJustIntonationChromaticScaleWithLimit(parseIntegerQueryParameter(q, "limit", defaultJustLimit))
	case "meantone":
		scale = music.NewQuarterCommaMeantoneScale()
	case "partch":
		scale = music.NewPartch43ToneScale()
	case "ptolemy":
		scale = music.NewIntenseDiatonicScale(validDiatonicModeOrDefault(parseStringQueryParameter(q, "mode", "")))
	case "pythagorean":
		scale = music.NewPythagoreanScale()
	case "pythagorean5":
		scale = music.New5LimitPythagoreanScale()
	case "saz":
		scale = music.NewSazScale()
	default:
		return errorResponse(http.StatusUnprocessableEntity, `{"error":"please provide a valid tuning system"}`), nil
	}

	if q["format"] == "scala" {
		body := scala.NewScalaScaleFileFromScale(makeScalaFilename(q), scale)
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Headers: scalaHeaders, Body: body}, nil
	}

	body, _ := json.Marshal(newScale(scale))
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(body)}, nil
}

func makeScalaFilename(q map[string]string) string {
	if q["mode"] != "" {
		return q["tuningSystem"] + "-" + strings.ToLower(parseStringQueryParameter(q, "mode", defaultDiatonicMode.String())) + ".scl"
	}
	if q["tuningSystem"] == "edo" {
		return strconv.Itoa(parseIntegerQueryParameter(q, "divisions", defaultEqualTemperamentDivisions)) + "-edo.scl"
	}
	return q["tuningSystem"] + ".scl"
}

func parseStringQueryParameter(q map[string]string, key string, fallback string) string {
	if q[key] == "" {
		return fallback
	}
	return q[key]
}

func validDiatonicModeOrDefault(mode string) music.MusicalMode {
	if mode == "" || !music.MusicalMode(mode).IsDiatonic() {
		return defaultDiatonicMode
	}
	return music.MusicalMode(mode)
}

func errorResponse(status int, body string) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{StatusCode: status, Headers: headers, Body: body}
}

func parseIntegerQueryParameter(q map[string]string, key string, fallback int) int {
	if q[key] == "" {
		return fallback
	}
	atoi, err := strconv.Atoi(q[key])
	if err != nil || atoi <= 0 {
		return fallback
	}
	return atoi
}
