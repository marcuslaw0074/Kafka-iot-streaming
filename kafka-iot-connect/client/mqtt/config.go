package mqtt

type MqttConf struct {
	Host    string `json:"host,omitempty"`
	Port    int    `json:"port,omitempty"`
	TlsPath string `json:"tlsPath,omitempty"`
}