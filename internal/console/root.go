package console

import (
	"github.com/irvankadhafi/talent-hub-service/internal/config"
	"github.com/irvankadhafi/talent-hub-service/internal/db"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var redisOptions *db.RedisConnectionPoolOptions

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cobra-example",
	Short: "An example of cobra",
	Long: `This application shows how to create modern CLI
			applications in go using Cobra CLI library`,
}

// Execute runs the root command for the application and handles any errors that may occur during execution.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolP("integration", "i", false, "use integration config file")
	config.GetConf()

	redisOptions = &db.RedisConnectionPoolOptions{
		DialTimeout:     config.RedisDialTimeout(),
		ReadTimeout:     config.RedisReadTimeout(),
		WriteTimeout:    config.RedisWriteTimeout(),
		IdleCount:       config.RedisMaxIdleConn(),
		PoolSize:        config.RedisMaxActiveConn(),
		IdleTimeout:     240 * time.Second,
		MaxConnLifetime: 1 * time.Minute,
	}
}
