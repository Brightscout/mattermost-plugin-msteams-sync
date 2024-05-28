// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	time "time"

	prometheus "github.com/prometheus/client_golang/prometheus"
	mock "github.com/stretchr/testify/mock"
)

// Metrics is an autogenerated mock type for the Metrics type
type Metrics struct {
	mock.Mock
}

// DecrementActiveWorkers provides a mock function with given fields: worker
func (_m *Metrics) DecrementActiveWorkers(worker string) {
	_m.Called(worker)
}

// DecrementChangeEventQueueLength provides a mock function with given fields: changeType
func (_m *Metrics) DecrementChangeEventQueueLength(changeType string) {
	_m.Called(changeType)
}

// GetRegistry provides a mock function with given fields:
func (_m *Metrics) GetRegistry() *prometheus.Registry {
	ret := _m.Called()

	var r0 *prometheus.Registry
	if rf, ok := ret.Get(0).(func() *prometheus.Registry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*prometheus.Registry)
		}
	}

	return r0
}

// IncrementActiveWorkers provides a mock function with given fields: worker
func (_m *Metrics) IncrementActiveWorkers(worker string) {
	_m.Called(worker)
}

// IncrementChangeEventQueueLength provides a mock function with given fields: changeType
func (_m *Metrics) IncrementChangeEventQueueLength(changeType string) {
	_m.Called(changeType)
}

// IncrementHTTPErrors provides a mock function with given fields:
func (_m *Metrics) IncrementHTTPErrors() {
	_m.Called()
}

// IncrementHTTPRequests provides a mock function with given fields:
func (_m *Metrics) IncrementHTTPRequests() {
	_m.Called()
}

// ObserveAPIEndpointDuration provides a mock function with given fields: handler, method, statusCode, elapsed
func (_m *Metrics) ObserveAPIEndpointDuration(handler string, method string, statusCode string, elapsed float64) {
	_m.Called(handler, method, statusCode, elapsed)
}

// ObserveActiveUsersReceiving provides a mock function with given fields: count
func (_m *Metrics) ObserveActiveUsersReceiving(count int64) {
	_m.Called(count)
}

// ObserveActiveUsersSending provides a mock function with given fields: count
func (_m *Metrics) ObserveActiveUsersSending(count int64) {
	_m.Called(count)
}

// ObserveChangeEvent provides a mock function with given fields: changeType, discardedReason
func (_m *Metrics) ObserveChangeEvent(changeType string, discardedReason string) {
	_m.Called(changeType, discardedReason)
}

// ObserveChangeEventQueueCapacity provides a mock function with given fields: count
func (_m *Metrics) ObserveChangeEventQueueCapacity(count int64) {
	_m.Called(count)
}

// ObserveChangeEventQueueRejected provides a mock function with given fields:
func (_m *Metrics) ObserveChangeEventQueueRejected() {
	_m.Called()
}

// ObserveClientSecretEndDateTime provides a mock function with given fields: expireDate
func (_m *Metrics) ObserveClientSecretEndDateTime(expireDate time.Time) {
	_m.Called(expireDate)
}

// ObserveConnectedUsers provides a mock function with given fields: count
func (_m *Metrics) ObserveConnectedUsers(count int64) {
	_m.Called(count)
}

// ObserveFile provides a mock function with given fields: action, source, discardedReason, isDirectOrGroupMessage
func (_m *Metrics) ObserveFile(action string, source string, discardedReason string, isDirectOrGroupMessage bool) {
	_m.Called(action, source, discardedReason, isDirectOrGroupMessage)
}

// ObserveFiles provides a mock function with given fields: action, source, discardedReason, isDirectOrGroupMessage, count
func (_m *Metrics) ObserveFiles(action string, source string, discardedReason string, isDirectOrGroupMessage bool, count int64) {
	_m.Called(action, source, discardedReason, isDirectOrGroupMessage, count)
}

// ObserveGoroutineFailure provides a mock function with given fields:
func (_m *Metrics) ObserveGoroutineFailure() {
	_m.Called()
}

// ObserveLifecycleEvent provides a mock function with given fields: lifecycleEventType, discardedReason
func (_m *Metrics) ObserveLifecycleEvent(lifecycleEventType string, discardedReason string) {
	_m.Called(lifecycleEventType, discardedReason)
}

// ObserveLinkedChannels provides a mock function with given fields: count
func (_m *Metrics) ObserveLinkedChannels(count int64) {
	_m.Called(count)
}

// ObserveMSGraphClientMethodDuration provides a mock function with given fields: method, success, statusCode, elapsed
func (_m *Metrics) ObserveMSGraphClientMethodDuration(method string, success string, statusCode string, elapsed float64) {
	_m.Called(method, success, statusCode, elapsed)
}

// ObserveMessage provides a mock function with given fields: action, source, isDirectOrGroupMessage
func (_m *Metrics) ObserveMessage(action string, source string, isDirectOrGroupMessage bool) {
	_m.Called(action, source, isDirectOrGroupMessage)
}

// ObserveMessageDelay provides a mock function with given fields: action, source, isDirectOrGroupMessage, delay
func (_m *Metrics) ObserveMessageDelay(action string, source string, isDirectOrGroupMessage bool, delay time.Duration) {
	_m.Called(action, source, isDirectOrGroupMessage, delay)
}

// ObserveNotification provides a mock function with given fields: isGroupChat, hasAttachments
func (_m *Metrics) ObserveNotification(isGroupChat bool, hasAttachments bool) {
	_m.Called(isGroupChat, hasAttachments)
}

// ObserveOAuthTokenInvalidated provides a mock function with given fields:
func (_m *Metrics) ObserveOAuthTokenInvalidated() {
	_m.Called()
}

// ObserveReaction provides a mock function with given fields: action, source, isDirectOrGroupMessage
func (_m *Metrics) ObserveReaction(action string, source string, isDirectOrGroupMessage bool) {
	_m.Called(action, source, isDirectOrGroupMessage)
}

// ObserveStoreMethodDuration provides a mock function with given fields: method, success, elapsed
func (_m *Metrics) ObserveStoreMethodDuration(method string, success string, elapsed float64) {
	_m.Called(method, success, elapsed)
}

// ObserveSubscription provides a mock function with given fields: action
func (_m *Metrics) ObserveSubscription(action string) {
	_m.Called(action)
}

// ObserveSyncMsgFileDelay provides a mock function with given fields: action, delayMillis
func (_m *Metrics) ObserveSyncMsgFileDelay(action string, delayMillis int64) {
	_m.Called(action, delayMillis)
}

// ObserveSyncMsgPostDelay provides a mock function with given fields: action, delayMillis
func (_m *Metrics) ObserveSyncMsgPostDelay(action string, delayMillis int64) {
	_m.Called(action, delayMillis)
}

// ObserveSyncMsgReactionDelay provides a mock function with given fields: action, delayMillis
func (_m *Metrics) ObserveSyncMsgReactionDelay(action string, delayMillis int64) {
	_m.Called(action, delayMillis)
}

// ObserveSyntheticUsers provides a mock function with given fields: count
func (_m *Metrics) ObserveSyntheticUsers(count int64) {
	_m.Called(count)
}

// ObserveUpstreamUsers provides a mock function with given fields: count
func (_m *Metrics) ObserveUpstreamUsers(count int64) {
	_m.Called(count)
}

// ObserveWhitelistLimit provides a mock function with given fields: limit
func (_m *Metrics) ObserveWhitelistLimit(limit int) {
	_m.Called(limit)
}

// ObserveWorker provides a mock function with given fields: worker
func (_m *Metrics) ObserveWorker(worker string) func() {
	ret := _m.Called(worker)

	var r0 func()
	if rf, ok := ret.Get(0).(func(string) func()); ok {
		r0 = rf(worker)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(func())
		}
	}

	return r0
}

// ObserveWorkerDuration provides a mock function with given fields: worker, elapsed
func (_m *Metrics) ObserveWorkerDuration(worker string, elapsed float64) {
	_m.Called(worker, elapsed)
}

type mockConstructorTestingTNewMetrics interface {
	mock.TestingT
	Cleanup(func())
}

// NewMetrics creates a new instance of Metrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMetrics(t mockConstructorTestingTNewMetrics) *Metrics {
	mock := &Metrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
