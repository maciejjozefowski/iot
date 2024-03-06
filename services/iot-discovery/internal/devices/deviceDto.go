package devices

type DeviceDto struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Mac   string `json:"mac"`
}

func NewDeviceDto(device Device) DeviceDto {
	return DeviceDto{
		ID:    device.ID,
		Name:  device.Name,
		Brand: device.Brand,
		Model: device.Model,
		Mac:   device.Mac,
	}
}
