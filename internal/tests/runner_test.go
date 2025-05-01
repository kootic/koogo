package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type testPlan struct {
	stepsFactory []testStepFactory
}

// testStepFactory is a function that receives the responses from previous steps
// and returns a new testStep in order to build steps that depend on previous steps.
//
// For example, if the first step creates a user, and the second step gets the user,
// the second step's testStepFactory will a slice of testResponse with length 1,
// where the first element is the response from the first step.
type testStepFactory func(responses []testResponse) testStep

// testResponse is a struct that contains the raw response body from a test step
// and anything else that might be needed to build the following steps can be added here.
type testResponse struct {
	rawBody []byte
}

// decodeTestResponse decodes the raw response body into a struct.
func decodeTestResponse[T any](r testResponse) (T, error) {
	var decoded T

	err := json.Unmarshal(r.rawBody, &decoded)
	if err != nil {
		return decoded, err
	}

	return decoded, nil
}

type testStep struct {
	name        string
	path        string
	method      string
	contentType string
	// body should be a map[string]any most of the time so we can also validate the JSON marshalling
	body             any
	expectStatusCode int
	validateResponse func(t *testing.T, response testResponse) error
}

func runTestPlan(t *testing.T, plan testPlan, responses []testResponse) {
	t.Helper()

	for _, stepFactory := range plan.stepsFactory {
		responses = runTestStep(t, stepFactory, responses)
	}
}

func runTestStep(t *testing.T, stepFactory testStepFactory, responses []testResponse) []testResponse {
	t.Helper()

	step := stepFactory(responses)

	var newResponse testResponse

	t.Run(step.name, func(t *testing.T) {
		// Marshal the body if it exists
		var body io.Reader
		if step.body != nil {
			jsonBody, err := json.Marshal(step.body)
			if err != nil {
				t.Fatalf("Failed to marshal body: %v", err)
			}

			body = bytes.NewReader(jsonBody)
		}

		req, _ := http.NewRequest(step.method, step.path, body)
		if step.contentType != "" {
			req.Header.Set("Content-Type", step.contentType)
		}

		// Send the request and get the response
		resp, err := testApp.FiberApp().Test(req)
		if err != nil {
			t.Fatalf("Failed to send test request: %v", err)
		}
		defer resp.Body.Close() //nolint:errcheck

		// Read the response body
		newResponse, err = getTestResponse(resp)
		if err != nil {
			t.Fatalf("Failed to get test response: %v", err)
		}

		// Run validations
		if step.expectStatusCode != 0 {
			if resp.StatusCode != step.expectStatusCode {
				t.Fatalf("Expected status %d, got %d", step.expectStatusCode, resp.StatusCode)
			}
		}

		if step.validateResponse != nil {
			err := step.validateResponse(t, newResponse)
			if err != nil {
				t.Fatalf("Failed to validate response: %v", err)
			}
		}
	})

	return append(responses, newResponse)
}

// getTestResponse reads the response body and returns a testResponse.
func getTestResponse(resp *http.Response) (testResponse, error) {
	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return testResponse{}, err
	}

	return testResponse{
		rawBody: respBody,
	}, nil
}
