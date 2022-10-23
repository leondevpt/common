package consul


import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type consul struct {
	client *api.Client
}

type RegisterOption struct {
	ServiceID string     // 服务唯一ID
	ServiceName string   // 服务名称
	ServiceIP  string
	ServicePort int
	Tags []string      // 为服务打标签
	HealthCheck bool
	HealthCheckTimeout string
	HealthCheckInterval string
	HealthCheckGrpcEndPoint string
	DeregisterCriticalServiceAfter string
}

// NewConsul 连接至consul服务返回一个consul对象
func NewConsul(addr string) (*consul, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &consul{c}, nil
}


// RegisterService 将gRPC服务注册到consul
func (c *consul) RegisterService(opt RegisterOption) error {
	if opt.ServiceID == "" {
		opt.ServiceID = fmt.Sprintf("%s-%s-%d", opt.ServiceName, opt.ServiceIP, opt.ServicePort)
	}
	if opt.HealthCheckGrpcEndPoint == "" {
		opt.HealthCheckGrpcEndPoint = fmt.Sprintf("%s:%d", opt.ServiceIP, opt.ServicePort)
	}
	if opt.HealthCheckTimeout == "" {
		opt.HealthCheckTimeout = "5s"
	}
	if opt.HealthCheckInterval == "" {
		opt.HealthCheckInterval =   "10s"
	}
	if opt.DeregisterCriticalServiceAfter == "" {
		opt.DeregisterCriticalServiceAfter = "15s"
	}

	srv := &api.AgentServiceRegistration{
		ID:      opt.ServiceID, // 服务唯一ID
		Name:    opt.ServiceName,       // 服务名称
		Tags:    opt.Tags,             // 为服务打标签
		Address: opt.ServiceIP,
		Port:    opt.ServicePort,
	}
	// 健康检查
	if opt.HealthCheck {
		check := &api.AgentServiceCheck{
			GRPC:                           opt.HealthCheckGrpcEndPoint, // 这里一定是外部可以访问的地址
			Timeout:                        opt.HealthCheckTimeout,
			Interval:                       opt.HealthCheckInterval,
			DeregisterCriticalServiceAfter: opt.DeregisterCriticalServiceAfter,
		}
		srv.Check = check
	}

	return c.client.Agent().ServiceRegister(srv)
}
// Deregister 注销服务
func (c *consul) Deregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}