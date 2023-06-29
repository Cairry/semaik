package volumes

type Volume struct {
	CreatedAt string
	Driver    string
	Name      string
	Scope     string
	Status    map[string]interface{}
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
