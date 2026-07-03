package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
	encrypDecryptV1 "ideyanale-be/pkg/middleware/encryption/v1"

	IAdmodel "ideyanale-be/pkg/modules/insti-admin/model"
	SAdmodel "ideyanale-be/pkg/modules/super-admin/model"
	Umodel "ideyanale-be/pkg/modules/users/model"
	Tmodel "ideyanale-be/pkg/modules/tickets/model"
	Instimodel "ideyanale-be/pkg/modules/institution/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DBConnList         []*gorm.DB
	DBConnListInternal []*gorm.DB
	DBErr              error
	RedisClient        *redis.Client
	RedisError         error

	JWTSecret = "your-super-secret-key"

	SecretKey string // ✅ GLOBAL SECRET KEY
)

type DatabaseConfig struct {
	HostNum  int
	Host     string
	Username string
	Password string
	Port     int
	SSLMode  string
	Timezone string
	DBNames  []string
}

// DecryptDBConfig reads environment variables and returns configs for all database hosts
func DecryptDBConfig() (map[int]*DatabaseConfig, error) {
	configs := make(map[int]*DatabaseConfig)

	// ============================================
	// STEP 1: Get the secret key from environment
	// ============================================
	secretKey := strings.TrimSpace(utils_v1.GetEnv("SECRET_KEY"))
	if secretKey == "" {
		return nil, fmt.Errorf("SECRET_KEY environment variable is required for decryption")
	}

	log.Printf("[CONFIG] Secret key loaded (length: %d chars)", len(secretKey))

	// ============================================
	// STEP 2: Find all unique host numbers from environment variables
	// ============================================
	hostNums := make(map[int]bool)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "POSTGRES_HOST_") {
			parts := strings.SplitN(env, "=", 2)
			numStr := strings.TrimPrefix(parts[0], "POSTGRES_HOST_")
			numStr = strings.TrimSpace(numStr)
			num, err := strconv.Atoi(numStr)
			if err == nil {
				hostNums[num] = true
			}
		}
	}

	if len(hostNums) == 0 {
		return nil, fmt.Errorf("no database hosts found in environment variables (looking for POSTGRES_HOST_1, POSTGRES_HOST_2, etc.)")
	}

	// ============================================
	// STEP 3: Load configuration for each host
	// ============================================
	for hostNum := range hostNums {
		config := &DatabaseConfig{
			HostNum: hostNum,
		}

		// Read encrypted credentials from environment
		encryptedHost := strings.TrimSpace(utils_v1.GetEnv(fmt.Sprintf("POSTGRES_HOST_%d", hostNum)))
		encryptedUsername := strings.TrimSpace(utils_v1.GetEnv(fmt.Sprintf("POSTGRES_USERNAME_%d", hostNum)))
		encryptedPassword := strings.TrimSpace(utils_v1.GetEnv(fmt.Sprintf("POSTGRES_PASSWORD_%d", hostNum)))

		// Decrypt using the same function from your encryption package
		var err error
		config.Host, err = encrypDecryptV1.DecryptV1(encryptedHost, secretKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt POSTGRES_HOST_%d: %v (encrypted value: %s)", hostNum, err, encryptedHost)
		}

		config.Username, err = encrypDecryptV1.DecryptV1(encryptedUsername, secretKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt POSTGRES_USERNAME_%d: %v", hostNum, err)
		}

		config.Password, err = encrypDecryptV1.DecryptV1(encryptedPassword, secretKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt POSTGRES_PASSWORD_%d: %v", hostNum, err)
		}

		// Read non-encrypted configuration
		portStr := strings.TrimSpace(utils_v1.GetEnv(fmt.Sprintf("POSTGRES_PORT_%d", hostNum)))
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port for host %d: %v", hostNum, err)
		}
		config.Port = port
		config.SSLMode = strings.TrimSpace(utils_v1.GetEnv(fmt.Sprintf("POSTGRES_SSL_MODE_%d", hostNum)))
		config.Timezone = strings.TrimSpace(utils_v1.GetEnv(fmt.Sprintf("POSTGRES_TIMEZONE_%d", hostNum)))

		// Validate decrypted values
		if config.Host == "" {
			return nil, fmt.Errorf("POSTGRES_HOST_%d decrypted to empty string", hostNum)
		}
		if config.Username == "" {
			return nil, fmt.Errorf("POSTGRES_USERNAME_%d decrypted to empty string", hostNum)
		}

		// ============================================
		// STEP 4: Get all database names for this host
		// ============================================
		prefix := fmt.Sprintf("DB_%d_NAME_", hostNum)
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, prefix) {
				parts := strings.SplitN(env, "=", 2)
				if len(parts) == 2 {
					encryptedDBName := strings.TrimSpace(parts[1])
					if encryptedDBName != "" {
						// Decrypt the database name
						dbName, err := encrypDecryptV1.DecryptV1(encryptedDBName, secretKey)
						if err != nil {
							return nil, fmt.Errorf("failed to decrypt database name for host %d: %v", hostNum, err)
						}
						if dbName != "" {
							config.DBNames = append(config.DBNames, dbName)
						}
					}
				}
			}
		}

		if len(config.DBNames) == 0 {
			return nil, fmt.Errorf("no databases found for host %d (looking for %s* variables)", hostNum, prefix)
		}

		configs[hostNum] = config
		log.Printf("[CONFIG] Loaded Host %d: %s (decrypted) with %d database(s)", hostNum, config.Host, len(config.DBNames))
	}

	return configs, nil
}

// PostgreSQLConnect establishes connections to all configured databases
func PostgreSQLConnect() bool {
	log.Println("========================================")
	log.Println("Starting PostgreSQL Connections...")
	log.Println("========================================")

	SecretKey = strings.TrimSpace(utils_v1.GetEnv("SECRET_KEY"))
	if SecretKey == "" {
		log.Fatal("SECRET_KEY is required")
	}

	configs, err := DecryptDBConfig()
	if err != nil {
		fmt.Printf("❌ Database config error: %s\n", err.Error())
		return false
	}

	var sortedHostNums []int
	for hostNum := range configs {
		sortedHostNums = append(sortedHostNums, hostNum)
	}
	sort.Ints(sortedHostNums)

	connectionIndex := 0
	for _, hostNum := range sortedHostNums {
		config := configs[hostNum]

		fmt.Printf("\n📍 Host %d: %s (User: %s)\n", hostNum, config.Host, config.Username)
		fmt.Printf("   Timezone: %s | SSL: %s\n", config.Timezone, config.SSLMode)

		for _, dbName := range config.DBNames {
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
				config.Host,
				config.Username,
				config.Password,
				dbName,
				config.Port,
				config.SSLMode,
				config.Timezone,
			)

			dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Error),
			})
			if err != nil {
				fmt.Printf("❌ FAILED TO CONNECT TO DATABASE '%s' on %s: %v\n", dbName, config.Host, err)
				return false
			}

			sqlDB, err := dbConn.DB()
			if err != nil {
				fmt.Printf("❌ FAILED TO GET DATABASE INSTANCE for '%s': %v\n", dbName, err)
				return false
			}

			err = sqlDB.Ping()
			if err != nil {
				fmt.Printf("❌ FAILED TO PING DATABASE '%s': %v\n", dbName, err)
				return false
			}

			// AUTO MIGRATE TABLES
			err = dbConn.AutoMigrate(

				//superadmin model
				&SAdmodel.SuperAccount{},

				// user model
				&Umodel.UserDetails{},
				&Umodel.LoginOTP{},

				//insti admin model
				&IAdmodel.JobPosition{},
				&IAdmodel.TicketType{},
				&IAdmodel.Category{},
				&IAdmodel.SubCategory{},
				&IAdmodel.Roles{},

				//Ticket Model
				&Tmodel.Ticket{},
				&Tmodel.TicketAttachment{},

				//Institution Model
				&Instimodel.Institution{},
				&Instimodel.InstitutionLogo{},
			)

			if err != nil {
				fmt.Printf("❌ AutoMigrate failed: %v\n", err)
				return false
			}

			DBConnList = append(DBConnList, dbConn)
			fmt.Printf("   ✔ [Index %d] %s\n", connectionIndex, strings.ToUpper(dbName))
			connectionIndex++
		}
	}

	fmt.Println("\n========================================")
	fmt.Printf("✅ Total database connections: %d\n", len(DBConnList))
	fmt.Println("========================================\n")

	return true
}

func RedisConnect(address, password string) bool {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	ping, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("❌ Can't ping redis:", err)
		return false
	}

	fmt.Println("✅ PING REDIS:", ping)
	return true
}
