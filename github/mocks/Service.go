// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	github "github.com/hardyantz/sonar-profile/github"
	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// GetPRDetail provides a mock function with given fields: urlPR
func (_m *Service) GetPRDetail(urlPR string) (github.GhPRResponse, error) {
	ret := _m.Called(urlPR)

	var r0 github.GhPRResponse
	if rf, ok := ret.Get(0).(func(string) github.GhPRResponse); ok {
		r0 = rf(urlPR)
	} else {
		r0 = ret.Get(0).(github.GhPRResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(urlPR)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendComment provides a mock function with given fields: issueUrl, commentBody
func (_m *Service) SendComment(issueUrl string, commentBody string) error {
	ret := _m.Called(issueUrl, commentBody)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(issueUrl, commentBody)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
