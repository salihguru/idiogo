package state

import "context"

type AgentDevice struct {
	Name string `json:"name"`
	Type string `json:"type"`
	OS   string `json:"os"`
	IP   string `json:"ip"`
}

func Device(ctx context.Context) *AgentDevice {
	if device, ok := ctx.Value(KeyDevice).(*AgentDevice); ok {
		return device
	}
	return nil
}

func SetDevice(ctx context.Context, device *AgentDevice) context.Context {
	return context.WithValue(ctx, KeyDevice, device)
}
