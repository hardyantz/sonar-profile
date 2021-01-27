# Sonar status
Library for send comment in pull request from sonarqube profile.

# Requirement

- Sonarqube token
- Sonarqube base URL
- Sonarqube project key
- Github username
- Github password / token
- Github Pull Request URL

# How To 

```go
package main

import (
    github "github.com/hardyantz/sonar-profile/github"
    services "github.com/hardyantz/sonar-profile/sonar"
    sonar "github.com/hardyantz/sonar-profile/sonar"
)

func main() {

    githubService := github.NewService("your_username", "your_password")
    sonarService := sonar.NewService("sonar_token", "sonar_base_url")
    service := services.NewStatus(githubService, sonarService)
    _ = service.Send(projectKey, urlPr)
    
}

```

# Result 

<img width="979" alt="Screen Shot 2020-12-08 at 15 39 24" src="https://user-images.githubusercontent.com/3294869/101460199-ddc6a880-396b-11eb-9d26-679e5c6d1a5e.png">


# Reference 

- generate sonar token : https://docs.sonarqube.org/latest/user-guide/user-token/
- github token : https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token