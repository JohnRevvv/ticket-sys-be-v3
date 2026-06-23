#########################################

# MULTI DATABASE CONNECTION SETUP GUIDE

#########################################
Multi-Database Connection Setup Guide

Overview
This guide walks you through setting up multiple PostgreSQL database connections in the Golang application using encrypted credentials. Each database can have its own host, username, and password.

Step 1: Create Local Databases
Before running the application, create the required databases in your local PostgreSQL instance.
Connection Details:
Host: 127.0.0.1
Username: postgres
Password: postgres
Port: 5432
Run the following SQL commands:
sqlCREATE DATABASE cfiad;
CREATE DATABASE innovation;
You can do this via psql, pgAdmin, or any PostgreSQL client of your choice.

Step 2: Configure & Encrypt Your Credentials
The application uses encrypted environment variables for all database credentials. Follow the steps below.
2a. Secret Key
If you don't have a SECRET_KEY yet, the system will auto-generate one when you first run the encryption route. Copy it and save it immediately into your .env file — it won't be shown again.
envSECRET_KEY = <generated_or_your_own_secret_key>

⚠️ Warning: Losing your SECRET_KEY means you will not be able to decrypt your credentials. Keep it safe.

2b. Encrypt Your Credentials
Run the application and call the encryption route provided in the app to encrypt each of your connection values. The route accepts plain text and returns the encrypted string.
Once encrypted, structure your .env file like this:
env

# DO NOT COMMIT THIS FILE — SENSITIVE CONFIGURATION

SECRET_KEY = 9d19708d717e927ef94ea07b6b620873

# REDIS CONNECTION PARAMETERS

REDIS_ADDRESS =
REDIS_PASSWORD =

# DATABASE 1 — CFIAD

DB_CFIAD_HOST = <encrypted_host>
DB_CFIAD_USER = <encrypted_username>
DB_CFIAD_PASS = <encrypted_password>
DB_CFIAD_PORT = 5432
DB_CFIAD_NAME = <encrypted_db_name>
DB_CFIAD_SSL = disable
DB_CFIAD_TZ = Asia/Manila

# DATABASE 2 — INNOVATION

DB_INNOVATION_HOST = <encrypted_host>
DB_INNOVATION_USER = <encrypted_username>
DB_INNOVATION_PASS = <encrypted_password>
DB_INNOVATION_PORT = 5432
DB_INNOVATION_NAME = <encrypted_db_name>
DB_INNOVATION_SSL = disable
DB_INNOVATION_TZ = Asia/Manila

# SSL CONFIGURATION

POSTGRES*TIMEZONE = Asia/Manila
POSTGRES_SSL_MODE = disable
SSL_CERTIFICATE =
SSL_KEY =
2c. Adding More Databases
To add a new database, simply follow the same naming pattern using the prefix DB*<NAME>\_ for each variable. The application will automatically detect and connect to it on startup.
envDB_ANALYTICS_HOST = <encrypted_host>
DB_ANALYTICS_USER = <encrypted_username>
DB_ANALYTICS_PASS = <encrypted_password>
DB_ANALYTICS_PORT = 5432
DB_ANALYTICS_NAME = <encrypted_db_name>
DB_ANALYTICS_SSL = disable
DB_ANALYTICS_TZ = Asia/Manila

✅ No code changes needed — just add the env variables and restart the app.

Step 3: Using the Database Connections
Once the application starts, all configured databases are connected and stored in config.DBConnMap, keyed by their prefix name.
Basic Usage
go// Accessing the CFIAD database
config.DBConnMap["CFIAD"].Table("table_name").Find(&result)

// Accessing the INNOVATION database
config.DBConnMap["INNOVATION"].Table("table_name").Find(&result)
Common GORM Operations
go// SELECT
var users []User
config.DBConnMap["CFIAD"].Table("users").Where("active = ?", true).Find(&users)

// INSERT
config.DBConnMap["INNOVATION"].Table("projects").Create(&newProject)

// UPDATE
config.DBConnMap["CFIAD"].Table("users").Where("id = ?", id).Updates(&updatedUser)

// DELETE
config.DBConnMap["INNOVATION"].Table("projects").Where("id = ?", id).Delete(&Project{})
Using with Models
go// With GORM model
config.DBConnMap["CFIAD"].Model(&User{}).Where("id = ?", userID).First(&user)

// With raw SQL
config.DBConnMap["INNOVATION"].Raw("SELECT \* FROM projects WHERE status = ?", "active").Scan(&projects)

ℹ️ Note: The key name (e.g. "CFIAD", "INNOVATION") must match the prefix used in your .env file exactly, in uppercase.

Quick Reference
Env Prefix PatternDescriptionDB*<NAME>\_HOSTEncrypted database hostDB*<NAME>_USEREncrypted database usernameDB_<NAME>_PASSEncrypted database passwordDB_<NAME>_PORTDatabase port (plain integer)DB_<NAME>_NAMEEncrypted database nameDB_<NAME>_SSLSSL mode (disable / require)DB_<NAME>\_TZTimezone (e.g. Asia/Manila)

Troubleshooting
FAILED TO CONNECT TO <NAME> — Check that the encrypted values are correct and that the database exists.
FAILED TO PING <NAME> — The database is unreachable. Verify your host, port, and firewall settings.
Decryption error — Your SECRET*KEY may not match the one used during encryption. Make sure it's consistent across your .env.
Connection not found (nil map key) — The key you're using in DBConnMap["KEY"] doesn't match any DB*<KEY>\_HOST prefix in your .env. Keys are case-sensitive and uppercase.

##################################

# DOCKER CONFIGURATION

##################################
-- Docker Build
docker build -t test:latest .

-- Docker Run
docker run -p 127.0.0.1:8000:8000 test:latest
