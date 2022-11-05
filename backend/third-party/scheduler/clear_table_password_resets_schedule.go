package scheduler

import (
	"context"
	"log"
)

type ClearTablePasswordResetsSchedule struct {
}

func (scheduler *ClearTablePasswordResetsSchedule) Run() {
	repo := repository.NewPasswordResetRepository()

	ctx := context.Background()
	configuration := config.New()
	db, err := server.NewInitializedDatabase(configuration)
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer util.CommitOrRollback(tx)

	err = repo.Truncate(ctx, tx)
	if err != nil {
		log.Fatalln(err)
	}
}
