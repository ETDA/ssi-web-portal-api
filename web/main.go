package web

import (
	"fmt"
	"os"

	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
)

func Run() {
	env := core.NewEnv()

	mysql, err := core.NewDatabase(env.Config()).Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "MySQL: %v", err)
		os.Exit(1)
	}
	contextOptions := &core.ContextOptions{
		DB:  mysql,
		ENV: env,
	}
	ctx := core.NewContext(contextOptions)
	vcService := services.NewVCService(ctx)
	go vcService.VCSigning()

	e := core.NewHTTPServer(&core.HTTPContextOptions{
		ContextOptions: contextOptions,
	})

	NewHomeHTTPHandler(e)

	core.StartHTTPServer(e, env)
}
