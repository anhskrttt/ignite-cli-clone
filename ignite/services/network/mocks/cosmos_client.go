// Code generated by mockery v2.10.6. DO NOT EDIT.

package mocks

import (
	client "github.com/cosmos/cosmos-sdk/client"
	types "github.com/cosmos/cosmos-sdk/types"
	mock "github.com/stretchr/testify/mock"

	cosmosaccount "github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	cosmosclient "github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
)

// CosmosClient is an autogenerated mock type for the CosmosClient type
type CosmosClient struct {
	mock.Mock
}

// Account provides a mock function with given fields: accountName
func (_m *CosmosClient) Account(accountName string) (cosmosaccount.Account, error) {
	ret := _m.Called(accountName)

	var r0 cosmosaccount.Account
	if rf, ok := ret.Get(0).(func(string) cosmosaccount.Account); ok {
		r0 = rf(accountName)
	} else {
		r0 = ret.Get(0).(cosmosaccount.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accountName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Address provides a mock function with given fields: accountName
func (_m *CosmosClient) Address(accountName string) (types.AccAddress, error) {
	ret := _m.Called(accountName)

	var r0 types.AccAddress
	if rf, ok := ret.Get(0).(func(string) types.AccAddress); ok {
		r0 = rf(accountName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.AccAddress)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accountName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BroadcastTx provides a mock function with given fields: accountName, msgs
func (_m *CosmosClient) BroadcastTx(accountName string, msgs ...types.Msg) (cosmosclient.Response, error) {
	_va := make([]interface{}, len(msgs))
	for _i := range msgs {
		_va[_i] = msgs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, accountName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 cosmosclient.Response
	if rf, ok := ret.Get(0).(func(string, ...types.Msg) cosmosclient.Response); ok {
		r0 = rf(accountName, msgs...)
	} else {
		r0 = ret.Get(0).(cosmosclient.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...types.Msg) error); ok {
		r1 = rf(accountName, msgs...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BroadcastTxWithProvision provides a mock function with given fields: accountName, msgs
func (_m *CosmosClient) BroadcastTxWithProvision(accountName string, msgs ...types.Msg) (uint64, func() (cosmosclient.Response, error), error) {
	_va := make([]interface{}, len(msgs))
	for _i := range msgs {
		_va[_i] = msgs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, accountName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(string, ...types.Msg) uint64); ok {
		r0 = rf(accountName, msgs...)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 func() (cosmosclient.Response, error)
	if rf, ok := ret.Get(1).(func(string, ...types.Msg) func() (cosmosclient.Response, error)); ok {
		r1 = rf(accountName, msgs...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func() (cosmosclient.Response, error))
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, ...types.Msg) error); ok {
		r2 = rf(accountName, msgs...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Context provides a mock function with given fields:
func (_m *CosmosClient) Context() client.Context {
	ret := _m.Called()

	var r0 client.Context
	if rf, ok := ret.Get(0).(func() client.Context); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(client.Context)
	}

	return r0
}
