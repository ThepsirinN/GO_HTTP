package servicesV1

import "context"

func (s *svcV1) GetHelloWorldSvc(ctx context.Context) string {
	return "Hello, World!"
}