package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/testcontainernetwork-go"
	"github.com/stretchr/testify/assert"

	"github.com/cucumber/godog"
)

type scaleResponse struct {
	System      string   `json:"system"`
	Description string   `json:"description"`
	Intervals   []string `json:"intervals"`
}

func TestFeatures(t *testing.T) {
	var steps steps
	steps.t = t

	suite := godog.TestSuite{
		TestSuiteInitializer: func(ctx *godog.TestSuiteContext) {
			ctx.BeforeSuite(steps.startContainerNetwork)
			ctx.AfterSuite(steps.stopContainerNetwork)
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Step(`^I request a scale for a Turkish saz$`, steps.iRequestAScaleForATurkishSaz)
			ctx.Step(`^I am provided with the just interval ratios for the saz$`, steps.iAmProvidedWithTheJustIntervalRatiosForTheSaz)

			ctx.Step(`^I request a Scala file for a Turkish saz scale$`, steps.iRequestAScalaFileForATurkishSazScale)
			ctx.Step(`^I am provided with a Scala file containing the just interval ratios for the saz$`, steps.iAmProvidedWithTheScalaFileForTheSaz)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

type steps struct {
	t                   *testing.T
	networkOfContainers testcontainernetwork.NetworkOfDockerContainers
	lambdaContainer     testcontainernetwork.LambdaDockerContainer
}

var responseFromLambda events.LambdaFunctionURLResponse

func (s *steps) startContainerNetwork() {
	s.lambdaContainer = testcontainernetwork.LambdaDockerContainer{
		Config: testcontainernetwork.LambdaDockerContainerConfig{
			Hostname:    "lambda",
			Executable:  "../main",
			Environment: map[string]string{},
		},
	}

	s.networkOfContainers =
		testcontainernetwork.NetworkOfDockerContainers{}.
			WithDockerContainer(&s.lambdaContainer)
	_ = s.networkOfContainers.StartWithDelay(5 * time.Second)
}

func (s *steps) stopContainerNetwork() {
	if err := s.networkOfContainers.Stop(); err != nil {
		log.Fatalf("stopping docker containers: %v", err)
	}
}

func (s *steps) iRequestAScaleForATurkishSaz() {
	s.invokeLambdaUsingRequest(events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "saz"},
	})
}

func (s *steps) iRequestAScalaFileForATurkishSazScale() {
	s.invokeLambdaUsingRequest(events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"tuningSystem": "saz", "format": "scala"},
	})
}

func (s *steps) invokeLambdaUsingRequest(request events.LambdaFunctionURLRequest) {
	requestJsonBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("marshalling lambda request %v", err)
	}
	response, err := http.Post(s.lambdaContainer.InvocationUrl(), "application/json", bytes.NewReader(requestJsonBytes))
	if err != nil {
		log.Fatalf("triggering lambda: %v", err)
	}

	if response.StatusCode != 200 {
		log.Fatalf("invoking Lambda: %d", response.StatusCode)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(response.Body); err != nil {
		log.Fatalf("reading response body: %v", err)
	}

	if err := json.Unmarshal(buf.Bytes(), &responseFromLambda); err != nil {
		log.Fatalf("unmarshalling response: %v", err)
	}
}

func (s *steps) iAmProvidedWithTheJustIntervalRatiosForTheSaz() error {
	assert.Equal(s.t, "application/json", responseFromLambda.Headers["Content-Type"])
	assert.Equal(s.t, http.StatusOK, responseFromLambda.StatusCode)

	var scale scaleResponse
	if err := json.Unmarshal([]byte(responseFromLambda.Body), &scale); err != nil {
		return fmt.Errorf("unmarshalling result: %s", err)
	}

	assert.Equal(s.t, "Saz", scale.System)
	assert.Equal(s.t, "Turkish Saz tuning ratios.", scale.Description)
	assert.Equal(s.t, 18, len(scale.Intervals))
	return nil
}

func (s *steps) iAmProvidedWithTheScalaFileForTheSaz() error {
	assert.Equal(s.t, "text/plain", responseFromLambda.Headers["Content-Type"])
	assert.Equal(s.t, http.StatusOK, responseFromLambda.StatusCode)

	contents := strings.Split(responseFromLambda.Body, "\n")
	assert.Equal(s.t, "! saz.scl", contents[0])
	assert.Equal(s.t, "! generated using github.com/mikebharris/music/scala", contents[1])
	assert.Equal(s.t, "!", contents[2])
	assert.Equal(s.t, "Saz scale using Turkish Saz tuning ratios.", contents[3])
	assert.Equal(s.t, "17", contents[4])
	assert.Equal(s.t, "2/1", contents[21])

	return nil
}
