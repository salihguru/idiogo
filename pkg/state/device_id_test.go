package state

import (
	"context"
	"testing"
)

func TestSetDeviceID(t *testing.T) {
	ctx := context.Background()
	deviceId := "test-device-id"

	ctx = SetDeviceID(ctx, deviceId)

	got := ctx.Value(KeyDeviceID)
	if got == nil {
		t.Errorf("SetDeviceId() did not set deviceId in context")
	}

	if got != deviceId {
		t.Errorf("SetDeviceId() = %v, want %v", got, deviceId)
	}
}

func TestDeviceID(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want string
	}{
		{
			name: "DeviceIdExists",
			ctx:  context.WithValue(context.Background(), KeyDeviceID, "device123"),
			want: "device123",
		},
		{
			name: "DeviceIdDoesNotExist",
			ctx:  context.Background(),
			want: "",
		},
		{
			name: "WrongTypeInContext",
			ctx:  context.WithValue(context.Background(), KeyDeviceID, 123),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeviceID(tt.ctx)
			if got != tt.want {
				t.Errorf("DeviceID() = %v, want %v", got, tt.want)
			}
		})
	}
}
