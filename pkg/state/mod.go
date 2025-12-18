package state

type contextKeyType string

const (
	KeyLocale       contextKeyType = "locale"
	KeyIP           contextKeyType = "ip"
	KeyDeviceID     contextKeyType = "device_id"
	KeyUser         contextKeyType = "user"
	KeyAccessToken  contextKeyType = "access_token"
	KeyRefreshToken contextKeyType = "refresh_token"
	KeyUserRefresh  contextKeyType = "user_refresh"
	KeyCurrency     contextKeyType = "currency"
	KeyDevice       contextKeyType = "device"
)
