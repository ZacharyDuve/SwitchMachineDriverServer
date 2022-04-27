package environment

import (
	"os"

	"github.com/ZacharyDuve/apireg/environment"
)

func GetCurrent() environment.Environment {
	if os.Getenv("environment") == "production" {
		return environment.Prod
	} else {
		return environment.NonProd
	}
}
