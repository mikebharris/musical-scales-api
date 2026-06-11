package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/music/music"
	"github.com/stretchr/testify/assert"
)

func Test_ShouldReturnErrorWhenTuningSystemIsNotProvided(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide a valid tuning system"}`}, response)
}

func Test_ShouldReturnErrorWhenTuningSystemIsInvalid(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "invalid"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide a valid tuning system"}`}, response)
}

func Test_ShouldReturnScaleForSazTuning(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "saz"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "Saz", scale.System)
	assert.Equal(t, "Turkish Saz tuning ratios.", scale.Description)

	// a minimum test to see that it returns something, but we don't want to re-test the music module
	assert.Equal(t, 18, len(scale.Intervals))
	assert.Equal(t, "1:1", scale.Intervals[0])
	assert.Equal(t, "2:1", scale.Intervals[17])
}

func Test_ShouldReturnScaleForPythagoreanTuning(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "pythagorean"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "Pythagorean", scale.System)
	assert.Equal(t, "3-limit Pythagorean ratios.", scale.Description)
}

func Test_ShouldReturnScaleForFiveLimitPythagoreanTuning(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "pythagorean5"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "5-limit Pythagorean", scale.System)
	assert.Equal(t, "5-limit just intonation pure ratios chromatic scale derived from applying syntonic comma to Pythagorean ratios.", scale.Description)
}

func Test_ShouldReturnFiveLimitJustScaleFromPureRatios(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "justFromRatios"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "5-limit Just Intonation", scale.System)
	assert.Equal(t, "Just Intonation chromatic scale based on 5-limit pure ratios.", scale.Description)
	assert.Equal(t, 13, len(scale.Intervals))
}

func Test_ShouldReturnJustScaleFromPureRatiosForGivenLimit(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "justFromRatios", "limit": "7"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "7-limit Just Intonation", scale.System)
	assert.Equal(t, "Just Intonation chromatic scale based on 7-limit pure ratios.", scale.Description)
	assert.Equal(t, 13, len(scale.Intervals))
}

func Test_ShouldReturnFiveLimitJustScaleFromPureRatiosWhenLimitIsInvalid(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "justFromRatios", "limit": "X"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "5-limit Just Intonation", scale.System)
	assert.Equal(t, "Just Intonation chromatic scale based on 5-limit pure ratios.", scale.Description)
	assert.Equal(t, 13, len(scale.Intervals))
}

func Test_ShouldReturnFiveLimitJustScaleFromPureRatiosWhenLimitIsNegative(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "justFromRatios", "limit": "-1"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "5-limit Just Intonation", scale.System)
}

func Test_ShouldReturnFiveLimitJustScaleFromPureRatiosWhenLimitIsZero(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "justFromRatios", "limit": "0"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "5-limit Just Intonation", scale.System)
}

func Test_ShouldReturnPtolemyIntenseDiatonicScaleInIonianModeByDefault(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "ptolemy"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "Ptolemy Intense Diatonic", scale.System)
	assert.Equal(t, "Ptolemy's 5-limit intense diatonic scale in Ionian mode.", scale.Description)
	assert.Equal(t, 8, len(scale.Intervals))
}

func Test_ShouldReturnPtolemyIntenseDiatonicScaleInIonianWhenInvalidModeIsProvided(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "ptolemy", "mode": "Athenian"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "Ptolemy Intense Diatonic", scale.System)
	assert.Equal(t, "Ptolemy's 5-limit intense diatonic scale in Ionian mode.", scale.Description)
	assert.Equal(t, 8, len(scale.Intervals))
}

func Test_ShouldReturnPtolemyIntenseDiatonicScaleForProvidedMusicalMode(t *testing.T) {
	// Given
	tuningSystem := "ptolemy"
	mode := music.LydianMode

	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": tuningSystem, "mode": mode.String()},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "Ptolemy's 5-limit intense diatonic scale in Lydian mode.", scale.Description)
}

func Test_ShouldBachWellTemperamentScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "bachWellTemperament"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "Bach's Well-Tempered Tuning", scale.System)
	assert.Equal(t, "Derived from Lehman's decoding of Bach's Well-Tempered tuning, using sixth-comma, twelfth-comma, and pure fifths.", scale.Description)
}

func Test_ShouldReturnQuarterCommaMeantoneScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "meantone"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "Quarter-Comma Meantone", scale.System)
	assert.Equal(t, "Meantone temperament achieved by narrowing of fifths by 0.25 of a syntonic comma (81/80).", scale.Description)
}

func Test_ShouldReturnExtendedQuarterCommaMeantoneScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "extendedMeantone"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "Extended Quarter-Comma Meantone", scale.System)
	assert.Equal(t, "Meantone temperament achieved by narrowing of fifths by 0.25 of a syntonic comma (81/80).", scale.Description)
}

func Test_ShouldReturnTwelveToneEqualTemperamentScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "edo"}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "Equal Temperament", scale.System)
	assert.Equal(t, "12-tone equal temperament.", scale.Description)
	assert.Equal(t, 13, len(scale.Intervals))
}

func Test_ShouldReturnNineteenToneEqualTemperamentScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "edo", "divisions": "19"}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "19-tone equal temperament.", scale.Description)
	assert.Equal(t, 20, len(scale.Intervals))
}

func Test_ShouldReturnTwelveToneEqualTemperamentScaleWhenInvalidNumberOfDivisionsOfOctaveProvided(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "edo", "divisions": "-1"}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	var scale Scale
	_ = json.Unmarshal([]byte(response.Body), &scale)
	assert.Equal(t, "12-tone equal temperament.", scale.Description)
}

func Test_ShouldReturnScalaFileContentsForRequestedPtolemyScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "ptolemy", "mode": music.PhrygianMode.String(), "format": "scala"}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, scalaHeaders, response.Headers)

	contents := strings.Split(response.Body, "\n")
	assert.Equal(t, "! ptolemy-phrygian.scl", contents[0])
	assert.Equal(t, "! generated using github.com/mikebharris/music/scala", contents[1])
	assert.Equal(t, "!", contents[2])
	assert.Equal(t, "Ptolemy Intense Diatonic scale using Ptolemy's 5-limit intense diatonic scale in Phrygian mode.", contents[3])
	assert.Equal(t, "7", contents[4])
	assert.Equal(t, "2/1", contents[11])
}

func Test_ShouldReturnScalaFileContentsForRequestedTemperedScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "edo", "divisions": "19", "format": "scala"}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, scalaHeaders, response.Headers)

	contents := strings.Split(response.Body, "\n")
	assert.Equal(t, "! 19-edo.scl", contents[0])
	assert.Equal(t, "! generated using github.com/mikebharris/music/scala", contents[1])
	assert.Equal(t, "!", contents[2])
	assert.Equal(t, "Equal Temperament scale using 19-tone equal temperament.", contents[3])
	assert.Equal(t, "19", contents[4])
	assert.Equal(t, "1200.00", contents[23])
}
