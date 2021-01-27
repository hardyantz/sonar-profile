package sonar

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/eapache/go-resiliency.v1/retrier"
)

type SearchProfile struct {
	Analyses []Analyses `json:"analyses"`
}

type Analyses struct {
	Key                         string `json:"key"`
	Date                        string `json:"date"`
	ProjectVersion              string `json:"projectVersion"`
	ManualNewCodePeriodBaseline bool   `json:"manualNewCodePeriodBaseline"`
	Revision                    string `json:"revision"`
}

type Response struct {
	ProjectStatus ProjectStatus `json:"projectStatus"`
}

type ProjectStatus struct {
	Status            string      `json:"status"`
	Conditions        []Condition `json:"conditions"`
	Periods           []Period    `json:"periods"`
	IgnoredConditions bool        `json:"ignoredConditions"`
}

type Condition struct {
	Status         string `json:"status"`
	MetricKey      string `json:"metricKey"`
	Comparator     string `json:"comparator"`
	PeriodIndex    int    `json:"periodIndex"`
	ErrorThreshold string `json:"errorThreshold"`
	ActualValue    string `json:"actualValue"`
}

type Period struct {
	Index int    `json:"index"`
	Mode  string `json:"mode"`
	Date  string `json:"date"`
}

type service struct {
	URL   string
	Token string
}

type Service interface {
	GetProfile(projectKey, shaCommit string) (Response, error)
	SearchProfile(projectKey string) (SearchProfile, error)
}

func NewService(token, url string) Service {
	return &service{
		URL:   url,
		Token: token,
	}
}

func (s *service) SearchProfile(projectKey string) (SearchProfile, error) {
	var result SearchProfile

	url := fmt.Sprintf("%s/api/project_analyses/search?project=%s", s.URL, projectKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	req.SetBasicAuth(s.Token, "")

	err = s.fetchApi(req, &result)
	if err != nil {
		return result, err
	}

	if len(result.Analyses) == 0 {
		return result, errors.New("analyses not found")
	}

	return result, nil
}

func (s *service) fetchApi(req *http.Request, response interface{}) error {
	req.SetBasicAuth(s.Token, "")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(respBody, &response); err != nil {
		return err
	}

	return nil
}

func (s *service) GetProfile(projectKey, shaCommit string) (Response, error) {
	var v Response
	var sonarAnalyses Analyses

	r := retrier.New(retrier.ConstantBackoff(5, 5*time.Second), nil)
	attempt := 0
	err := r.Run(func() error {
		searchProfile, err := s.SearchProfile(projectKey)
		if err != nil {
			attempt++
			return err
		}

		for _, analyse := range searchProfile.Analyses {
			if analyse.Revision == shaCommit {
				sonarAnalyses = analyse
				break
			}
		}

		if sonarAnalyses.Key == "" {
			err = errors.New("error sha commit not found")
		}

		attempt++
		return err
	})

	if err != nil {
		return v, err
	}

	req, err := http.NewRequest(http.MethodGet, s.URL+"/api/qualitygates/project_status?analysisId="+sonarAnalyses.Key, nil)
	if err != nil {
		return v, err
	}

	if err = s.fetchApi(req, &v); err != nil {
		return v, err
	}

	return v, nil
}
