package config

import (
	"awsx-api/log"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/model"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sync"
)

// Global configuration for the application.
var configuration Config
var rwMutex sync.RWMutex

var awsClientCache = make(map[string]*model.Auth)

// Server configuration
type Server struct {
	Address                    string `yaml:"address,omitempty"`
	AuditLog                   bool   `yaml:"audit_log,omitempty"` // When true, allows additional audit logging on Write operations
	CORSAllowAll               bool   `yaml:"cors_allow_all,omitempty"`
	GzipEnabled                bool   `yaml:"gzip_enabled,omitempty"`
	Port                       int    `yaml:"port,omitempty"`
	StaticContentRootDirectory string `yaml:"static_content_root_directory,omitempty"`
	WebFQDN                    string `yaml:"web_fqdn,omitempty"`
	WebPort                    string `yaml:"web_port,omitempty"`
	WebRoot                    string `yaml:"web_root,omitempty"`
	WebHistoryMode             string `yaml:"web_history_mode,omitempty"`
	WebSchema                  string `yaml:"web_schema,omitempty"`
	WhiteListUrls              string `yaml:"white_list_urls,omitempty"`
}

// Vault configuration
type Vault struct {
	Url string `yaml:"url,omitempty"`
}

type CloudElement struct {
	Url        string `yaml:"url,omitempty"`
	ElementKey string `yaml:"element_key,omitempty"`
}

type Config struct {
	Server       Server       `yaml:",omitempty"`
	Vault        Vault        `yaml:",omitempty"`
	CloudElement CloudElement `yaml:",omitempty"`
}

func LoadFromFile(filename string) (conf *Config, err error) {
	log.Debugf("Reading YAML config from [%s]", filename)
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file [%v]. error=%v", filename, err)
	}

	conf, err = Unmarshal(string(fileContent))
	if err != nil {
		return
	}

	return
}

func Unmarshal(yamlString string) (conf *Config, err error) {
	conf = NewConfig()
	err = yaml.Unmarshal([]byte(yamlString), &conf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml data. error=%v", err)
	}
	return
}

func NewConfig() (c *Config) {
	c = &Config{
		Server: Server{
			AuditLog:                   true,
			GzipEnabled:                true,
			Port:                       7000,
			StaticContentRootDirectory: "/opt/awsx-api/console",
			WebFQDN:                    "",
			WebRoot:                    "/",
			WebHistoryMode:             "browser",
			WebSchema:                  "",
		},
		Vault:        Vault{},
		CloudElement: CloudElement{},
	}

	return
}

// Get the global Config
func Get() (conf *Config) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	copy := configuration
	return &copy
}

// Set the global Config
// This function should not be called outside of main or tests.
// If possible keep config unmutated and use globals and/or appstate package for mutable states to avoid concurrent writes risk.
func Set(conf *Config) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	configuration = *conf
}

// Get the aws client based on element type
func GetAwsClient(elementType string) (auth *model.Auth) {
	return awsClientCache[elementType]
}

func SetAwsClient(elementType string, auth *model.Auth) {
	awsClientCache[elementType] = auth
}
