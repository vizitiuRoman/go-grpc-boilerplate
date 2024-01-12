package v1

type Config struct {
	Host     string `yaml:"host"`
	GrpcPort int    `yaml:"grpc_port"`
}
