package environment

import (
	"os"

	"git.zmanhobbies.com/software/apireg/environment"
)

func GetCurrent() environment.Environment {
	if os.Getenv("environment") == "production" {
		return environment.Prod
	} else {
		return environment.NonProd
	}
}
