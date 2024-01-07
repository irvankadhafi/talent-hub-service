package console

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/irvankadhafi/talent-hub-service/internal/config"
	"github.com/irvankadhafi/talent-hub-service/internal/db"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
	"github.com/irvankadhafi/talent-hub-service/internal/repository"
	"github.com/irvankadhafi/talent-hub-service/internal/usecase"
	"github.com/irvankadhafi/talent-hub-service/pkg/cacher"
	"github.com/irvankadhafi/talent-hub-service/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed-candidate",
	Short: "run seed-candidate",
	Long:  `This subcommand seeding candidate`,
	Run:   seedUser,
}

func init() {
	RootCmd.AddCommand(seedCmd)
}

func seedUser(cmd *cobra.Command, args []string) {
	// Initiate all connection like db, redis, etc
	db.InitializePostgresConn()

	cacheManager := cacher.ConstructCacheManager()

	if !config.DisableCaching() {
		redisDB, err := db.InitializeRedigoRedisConnectionPool(config.RedisCacheHost(), redisOptions)
		continueOrFatal(err)
		defer utils.WrapCloser(redisDB.Close)

		cacheManager.SetConnectionPool(redisDB)
	}

	cacheManager.SetDisableCaching(config.DisableCaching())

	candidateRepo := repository.NewCandidateRepository(db.PostgreSQL, cacheManager)
	candidateUsecase := usecase.NewCandidateUsecase(candidateRepo)

	for i := 0; i < 10; i++ { // Number of candidates to seed
		var candidate model.CreateCandidateInput

		// Generate random data for each field
		candidate.FullName = faker.Name()              // Random full name
		candidate.Email = faker.Email()                // Random email
		candidate.Password = "Password123"             // Use a fixed password
		candidate.PasswordConfirmation = "Password123" // Same as password

		// Set Gender manually as faker does not provide Gender
		candidate.Gender = model.GenderMale // or model.GenderFemale

		// Create the candidate
		_, err := candidateUsecase.Create(context.Background(), candidate)
		if err != nil {
			logrus.Error("Error creating candidate: %v", err)
			continue
		}
	}

	logrus.Warn("DONE!")
}
