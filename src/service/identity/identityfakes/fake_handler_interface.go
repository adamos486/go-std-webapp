// Code generated by counterfeiter. DO NOT EDIT.
package identityfakes

import (
	"net/http"
	"service/identity"
	"sync"
)

type FakeHandlerInterface struct {
	HandlerStub        func(w http.ResponseWriter, req *http.Request)
	handlerMutex       sync.RWMutex
	handlerArgsForCall []struct {
		w   http.ResponseWriter
		req *http.Request
	}
	CreateIdentityStub        func(w http.ResponseWriter, req *http.Request)
	createIdentityMutex       sync.RWMutex
	createIdentityArgsForCall []struct {
		w   http.ResponseWriter
		req *http.Request
	}
	AuthIdentityStub        func(w http.ResponseWriter, req *http.Request)
	authIdentityMutex       sync.RWMutex
	authIdentityArgsForCall []struct {
		w   http.ResponseWriter
		req *http.Request
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHandlerInterface) Handler(w http.ResponseWriter, req *http.Request) {
	fake.handlerMutex.Lock()
	fake.handlerArgsForCall = append(fake.handlerArgsForCall, struct {
		w   http.ResponseWriter
		req *http.Request
	}{w, req})
	fake.recordInvocation("Handler", []interface{}{w, req})
	fake.handlerMutex.Unlock()
	if fake.HandlerStub != nil {
		fake.HandlerStub(w, req)
	}
}

func (fake *FakeHandlerInterface) HandlerCallCount() int {
	fake.handlerMutex.RLock()
	defer fake.handlerMutex.RUnlock()
	return len(fake.handlerArgsForCall)
}

func (fake *FakeHandlerInterface) HandlerArgsForCall(i int) (http.ResponseWriter, *http.Request) {
	fake.handlerMutex.RLock()
	defer fake.handlerMutex.RUnlock()
	return fake.handlerArgsForCall[i].w, fake.handlerArgsForCall[i].req
}

func (fake *FakeHandlerInterface) CreateIdentity(w http.ResponseWriter, req *http.Request) {
	fake.createIdentityMutex.Lock()
	fake.createIdentityArgsForCall = append(fake.createIdentityArgsForCall, struct {
		w   http.ResponseWriter
		req *http.Request
	}{w, req})
	fake.recordInvocation("CreateIdentity", []interface{}{w, req})
	fake.createIdentityMutex.Unlock()
	if fake.CreateIdentityStub != nil {
		fake.CreateIdentityStub(w, req)
	}
}

func (fake *FakeHandlerInterface) CreateIdentityCallCount() int {
	fake.createIdentityMutex.RLock()
	defer fake.createIdentityMutex.RUnlock()
	return len(fake.createIdentityArgsForCall)
}

func (fake *FakeHandlerInterface) CreateIdentityArgsForCall(i int) (http.ResponseWriter, *http.Request) {
	fake.createIdentityMutex.RLock()
	defer fake.createIdentityMutex.RUnlock()
	return fake.createIdentityArgsForCall[i].w, fake.createIdentityArgsForCall[i].req
}

func (fake *FakeHandlerInterface) AuthIdentity(w http.ResponseWriter, req *http.Request) {
	fake.authIdentityMutex.Lock()
	fake.authIdentityArgsForCall = append(fake.authIdentityArgsForCall, struct {
		w   http.ResponseWriter
		req *http.Request
	}{w, req})
	fake.recordInvocation("AuthIdentity", []interface{}{w, req})
	fake.authIdentityMutex.Unlock()
	if fake.AuthIdentityStub != nil {
		fake.AuthIdentityStub(w, req)
	}
}

func (fake *FakeHandlerInterface) AuthIdentityCallCount() int {
	fake.authIdentityMutex.RLock()
	defer fake.authIdentityMutex.RUnlock()
	return len(fake.authIdentityArgsForCall)
}

func (fake *FakeHandlerInterface) AuthIdentityArgsForCall(i int) (http.ResponseWriter, *http.Request) {
	fake.authIdentityMutex.RLock()
	defer fake.authIdentityMutex.RUnlock()
	return fake.authIdentityArgsForCall[i].w, fake.authIdentityArgsForCall[i].req
}

func (fake *FakeHandlerInterface) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.handlerMutex.RLock()
	defer fake.handlerMutex.RUnlock()
	fake.createIdentityMutex.RLock()
	defer fake.createIdentityMutex.RUnlock()
	fake.authIdentityMutex.RLock()
	defer fake.authIdentityMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeHandlerInterface) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ identity.HandlerInterface = new(FakeHandlerInterface)
