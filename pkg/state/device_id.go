package state

import "context"

// SetDeviceId sets the device id in the context
func SetDeviceID(ctx context.Context, deviceId string) context.Context {
	return context.WithValue(ctx, KeyDeviceID, deviceId)
}

// GetDeviceId gets the device id from the context
func DeviceID(ctx context.Context) string {
	if deviceId, ok := ctx.Value(KeyDeviceID).(string); ok {
		return deviceId
	}
	return ""
}
