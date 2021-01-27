package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type service struct {
	Username string
	Password string
	PRUrl    string
}

type GhPRResponse struct {
	URL      string `json:"url"`
	ID       int    `json:"id"`
	IssueURL string `json:"issue_url"`
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Head     struct {
		Sha string `json:"sha"`
	} `json:"head"`
}

type Service interface {
	SendComment(issueUrl, commentBody string) error
	GetPRDetail(urlPR string) (GhPRResponse, error)
}

func NewService(username, password string) Service {
	return &service{
		Username: username,
		Password: password,
	}
}

func (s *service) SendComment(issueUrl, commentBody string) error {
	str, _ := json.Marshal(map[string]string{"body": commentBody})

	byteReq := strings.NewReader(string(str))

	req, err := http.NewRequest(http.MethodPost, issueUrl+"/comments", byteReq)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")

	req.SetBasicAuth(s.Username, s.Password)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("post comment failed with error %s", res.Status)
	}

	return nil
}

func (s *service) GetPRDetail(urlPR string) (GhPRResponse, error) {
	apiPR := s.parseUrlPR(urlPR)
	var result GhPRResponse
	req, err := http.NewRequest(http.MethodGet, apiPR, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")

	req.SetBasicAuth(s.Username, s.Password)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	// unmarshal to our target
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return result, err
	}

	if res.StatusCode != http.StatusOK {
		return result, fmt.Errorf("get PR detail failed with error %s", res.Status)
	}

	return result, nil
}

func (s *service) parseUrlPR(u string) string {
	issueURL := strings.Replace(u, "/pull/", "/pulls/", 1)
	issueURL = strings.Replace(issueURL, "https://github.com", "https://api.github.com/repos", 1)
	return issueURL
}
