package definitions

type Credentials struct {
	DeviceName string `json:"deviceName"`
	DeviceKey string `json:"deviceKey"`
	DeviceGroup string `json:"deviceGroup"`
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	KeyID string `json:"keyId"`
}
