# Overview
Generate CLI for sonar-profile. 

# How-to

run this command
```shell script
$ go install $GOPATH/bin/sonar-profile
```

execute CLI for sending sonar analysis to github

```shell script
$ sonar-profile -ghPass {github password} -ghUname {github username} -sonarKey {sonar project key} -sonarToken {sonar token authentication} -sonarURL {sonar base url} -urlPR {github pull request URL}
```

type help for detail 

```shell script
$ sonar-profile --help
```