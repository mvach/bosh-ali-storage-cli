package main

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry/bosh-ali-storage-cli/client"
	"github.com/cloudfoundry/bosh-ali-storage-cli/config"
	"log"
	"os"
)

var version string

func main() {

	configPath := flag.String("c", "", "configuration path")
	showVer := flag.Bool("v", false, "version")
	flag.Parse()

	if *showVer {
		fmt.Printf("version %s\n", version)
		os.Exit(0)
	}

	configFile, err := os.Open(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	aliConfig, err := config.NewFromReader(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	storageClient, err := client.NewStorageClient(aliConfig)
	if err != nil {
		log.Fatalln(err)
	}

	blobstoreClient, err := client.New(storageClient)
	if err != nil {
		log.Fatalln(err)
	}

	nonFlagArgs := flag.Args()
	if len(nonFlagArgs) < 2 {
		log.Fatalf("Expected at least two arguments got %d\n", len(nonFlagArgs))
	}

	cmd := nonFlagArgs[0]

	switch cmd {
	case "put":
		if len(nonFlagArgs) != 3 {
			log.Fatalf("Put method expected 3 arguments got %d\n", len(nonFlagArgs))
		}
		sourceFilePath, destination := nonFlagArgs[1], nonFlagArgs[2]

		_, err := os.Stat(sourceFilePath)
		if err != nil {
			log.Fatalln(err)
		}

		err = blobstoreClient.Put(sourceFilePath, destination)
		fatalLog(cmd, err)

	case "get":
		if len(nonFlagArgs) != 3 {
			log.Fatalf("Get method expected 3 arguments got %d\n", len(nonFlagArgs))
		}
		source, destinationFilePath := nonFlagArgs[1], nonFlagArgs[2]

		err = blobstoreClient.Get(source, destinationFilePath)
		fatalLog(cmd, err)

	default:
		log.Fatalf("unknown command: '%s'\n", cmd)
	}
}

func fatalLog(cmd string, err error) {
	if err != nil {
		log.Fatalf("performing operation %s: %s\n", cmd, err)
	}
}
