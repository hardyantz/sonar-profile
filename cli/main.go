package main

import (
	"flag"
	"fmt"
	"os"

	sonarcomment "github.com/hardyantz/sonar-profile"
	"github.com/hardyantz/sonar-profile/github"
	"github.com/hardyantz/sonar-profile/sonar"
)

func main() {
	var urlPR, sonarKey, sonarURL, sonarToken, ghPass, ghUName string
	var errArg bool
	flag.StringVar(&urlPR, "urlPR", "", "Github Pull Request URL, ex: https://github.com/repository/pull/1")
	flag.StringVar(&sonarKey, "sonarKey", "", "Sonar project key")
	flag.StringVar(&sonarURL, "sonarURL", "", "Sonar base url")
	flag.StringVar(&sonarToken, "sonarToken", "", "Sonar token authentication")
	flag.StringVar(&ghPass, "ghPass", "", "Github password or github token authentication")
	flag.StringVar(&ghUName, "ghUName", "", "Github username")
	flag.Parse()

	if ghPass == "" {
		fmt.Println("error: GITHUB_PASSWORD or github token required")
		errArg = true
	}

	if ghUName == "" {
		fmt.Println("error: GITHUB_USERNAME required")
		errArg = true
	}

	if sonarToken == "" {
		fmt.Println("error: SONAR_TOKEN required")
		errArg = true
	}

	if sonarURL == "" {
		fmt.Println("error: SONAR_URL required")
		errArg = true
	}

	if sonarKey == "" {
		fmt.Println("error: SONAR_PROJECT_KEY required")
		errArg = true
	}

	if errArg {
		fmt.Println(`type --help for detail`)
		os.Exit(0)
	}

	githubAuth := github.NewService(ghUName, ghPass)
	sonarAuth := sonar.NewService(sonarToken, sonarURL)
	sonarComment := sonarcomment.NewStatus(githubAuth, sonarAuth)
	err := sonarComment.Send(sonarKey, urlPR)
	if err != nil {
		fmt.Println("error: failed processing sonar profile with error :" + err.Error())
	} else {
		fmt.Println("success: processing sonar profile")
	}
	os.Exit(0)
}
