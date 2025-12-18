package state

import "context"

func SetIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, KeyIP, ip)
}

func IP(ctx context.Context) string {
	if ip, ok := ctx.Value(KeyIP).(string); ok {
		return ip
	}
	return ""
}
