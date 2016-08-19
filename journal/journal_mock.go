package journal

import "github.com/stretchr/testify/mock"

// MockJournal is an autogenerated mock type for the Journal type
type MockJournal struct {
	mock.Mock
}

// Append provides a mock function with given fields: payload
func (_m *MockJournal) Append(payload ...Payload) (Index, error) {
	ret := _m.Called(payload)

	var r0 Index
	if rf, ok := ret.Get(0).(func(...Payload) Index); ok {
		r0 = rf(payload...)
	} else {
		r0 = ret.Get(0).(Index)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...Payload) error); ok {
		r1 = rf(payload...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Notify provides a mock function with given fields: _a0, _a1
func (_m *MockJournal) Notify(_a0 Notifiee, _a1 Index) {
	_m.Called(_a0, _a1)
}

// Restore provides a mock function with given fields: entry
func (_m *MockJournal) Restore(entry ...Entry) (Index, error) {
	ret := _m.Called(entry)

	var r0 Index
	if rf, ok := ret.Get(0).(func(...Entry) Index); ok {
		r0 = rf(entry...)
	} else {
		r0 = ret.Get(0).(Index)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...Entry) error); ok {
		r1 = rf(entry...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEntry is an autogenerated mock type for the Entry type
type MockEntry struct {
	mock.Mock
}

// GetIndex provides a mock function with given fields:
func (_m *MockEntry) GetIndex() Index {
	ret := _m.Called()

	var r0 Index
	if rf, ok := ret.Get(0).(func() Index); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(Index)
	}

	return r0
}

// GetPayload provides a mock function with given fields:
func (_m *MockEntry) GetPayload() Payload {
	ret := _m.Called()

	var r0 Payload
	if rf, ok := ret.Get(0).(func() Payload); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Payload)
		}
	}

	return r0
}

// MockNotifiee is an autogenerated mock type for the Notifiee type
type MockNotifiee struct {
	mock.Mock
}

// AppendedEntry provides a mock function with given fields: _a0
func (_m *MockNotifiee) AppendedEntry(_a0 Entry) {
	_m.Called(_a0)
}