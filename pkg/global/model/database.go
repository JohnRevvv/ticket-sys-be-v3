package model

type (
	DBEntry struct {
		Host     string
		Username string
		Password string
		Port     int
		SSLMode  string
		Timezone string
		DBName   string
	}

	Database struct {
		SecretKey string
		DBList    map[string]DBEntry // key = prefix (e.g. "GAMIFICATION")
	}

	Redis struct {
		RedisAddress string
		Password     string
	}
)
