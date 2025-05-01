package tests

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/kootic/koogo/internal/dto"
)

func TestKooUser(t *testing.T) {
	t.Parallel()

	responses := make([]testResponse, 0)

	plan := testPlan{
		stepsFactory: []testStepFactory{
			func(responses []testResponse) testStep {
				return testStep{
					name:        "create koo user",
					path:        "/api/v1/koo/users",
					method:      http.MethodPost,
					contentType: "application/json",
					body: map[string]any{
						"firstName": "John Doe",
					},
					expectStatusCode: http.StatusOK,
					validateResponse: func(t *testing.T, response testResponse) error {
						kooUser, err := decodeTestResponse[dto.KooUserResponse](response)
						if err != nil {
							return err
						}

						if kooUser.ID == uuid.Nil {
							return errors.New("user id is nil")
						}

						return nil
					},
				}
			},
			func(responses []testResponse) testStep {
				newUser, err := decodeTestResponse[dto.KooUserResponse](responses[0])
				if err != nil {
					t.Fatalf("failed to get user from previous step")
					return testStep{}
				}

				return testStep{
					name:             "get koo user",
					path:             "/api/v1/koo/users/" + newUser.ID.String(),
					method:           http.MethodGet,
					expectStatusCode: http.StatusOK,
					validateResponse: func(t *testing.T, response testResponse) error {
						kooUser, err := decodeTestResponse[dto.KooUserResponse](response)
						if err != nil {
							return err
						}

						if kooUser.ID != newUser.ID {
							return errors.New("user id does not match")
						}

						if kooUser.FirstName != newUser.FirstName {
							return errors.New("user first name does not match")
						}

						return nil
					},
				}
			},
		},
	}

	runTestPlan(t, plan, responses)
}
