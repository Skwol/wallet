package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/skwol/wallet/pkg/client/pgdb"
	"github.com/skwol/wallet/pkg/fixtures"
	"github.com/skwol/wallet/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()

	if err := run(logger); err != nil {
		logger.Fatal("run walletctl:", err.Error())
		os.Exit(1)
	}
}

type config struct {
	storage struct {
		dsn string
	}
}

func run(logger logging.Logger) (err error) {
	cfg := config{}
	logger.Info("starting service ", "service_name ", "walletctl")

	flagParser := kingpin.New("walletctl", "The tool to manage wallet")
	flagParser = flagParser.DefaultEnvars()
	flagParser.HelpFlag.Short('h')

	flagParser.Flag("storage.dsn", "Connection string for PostgresSQL").
		Default("postgres://wallet_user:psw@localhost:5432/wallet_db?sslmode=disable").
		Envar("STORAGE_DSN").
		StringVar(&cfg.storage.dsn)

	loadFixturesCmd := flagParser.Command("load-fixtures", "Load flagParser 'fake' set of data into flagParser database.")

	command := kingpin.MustParse(flagParser.Parse(os.Args[1:]))

	db, err := pgdb.NewClient("production")
	if err != nil {
		return err
	}
	defer db.Conn.Close()

	if loadFixturesCmd.FullCommand() == command {
		if err := fixtures.Load(db.Conn, "./db/fixtures"); err != nil {
			return err
		}
	}

	return nil
}
