package grpc_test

import (
	"fmt"
	"testing"

	grpc_server "git.casinomodule.org/casino27/saga-svc/internal/common/adapter/server/grpc"

	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/suite"
)

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
	descriptors grpc_server.MethodDescriptors
}

func (s *ServiceSuite) SetupSuite() {
	s.descriptors = grpc_server.NewMethodDescriptors()
}

func (s *ServiceSuite) TestAll() {
	methodName := randomdata.Alphanumeric(5)
	descriptor := grpc_server.MethodDescriptor{}

	s.descriptors.Set(methodName, &descriptor)
	desc := s.descriptors.Get(methodName)
	s.Equal(desc, &descriptor)

	desc = s.descriptors.GetMethodDescriptorByServerInfo(fmt.Sprintf("%s/%s", randomdata.Alphanumeric(5), methodName))
	s.Equal(desc, &descriptor)

	desc = s.descriptors.GetMethodDescriptorByServerInfo(randomdata.Alphanumeric(5))
	s.Nil(desc)
}
