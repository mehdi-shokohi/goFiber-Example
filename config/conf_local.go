package conf

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func CheckEnvFile() {
	GetHost()
	GetPort()
	GetMongoAddress()
	GetKeydbAddress()
	GetNatsAddress()
	GetMongodbName()
	GetSecretKey()
	GetEmailPassword()
	GetEmailAddress()
	GetSmtpPort()
	GetSmtpHost()
	GetRedisExpireTimeForEmail()
}

func GetHost() (host string) {
	host = get_env_value("HOST")
	if host == "" {
		panic(errors.New("HOST not found in .env file"))
	}
	return
}

func GetPort() (port string) {
	port = get_env_value("PORT")
	if port == "" {
		panic(errors.New("PORT not found in .env file"))
	}
	return
}

func GetMongoAddress() (mongoAddress string) {
	mongoAddress = get_env_value("MONGO_ADDRESS")
	if mongoAddress == "" {
		panic(errors.New("MONGO_ADDRESS not found in .env file"))
	}
	return
}

func GetKeydbAddress() (keydbAddress string) {
	keydbAddress = get_env_value("KEYDB_ADDRESS")
	if keydbAddress == "" {
		panic(errors.New("KEYDB_ADDRESS not found in .env file"))
	}
	return
}

func GetNatsAddress() (natsAddress string) {
	natsAddress = get_env_value("NATS_ADDRESS")
	if natsAddress == "" {
		panic(errors.New("NATS_ADDRESS not found in .env file"))
	}
	return
}

func GetMongodbName() (mongodbName string) {
	mongodbName = get_env_value("MONGODB_NAME")
	if mongodbName == "" {
		panic(errors.New("MONGODB_NAME not found in .env file"))
	}
	return
}

func GetSecretKey() (secretKey string) {
	secretKey = get_env_value("SECRET_KEY")
	if secretKey == "" {
		panic(errors.New("SECRET_KEY not found in .env file"))
	}
	return
}

func get_env_value(key string) string {
	if RunMode == "normal" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	} else if RunMode == "test" {
		path, _ := os.Getwd()
		err := godotenv.Load(strings.Split(path, ProjectName)[0] + ProjectName + "/.env.test")
		if err != nil {
			log.Fatalf("Error loading .env.test file")
		}
	}

	return os.Getenv(key)
}

func GetEmailAddress() (email string) {
	email = get_env_value("EMAIL-ADDRESS")
	if email == "" {
		panic(errors.New("EMAIL-ADDRESS not found in .env file"))
	}
	return
}

func GetEmailPassword() (password string) {
	password = get_env_value("EMAIL-PASSWORD")
	if password == "" {
		panic(errors.New("EMAIL-PASSWORD not found in .env file"))
	}

	return
}

func GetSmtpHost() (host string) {
	host = get_env_value("SMTP-HOST")
	if host == "" {
		panic(errors.New("SMTP HOST not found in .env file"))
	}
	return
}

func GetSmtpPort() (port string) {
	port = get_env_value("SMTP-PORT")
	if port == "" {
		panic(errors.New("SMTP PORT not found in .env file"))
	}

	return
}

func GetRedisExpireTimeForEmail() (expireTime string) {
	expireTime = get_env_value("REDIS-EXPIRE-TIME-IN-EMAIL")
	if expireTime == "" {
		panic(errors.New("REDIS-EXPIRE-TIME-FOR-EMAIL not found in .env file"))
	}

	return
}
