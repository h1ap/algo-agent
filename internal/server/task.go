package server

import (
	"algo-agent/internal/middleware/cron"
	"algo-agent/internal/service"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

func NewTaskServer(service *service.CronServer, logger log.Logger) *cron.Server {
	helper := log.NewHelper(logger)
	ctx := context.Background()

	checkTrainingTaskFunc := func() {
		err := service.RunCheckTrainingTaskService(ctx)
		if err != nil {
			helper.Error(err)
			return
		}
	}

	checkEvalTaskFunc := func() {
		err := service.RunCheckEvalTaskService(ctx)
		if err != nil {
			helper.Error(err)
			return
		}
	}

	checkDeployServiceFunc := func() {
		err := service.RunCheckDeployService(ctx)
		if err != nil {
			helper.Error(err)
			return
		}
	}

	checkExtractTaskFunc := func() {
		err := service.RunCheckExtractTaskService(ctx)
		if err != nil {
			helper.Error(err)
			return
		}
	}

	srv := cron.NewServer(
		cron.Logger(logger),
		cron.WithContext(ctx),

		cron.RegisterFunc("@every 30s", "checkTrainingTaskFunc", func() {
			checkTrainingTaskFunc()
		}),

		cron.RegisterFunc("@every 30s", "checkEvalTaskFunc", func() {
			checkEvalTaskFunc()
		}),

		cron.RegisterFunc("@every 30s", "checkDeployServiceFunc", func() {
			checkDeployServiceFunc()
		}),

		cron.RegisterFunc("@every 30s", "checkExtractTaskFunc", func() {
			checkExtractTaskFunc()
		}),
	)

	return srv
}
