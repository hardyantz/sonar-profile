package sonarcomment

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hardyantz/sonar-profile/github"
	"github.com/hardyantz/sonar-profile/sonar"
)

type Profile interface {
	Send(projectKey, UrlPR string) error
}

type profile struct {
	github github.Service
	sonar  sonar.Service
}

func NewStatus(github github.Service, sonar sonar.Service) Profile {
	return &profile{
		github,
		sonar,
	}
}

func (s *profile) Send(projectKey, UrlPR string) error {
	resDetailPR, err := s.github.GetPRDetail(UrlPR)
	if err != nil {
		return err
	}

	resSonar, err := s.sonar.GetProfile(projectKey, resDetailPR.Head.Sha)
	if err != nil {
		return err
	}

	profile := s.parseResponse(resSonar)
	if profile == "" {
		return errors.New("empty message")
	}

	err = s.github.SendComment(resDetailPR.IssueURL, profile)
	if err != nil {
		return err
	}

	return nil
}

func (s *profile) parseResponse(sonarProfile sonar.Response) string {
	var profile string

	if sonarProfile.ProjectStatus.Status == "" {
		return ""
	}

	overall := `:heavy_check_mark: `
	if sonarProfile.ProjectStatus.Status != "OK" {
		overall = `:no_entry_sign: `
	}

	profile = fmt.Sprintf("%s Sonarqube analysis : **%s**\n", overall, sonarProfile.ProjectStatus.Status)

	for _, v := range sonarProfile.ProjectStatus.Conditions {
		metric := strings.Replace(v.MetricKey, "_", " ", 5)
		conStatus := `:white_check_mark:`
		if v.Status != "OK" {
			conStatus = `:x:`
		}

		var comparator string
		switch v.Comparator {
		case "GT":
			comparator = " > "
		case "LT":
			comparator = " < "
		case "LTE":
			comparator = " <= "
		case "GTE":
			comparator = " >= "
		}

		status := fmt.Sprintf("- %s **%s** on new code is **%s**, actual value **%s**, error if %s %s \n", conStatus, strings.Title(metric), v.Status, v.ActualValue, comparator, v.ErrorThreshold)
		profile += status
	}

	return profile
}
