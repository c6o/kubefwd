// Code generated by MockGen. DO NOT EDIT.
// Source: fwdport.go

// Package fwdport is a generated GoMock package.
package fwdport

import (
	gomock "github.com/golang/mock/gomock"
	io "io"
	v1 "k8s.io/api/core/v1"
	httpstream "k8s.io/apimachinery/pkg/util/httpstream"
	rest "k8s.io/client-go/rest"
	portforward "k8s.io/client-go/tools/portforward"
	spdy "k8s.io/client-go/transport/spdy"
	http "net/http"
	reflect "reflect"
)

// MockServiceFWD is a mock of ServiceFWD interface
type MockServiceFWD struct {
	ctrl     *gomock.Controller
	recorder *MockServiceFWDMockRecorder
}

// MockServiceFWDMockRecorder is the mock recorder for MockServiceFWD
type MockServiceFWDMockRecorder struct {
	mock *MockServiceFWD
}

// NewMockServiceFWD creates a new mock instance
func NewMockServiceFWD(ctrl *gomock.Controller) *MockServiceFWD {
	mock := &MockServiceFWD{ctrl: ctrl}
	mock.recorder = &MockServiceFWDMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockServiceFWD) EXPECT() *MockServiceFWDMockRecorder {
	return m.recorder
}

// String mocks base method
func (m *MockServiceFWD) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockServiceFWDMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockServiceFWD)(nil).String))
}

// SyncPodForwards mocks base method
func (m *MockServiceFWD) SyncPodForwards(arg0 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SyncPodForwards", arg0)
}

// SyncPodForwards indicates an expected call of SyncPodForwards
func (mr *MockServiceFWDMockRecorder) SyncPodForwards(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncPodForwards", reflect.TypeOf((*MockServiceFWD)(nil).SyncPodForwards), arg0)
}

// ListServicePodNames mocks base method
func (m *MockServiceFWD) ListServicePodNames() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListServicePodNames")
	ret0, _ := ret[0].([]string)
	return ret0
}

// ListServicePodNames indicates an expected call of ListServicePodNames
func (mr *MockServiceFWDMockRecorder) ListServicePodNames() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListServicePodNames", reflect.TypeOf((*MockServiceFWD)(nil).ListServicePodNames))
}

// AddServicePod mocks base method
func (m *MockServiceFWD) AddServicePod(pfo *PortForwardOpts) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddServicePod", pfo)
}

// AddServicePod indicates an expected call of AddServicePod
func (mr *MockServiceFWDMockRecorder) AddServicePod(pfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddServicePod", reflect.TypeOf((*MockServiceFWD)(nil).AddServicePod), pfo)
}

// GetServicePodPortForwards mocks base method
func (m *MockServiceFWD) GetServicePodPortForwards(servicePodName string) []*PortForwardOpts {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServicePodPortForwards", servicePodName)
	ret0, _ := ret[0].([]*PortForwardOpts)
	return ret0
}

// GetServicePodPortForwards indicates an expected call of GetServicePodPortForwards
func (mr *MockServiceFWDMockRecorder) GetServicePodPortForwards(servicePodName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServicePodPortForwards", reflect.TypeOf((*MockServiceFWD)(nil).GetServicePodPortForwards), servicePodName)
}

// RemoveServicePod mocks base method
func (m *MockServiceFWD) RemoveServicePod(servicePodName string, stop bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveServicePod", servicePodName, stop)
}

// RemoveServicePod indicates an expected call of RemoveServicePod
func (mr *MockServiceFWDMockRecorder) RemoveServicePod(servicePodName, stop interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveServicePod", reflect.TypeOf((*MockServiceFWD)(nil).RemoveServicePod), servicePodName, stop)
}

// RemoveServicePodByPort mocks base method
func (m *MockServiceFWD) RemoveServicePodByPort(servicePodName, podPort string, stop bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveServicePodByPort", servicePodName, podPort, stop)
}

// RemoveServicePodByPort indicates an expected call of RemoveServicePodByPort
func (mr *MockServiceFWDMockRecorder) RemoveServicePodByPort(servicePodName, podPort, stop interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveServicePodByPort", reflect.TypeOf((*MockServiceFWD)(nil).RemoveServicePodByPort), servicePodName, podPort, stop)
}

// MockPortForwardHelper is a mock of PortForwardHelper interface
type MockPortForwardHelper struct {
	ctrl     *gomock.Controller
	recorder *MockPortForwardHelperMockRecorder
}

// MockPortForwardHelperMockRecorder is the mock recorder for MockPortForwardHelper
type MockPortForwardHelperMockRecorder struct {
	mock *MockPortForwardHelper
}

// NewMockPortForwardHelper creates a new mock instance
func NewMockPortForwardHelper(ctrl *gomock.Controller) *MockPortForwardHelper {
	mock := &MockPortForwardHelper{ctrl: ctrl}
	mock.recorder = &MockPortForwardHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPortForwardHelper) EXPECT() *MockPortForwardHelperMockRecorder {
	return m.recorder
}

// GetPortForwardRequest mocks base method
func (m *MockPortForwardHelper) GetPortForwardRequest(pfo *PortForwardOpts) *rest.Request {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPortForwardRequest", pfo)
	ret0, _ := ret[0].(*rest.Request)
	return ret0
}

// GetPortForwardRequest indicates an expected call of GetPortForwardRequest
func (mr *MockPortForwardHelperMockRecorder) GetPortForwardRequest(pfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPortForwardRequest", reflect.TypeOf((*MockPortForwardHelper)(nil).GetPortForwardRequest), pfo)
}

// NewOnAddresses mocks base method
func (m *MockPortForwardHelper) NewOnAddresses(dialer httpstream.Dialer, addresses, ports []string, stopChan <-chan struct{}, readyChan chan struct{}, out, errOut io.Writer) (*portforward.PortForwarder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewOnAddresses", dialer, addresses, ports, stopChan, readyChan, out, errOut)
	ret0, _ := ret[0].(*portforward.PortForwarder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewOnAddresses indicates an expected call of NewOnAddresses
func (mr *MockPortForwardHelperMockRecorder) NewOnAddresses(dialer, addresses, ports, stopChan, readyChan, out, errOut interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewOnAddresses", reflect.TypeOf((*MockPortForwardHelper)(nil).NewOnAddresses), dialer, addresses, ports, stopChan, readyChan, out, errOut)
}

// ForwardPorts mocks base method
func (m *MockPortForwardHelper) ForwardPorts(forwarder *portforward.PortForwarder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForwardPorts", forwarder)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForwardPorts indicates an expected call of ForwardPorts
func (mr *MockPortForwardHelperMockRecorder) ForwardPorts(forwarder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForwardPorts", reflect.TypeOf((*MockPortForwardHelper)(nil).ForwardPorts), forwarder)
}

// RoundTripperFor mocks base method
func (m *MockPortForwardHelper) RoundTripperFor(config *rest.Config) (http.RoundTripper, spdy.Upgrader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RoundTripperFor", config)
	ret0, _ := ret[0].(http.RoundTripper)
	ret1, _ := ret[1].(spdy.Upgrader)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RoundTripperFor indicates an expected call of RoundTripperFor
func (mr *MockPortForwardHelperMockRecorder) RoundTripperFor(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RoundTripperFor", reflect.TypeOf((*MockPortForwardHelper)(nil).RoundTripperFor), config)
}

// NewDialer mocks base method
func (m *MockPortForwardHelper) NewDialer(upgrader spdy.Upgrader, client *http.Client, method string, pfRequest *rest.Request) httpstream.Dialer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewDialer", upgrader, client, method, pfRequest)
	ret0, _ := ret[0].(httpstream.Dialer)
	return ret0
}

// NewDialer indicates an expected call of NewDialer
func (mr *MockPortForwardHelperMockRecorder) NewDialer(upgrader, client, method, pfRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewDialer", reflect.TypeOf((*MockPortForwardHelper)(nil).NewDialer), upgrader, client, method, pfRequest)
}

// MockHostsOperator is a mock of HostsOperator interface
type MockHostsOperator struct {
	ctrl     *gomock.Controller
	recorder *MockHostsOperatorMockRecorder
}

// MockHostsOperatorMockRecorder is the mock recorder for MockHostsOperator
type MockHostsOperatorMockRecorder struct {
	mock *MockHostsOperator
}

// NewMockHostsOperator creates a new mock instance
func NewMockHostsOperator(ctrl *gomock.Controller) *MockHostsOperator {
	mock := &MockHostsOperator{ctrl: ctrl}
	mock.recorder = &MockHostsOperatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHostsOperator) EXPECT() *MockHostsOperatorMockRecorder {
	return m.recorder
}

// AddHosts mocks base method
func (m *MockHostsOperator) AddHosts() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddHosts")
}

// AddHosts indicates an expected call of AddHosts
func (mr *MockHostsOperatorMockRecorder) AddHosts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHosts", reflect.TypeOf((*MockHostsOperator)(nil).AddHosts))
}

// RemoveHosts mocks base method
func (m *MockHostsOperator) RemoveHosts() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveHosts")
}

// RemoveHosts indicates an expected call of RemoveHosts
func (mr *MockHostsOperatorMockRecorder) RemoveHosts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveHosts", reflect.TypeOf((*MockHostsOperator)(nil).RemoveHosts))
}

// RemoveInterfaceAlias mocks base method
func (m *MockHostsOperator) RemoveInterfaceAlias() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveInterfaceAlias")
}

// RemoveInterfaceAlias indicates an expected call of RemoveInterfaceAlias
func (mr *MockHostsOperatorMockRecorder) RemoveInterfaceAlias() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveInterfaceAlias", reflect.TypeOf((*MockHostsOperator)(nil).RemoveInterfaceAlias))
}

// MockPodStateWaiter is a mock of PodStateWaiter interface
type MockPodStateWaiter struct {
	ctrl     *gomock.Controller
	recorder *MockPodStateWaiterMockRecorder
}

// MockPodStateWaiterMockRecorder is the mock recorder for MockPodStateWaiter
type MockPodStateWaiterMockRecorder struct {
	mock *MockPodStateWaiter
}

// NewMockPodStateWaiter creates a new mock instance
func NewMockPodStateWaiter(ctrl *gomock.Controller) *MockPodStateWaiter {
	mock := &MockPodStateWaiter{ctrl: ctrl}
	mock.recorder = &MockPodStateWaiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPodStateWaiter) EXPECT() *MockPodStateWaiterMockRecorder {
	return m.recorder
}

// WaitUntilPodRunning mocks base method
func (m *MockPodStateWaiter) WaitUntilPodRunning(stopChannel <-chan struct{}) (*v1.Pod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitUntilPodRunning", stopChannel)
	ret0, _ := ret[0].(*v1.Pod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WaitUntilPodRunning indicates an expected call of WaitUntilPodRunning
func (mr *MockPodStateWaiterMockRecorder) WaitUntilPodRunning(stopChannel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitUntilPodRunning", reflect.TypeOf((*MockPodStateWaiter)(nil).WaitUntilPodRunning), stopChannel)
}

// ListenUntilPodDeleted mocks base method
func (m *MockPodStateWaiter) ListenUntilPodDeleted(stopChannel <-chan struct{}, pod *v1.Pod) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ListenUntilPodDeleted", stopChannel, pod)
	return
}

// ListenUntilPodDeleted indicates an expected call of ListenUntilPodDeleted
func (mr *MockPodStateWaiterMockRecorder) ListenUntilPodDeleted(stopChannel interface{}, pod interface {}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenUntilPodDeleted", reflect.TypeOf((*MockPodStateWaiter)(nil).ListenUntilPodDeleted), stopChannel, pod)
}

