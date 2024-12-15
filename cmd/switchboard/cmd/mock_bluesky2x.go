// Code generated by MockGen. DO NOT EDIT.
// Source: bluesky2x.go
//
// Generated by this command:
//
//	mockgen -source=bluesky2x.go -destination=mock_bluesky2x.go -package=cmd
//

// Package cmd is a generated GoMock package.
package cmd

import (
	context "context"
	reflect "reflect"

	switchboard "github.com/go-zen-chu/switchboard"
	gomock "go.uber.org/mock/gomock"
)

// MockBluesky2XCmdRequirements is a mock of Bluesky2XCmdRequirements interface.
type MockBluesky2XCmdRequirements struct {
	ctrl     *gomock.Controller
	recorder *MockBluesky2XCmdRequirementsMockRecorder
	isgomock struct{}
}

// MockBluesky2XCmdRequirementsMockRecorder is the mock recorder for MockBluesky2XCmdRequirements.
type MockBluesky2XCmdRequirementsMockRecorder struct {
	mock *MockBluesky2XCmdRequirements
}

// NewMockBluesky2XCmdRequirements creates a new mock instance.
func NewMockBluesky2XCmdRequirements(ctrl *gomock.Controller) *MockBluesky2XCmdRequirements {
	mock := &MockBluesky2XCmdRequirements{ctrl: ctrl}
	mock.recorder = &MockBluesky2XCmdRequirementsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBluesky2XCmdRequirements) EXPECT() *MockBluesky2XCmdRequirementsMockRecorder {
	return m.recorder
}

// BlueskyClient mocks base method.
func (m *MockBluesky2XCmdRequirements) BlueskyClient() switchboard.BlueskyClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlueskyClient")
	ret0, _ := ret[0].(switchboard.BlueskyClient)
	return ret0
}

// BlueskyClient indicates an expected call of BlueskyClient.
func (mr *MockBluesky2XCmdRequirementsMockRecorder) BlueskyClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlueskyClient", reflect.TypeOf((*MockBluesky2XCmdRequirements)(nil).BlueskyClient))
}

// Context mocks base method.
func (m *MockBluesky2XCmdRequirements) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockBluesky2XCmdRequirementsMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockBluesky2XCmdRequirements)(nil).Context))
}

// XClient mocks base method.
func (m *MockBluesky2XCmdRequirements) XClient() switchboard.XClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "XClient")
	ret0, _ := ret[0].(switchboard.XClient)
	return ret0
}

// XClient indicates an expected call of XClient.
func (mr *MockBluesky2XCmdRequirementsMockRecorder) XClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "XClient", reflect.TypeOf((*MockBluesky2XCmdRequirements)(nil).XClient))
}
