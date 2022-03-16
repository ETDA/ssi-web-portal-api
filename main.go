package main

import (
	"fmt"
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/web"
	"ssi-gitlab.teda.th/ssi/core"
)

func main() {
	switch core.NewEnv().Config().Service {
	case string(consts.ServiceWeb):
		web.Run()
	default:
		fmt.Printf("Service not found")
	}
}
