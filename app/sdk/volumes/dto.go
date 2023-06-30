package volumes

type VolumeList struct {
	CreatedAt string                 `json:"createdAt"`
	Driver    string                 `json:"driver"`
	Name      string                 `json:"name"`
	Scope     string                 `json:"scope"`
	Status    map[string]interface{} `json:"status"`
}

type VolumeCreateStruct struct {
	Driver     string            `json:"driver"`
	DriverOpts map[string]string `json:"driverOpts"`
	Labels     map[string]string `json:"labels"`
	Name       string            `json:"name"`
}

type VolumeDeleteStruct struct {
	Name []string `json:"name"`
}
