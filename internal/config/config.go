package config

import (
	"os"

	"github.com/ian77-huang/golang-echo/pkg/cast"

	_ "github.com/joho/godotenv/autoload"
)

func Load() Config {
	minLengthAccount, err := cast.StringToInt(os.Getenv("USER_ACCOUNT_MIN_LENGTH"), 6)
	if err != nil {
		panic(err)
	}
	minLengthPassword, err := cast.StringToInt(os.Getenv("USER_PASSWORD_MIN_LENGTH"), 8)
	if err != nil {
		panic(err)
	}
	maxSizeUserProfileAvatar, err := cast.StringToInt(os.Getenv("MAX_SIZE_USER_PROFILE_AVATAR"), 1024*1024)
	if err != nil {
		panic(err)
	}
	return Config{
		SecretKey:  os.Getenv("SECRET_KEY"),
		AssetsPath: os.Getenv("ASSETS_PATH"),
		RedisURL:   os.Getenv("REDIS_URL"),
		Databases: ConfigDatabases{
			Path: os.Getenv("DATABASE_PATH"),
		},
		Users: ConfigUsers{
			MinLengthAccount:  minLengthAccount,
			MinLengthPassword: minLengthPassword,
		},
		MaxSizeUserProfileAvatar: cast.Megabytes(maxSizeUserProfileAvatar),
	}
}
