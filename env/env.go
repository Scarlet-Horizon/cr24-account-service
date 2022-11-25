package env

import (
	"github.com/joho/godotenv"
)

func Load(fileName string) error {
	return godotenv.Load(fileName)
}
