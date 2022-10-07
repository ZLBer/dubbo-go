package config

type TLSConfig struct {
	CACertFile    string ` yaml:"ca_cert_file" json:"ca_cert_file" property:"ca_cert_file"`
	TLSCertFile   string `yaml:"tls_cert_file" json:"tls_cert_file" property:"tls_cert_file"`
	TLSKeyFile    string `yaml:"tls_key_file" json:"tls_key_file" property:"tls_key_file"`
	TLSServerName string `yaml:"tls_server_name" json:"tls_server_name" property:"tls_server_name"`
}
