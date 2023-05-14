package grpc

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
)

type MethodDescriptors interface {
	GetMethodDescriptorByServerInfo(infoFullMethod string) *MethodDescriptor
	JoinServerDescriptor(descriptor *ServerDescriptor)
	Set(methodName string, descriptor *MethodDescriptor)
	Get(methodName string) *MethodDescriptor
}

type MethodDescriptor struct {
	Method any
}

func (m *MethodDescriptor) String() string {
	return m.MethodName()
}

func (m *MethodDescriptor) MethodName() string {
	if m.Method == nil {
		return ""
	}
	methodPointer := reflect.ValueOf(m.Method).Pointer()
	fullName := runtime.FuncForPC(methodPointer).Name()
	methodNameParts := strings.Split(fullName, ".")
	return methodNameParts[len(methodNameParts)-1]
}

type methodDescriptors struct {
	m  map[string]*MethodDescriptor
	mx sync.RWMutex
}

func NewMethodDescriptors() MethodDescriptors {
	return &methodDescriptors{
		m: map[string]*MethodDescriptor{},
	}
}

func (s *methodDescriptors) GetMethodDescriptorByServerInfo(infoFullMethod string) *MethodDescriptor {
	methodNameParts := strings.Split(infoFullMethod, "/")
	methodName := methodNameParts[len(methodNameParts)-1]
	return s.Get(methodName)
}

func (s *methodDescriptors) Set(methodName string, descriptor *MethodDescriptor) {
	s.mx.Lock()
	s.m[methodName] = descriptor
	s.mx.Unlock()
}

func (s *methodDescriptors) Get(methodName string) *MethodDescriptor {
	s.mx.RLock()
	defer s.mx.RUnlock()
	methodDescriptor, ok := s.m[methodName]
	if !ok {
		return nil
	}
	return methodDescriptor
}

func (s *methodDescriptors) JoinServerDescriptor(descriptor *ServerDescriptor) {
	for i := range descriptor.Methods {
		md := &descriptor.Methods[i]
		s.Set(md.MethodName(), md)
	}
}
