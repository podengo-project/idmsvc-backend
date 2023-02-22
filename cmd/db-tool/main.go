package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hmsidm/internal/config"
	"github.com/hmsidm/internal/infrastructure/datastore"
)

func printUsage(args []string) {
	fmt.Printf(`
Usage: %s {new,migrate} [opts...]

  new <migration-name>

    migration-name   eg 'create_table_todo'


	migrate {up,down} [step]

`, args[0])
}

func migrate(config *config.Config, direction string, steps int) error {
	switch direction {
	case "up":
		{
			return datastore.MigrateUp(config, steps)
		}
	case "down":
		{
			return datastore.MigrateDown(config, steps)
		}
	}
	return fmt.Errorf("func 'migrate' not implemented")
}

func main() {
	var err error
	// TODO Refactor in a better way
	//      Adopt error return as soon as possible
	//      Encapsulate argument parse and checks
	//      Encapsulate main body in a 'run' function
	if len(os.Args) == 1 {
		printUsage(os.Args)
		os.Exit(1)
	}

	config := config.Get()

	switch os.Args[1] {
	case "new":
		if len(os.Args) == 2 {
			printUsage(os.Args)
			os.Exit(1)
		}
		if err = datastore.CreateMigrationFile(os.Args[2]); err != nil {
			panic(err)
		}
	case "migrate":
		var steps int
		if len(os.Args) == 2 || len(os.Args) == 3 {
			printUsage(os.Args)
			panic("not enough parameters")
		}
		steps, err = strconv.Atoi(os.Args[3])
		if err != nil {
			printUsage(os.Args)
			panic(err)
		}
		if err = migrate(config, os.Args[2], steps); err != nil {
			panic(err)
		}
	}

	return
}
