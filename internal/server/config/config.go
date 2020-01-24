package config

type ServerConfig struct {
	ControlBindAddr,
	DataBindAddr,
	UserBindAddr,
	KeyFile,
	CertFile,
	ClientCaFile string

	ClientAuth bool
}
