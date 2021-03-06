// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	sonar "github.com/hardyantz/sonar-profile/sonar"
	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// GetProfile provides a mock function with given fields: projectKey, shaCommit
func (_m *Service) GetProfile(projectKey string, shaCommit string) (sonar.Response, error) {
	ret := _m.Called(projectKey, shaCommit)

	var r0 sonar.Response
	if rf, ok := ret.Get(0).(func(string, string) sonar.Response); ok {
		r0 = rf(projectKey, shaCommit)
	} else {
		r0 = ret.Get(0).(sonar.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(projectKey, shaCommit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchProfile provides a mock function with given fields: projectKey
func (_m *Service) SearchProfile(projectKey string) (sonar.SearchProfile, error) {
	ret := _m.Called(projectKey)

	var r0 sonar.SearchProfile
	if rf, ok := ret.Get(0).(func(string) sonar.SearchProfile); ok {
		r0 = rf(projectKey)
	} else {
		r0 = ret.Get(0).(sonar.SearchProfile)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(projectKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
