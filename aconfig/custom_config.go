package aconfig

const (
	defaultLogLevel     = "info"
	defaultLogDirectory = "/data/log"

	defaultMaxOpenConnections        = 100
	defaultMaxIdleConnections        = 5
	defaultConnectionMaxLifeSeconds  = 3600 // an hour
	defaultConnectionMaxIdleSeconds  = 300  // 5 minutes
	defaultSlowThresholdMilliseconds = 500  // 0.5 second
)

type Common struct {
	Log       Log       `json:"log,omitempty"`
	Database  Database  `json:"database,omitempty"`
	Encryptor Encryptor `json:"encryptor,omitempty"`
}

func (c *Common) Complete() {
	c.Log.complete()
	c.Database.complete()
}

type Log struct {
	Level     string `json:"level,omitempty"`
	Directory string `json:"directory,omitempty"`
	Format    string `json:"format,omitempty"`
}

func (l *Log) complete() {
	if l.Level == "" {
		l.Level = defaultLogLevel
	}
	if l.Directory == "" {
		l.Directory = defaultLogDirectory
	}
}

type Database struct {
	MaxOpenConnections        int   `json:"max_open_connections,omitempty"`
	MaxIdleConnections        int   `json:"max_idle_connections,omitempty"`
	ConnectionMaxLifeSeconds  int64 `json:"connection_max_life_seconds,omitempty"`
	ConnectionMaxIdleSeconds  int64 `json:"connection_max_idle_seconds,omitempty"`
	SlowThresholdMilliseconds int64 `json:"slow_threshold_milliseconds,omitempty"`
}

func (db *Database) complete() {
	if db.MaxOpenConnections == 0 {
		db.MaxOpenConnections = defaultMaxOpenConnections
	}
	if db.MaxIdleConnections == 0 {
		db.MaxIdleConnections = defaultMaxIdleConnections
	}
	if db.ConnectionMaxLifeSeconds == 0 {
		db.ConnectionMaxLifeSeconds = defaultConnectionMaxLifeSeconds
	}
	if db.ConnectionMaxIdleSeconds == 0 {
		db.ConnectionMaxIdleSeconds = defaultConnectionMaxIdleSeconds
	}
	if db.SlowThresholdMilliseconds == 0 {
		db.SlowThresholdMilliseconds = defaultSlowThresholdMilliseconds
	}
}

type Encryptor struct {
	S string `json:"s,omitempty"`
}
