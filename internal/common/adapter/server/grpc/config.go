package grpc

import "time"

type Config struct {
	Host                    string `yaml:"host"`
	GrpcPort                int    `yaml:"grpc_port"`
	HttpPort                int    `yaml:"http_port"`
	UseTLS                  bool   `yaml:"use_tls"`
	MaxSendMessageLength    int    `yaml:"max_send_message_length"`
	MaxReceiveMessageLength int    `yaml:"max_receive_message_length"`

	ReadDeadline  time.Duration `yaml:"read_deadline"`
	WriteDeadline time.Duration `yaml:"write_deadline"`
}
