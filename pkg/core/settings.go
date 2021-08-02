package core

// Settings ...
type Settings struct {
	Port    int32          `yaml:"port"`
	JWT     *JWT           `yaml:"jwt"`
	MongoDB *MongoDBConfig `yaml:"mongodb"`
}

// JWT ...
type JWT struct {
	ExpirationInMinutes int64  `yaml:"expiration_in_minutes"`
	Secret              string `yaml:"secret"`
}

// MongoDBConfig ...
type MongoDBConfig struct {
	Database         string `yaml:"database"`
	ConnectionString string `yaml:"connection_string"`
}
