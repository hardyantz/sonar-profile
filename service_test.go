package sonarcomment

import (
	"errors"
	"testing"

	"github.com/hardyantz/sonar-profile/github"
	githubMock "github.com/hardyantz/sonar-profile/github/mocks"
	"github.com/hardyantz/sonar-profile/sonar"
	sonarMock "github.com/hardyantz/sonar-profile/sonar/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStatus(t *testing.T) {
	githubAuth := github.NewService("", "")
	sonarAuth := sonar.NewService("", "")
	actual := NewStatus(githubAuth, sonarAuth)
	assert.Equal(t, &profile{
		github: githubAuth,
		sonar:  sonarAuth,
	}, actual)
}

func TestProfileSend(t *testing.T) {

	tests := []struct {
		name          string
		response      sonar.Response
		isError       bool
		githubError   error
		sonarError    error
		ghPRDetailErr error
	}{
		{
			name: "send comment success",
			response: sonar.Response{
				ProjectStatus: sonar.ProjectStatus{
					Status: "OK",
					Conditions: []sonar.Condition{
						{
							Status:         "OK",
							MetricKey:      "new_code",
							ErrorThreshold: "1",
							ActualValue:    "1",
						},
					},
				},
			},
			isError: false,
		},
		{
			name:       "get profile failed",
			response:   sonar.Response{},
			sonarError: errors.New("failed get profile"),
			isError:    true,
		},
		{
			name: "send comment failed",
			response: sonar.Response{
				ProjectStatus: sonar.ProjectStatus{
					Status: "OK",
					Conditions: []sonar.Condition{
						{
							Status:         "OK",
							MetricKey:      "new_code",
							ErrorThreshold: "1",
							ActualValue:    "1",
						},
					},
				},
			},
			githubError: errors.New("failed sendProfile"),
			isError:     true,
		},
		{
			name:     "sonar profile empty",
			response: sonar.Response{},
			isError:  true,
		},
		{
			name:          "get PR detail failed",
			response:      sonar.Response{},
			ghPRDetailErr: errors.New("failed get PR detail"),
			isError:       true,
		},
		{
			name: "parse profile",
			response: sonar.Response{
				ProjectStatus: sonar.ProjectStatus{
					Status: "Not OK",
					Conditions: []sonar.Condition{
						{
							Status:         "not OK",
							MetricKey:      "new_code",
							ErrorThreshold: "1",
							ActualValue:    "1",
							Comparator:     "GT",
						},
						{
							Status:         "OK",
							MetricKey:      "new_reliability",
							ErrorThreshold: "2",
							ActualValue:    "2",
							Comparator:     "LT",
						},
						{
							Status:         "OK",
							MetricKey:      "new_coverage",
							ErrorThreshold: "80",
							ActualValue:    "80",
							Comparator:     "LTE",
						},
						{
							Status:         "OK",
							MetricKey:      "new_bugs",
							ErrorThreshold: "2",
							ActualValue:    "2",
							Comparator:     "GTE",
						},
					},
				},
			},
			isError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			githubService := new(githubMock.Service)
			sonarService := new(sonarMock.Service)
			sonarService.On("GetProfile", mock.Anything, mock.Anything).Return(tt.response, tt.sonarError)
			githubService.On("SendComment", mock.Anything, mock.Anything).Return(tt.githubError)
			githubService.On("GetPRDetail", mock.Anything).Return(github.GhPRResponse{}, tt.ghPRDetailErr)
			projectKey := ""
			urlPr := ""
			res := NewStatus(githubService, sonarService)
			err := res.Send(projectKey, urlPr)
			if tt.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
