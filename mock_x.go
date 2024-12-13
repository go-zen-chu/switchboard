// Code generated by MockGen. DO NOT EDIT.
// Source: x.go
//
// Generated by this command:
//
//	mockgen -source=x.go -destination=mock_x.go -package=switchboard
//

// Package switchboard is a generated GoMock package.
package switchboard

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockXClient is a mock of XClient interface.
type MockXClient struct {
	ctrl     *gomock.Controller
	recorder *MockXClientMockRecorder
	isgomock struct{}
}

// MockXClientMockRecorder is the mock recorder for MockXClient.
type MockXClientMockRecorder struct {
	mock *MockXClient
}

// NewMockXClient creates a new mock instance.
func NewMockXClient(ctrl *gomock.Controller) *MockXClient {
	mock := &MockXClient{ctrl: ctrl}
	mock.recorder = &MockXClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockXClient) EXPECT() *MockXClientMockRecorder {
	return m.recorder
}

// Post mocks base method.
func (m *MockXClient) Post(ctx context.Context, content string) (*XPost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", ctx, content)
	ret0, _ := ret[0].(*XPost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Post indicates an expected call of Post.
func (mr *MockXClientMockRecorder) Post(ctx, content any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockXClient)(nil).Post), ctx, content)
}
