package etcd

type IConfig interface {
	KeyConfigHostPrefix() string
	KeyOnlinePrefix() string
	CaPem() string
	ClientPem() string
	ClientKeyPem() string
	Username() string
	Password() string
}

type Config struct {
	keyOnlinePrefix     string
	keyConfigHostPrefix string
	caPem               string
	clientPem           string
	clientKeyPem        string
	username            string
	password            string
}

func (c Config) KeyConfigHostPrefix() string {
	return c.keyConfigHostPrefix
}

func (c Config) KeyOnlinePrefix() string {
	return c.keyOnlinePrefix
}

func (c Config) CaPem() string {
	return c.caPem
}

func (c Config) ClientPem() string {
	return c.clientPem
}

func (c Config) ClientKeyPem() string {
	return c.clientKeyPem
}

func (c Config) Username() string {
	return c.username
}

func (c Config) Password() string {
	return c.password
}

func NewProdConfig() Config {
	return Config{
		keyOnlinePrefix:     Etcd_key_online_prefix,
		keyConfigHostPrefix: Etcd_key_config_host_prefix,
		caPem:               Etcd_ca_pem,
		clientPem:           Etcd_client_pem,
		clientKeyPem:        Etcd_client_key_pem,
		username:            Etcd_username,
		password:            Etcd_passwd,
	}
}

func NewTestConfig() Config {
	return Config{
		keyOnlinePrefix:     Etcd_key_online_prefix_test,
		keyConfigHostPrefix: Etcd_key_config_host_prefix_test,
		caPem:               Etcd_ca_pem_test,
		clientPem:           Etcd_client_pem_test,
		clientKeyPem:        Etcd_client_key_pem_test,
		username:            Etcd_username_test,
		password:            Etcd_passwd_test,
	}
}

func NewDevConfig() Config {
	return Config{
		keyOnlinePrefix:     Etcd_key_online_prefix_dev,
		keyConfigHostPrefix: Etcd_key_config_host_prefix_dev,
		caPem:               Etcd_ca_pem_dev,
		clientPem:           Etcd_client_pem_dev,
		clientKeyPem:        Etcd_client_key_pem_dev,
		username:            Etcd_username_dev,
		password:            Etcd_passwd_dev,
	}
}
