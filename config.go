package guerrilla

// AppConfig is the holder of the configuration of the app
type AppConfig struct {
	BackendName   string          `json:"backend_name"`
	BackendConfig BackendConfig   `json:"backend_config,omitempty"`
	Servers       []*ServerConfig `json:"servers"`
	AllowedHosts  []string        `json:"allowed_hosts"`
	AnalyticsConfig
}

// ServerConfig specifies config options for a single server
type ServerConfig struct {
	IsEnabled       bool     `json:"is_enabled"`
	Hostname        string   `json:"host_name"`
	AllowedHosts    []string `json:"allowed_hosts"`
	MaxSize         int64    `json:"max_size"`
	PrivateKeyFile  string   `json:"private_key_file"`
	PublicKeyFile   string   `json:"public_key_file"`
	Timeout         int      `json:"timeout"`
	ListenInterface string   `json:"listen_interface"`
	AdvertiseTLS    bool     `json:"advertise_tls,omitempty"`
	RequireTLS      bool     `json:"require_tls,omitempty"`
	MaxClients      int      `json:"max_clients"`
}

// Specifies config options for collecting data about the app while
// it is running, and displying that data in a web dashboard.
type AnalyticsConfig struct {
	// Whether the app should collect data about its performance
	Enabled bool
	// Credentials for accessing the web dashboard
	WebUsername     string
	WebPassword     string
	ListenInterface string
	PrivateKeyFile  string
	PublicKeyFile   string
}

type BackendConfig struct {
	LogReceivedMail bool
}
