package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sanches1984/json-parser/app/config"
	"github.com/sanches1984/json-parser/app/reporter"
	"os"
)

func main() {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	logger.Debug().Msg("read config")

	cfg, err := config.Load("config.json")
	if err != nil {
		logger.Fatal().Err(err).Msg("can't read config")
	}

	file, err := os.Open(cfg.Path)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't open file")
	}
	defer file.Close()

	logger.Debug().Str("path", cfg.Path).Msg("read data from file")

	reporterService := reporter.New(cfg, logger)
	if err := reporterService.Read(file); err != nil {
		logger.Fatal().Err(err).Msg("can't read file")
	}

	logger.Debug().Msg("make report")

	jsonData, err := reporterService.MakeReport().ToJSON()
	if err != nil {
		logger.Fatal().Err(err).Msg("can't marshal json")
	}

	fmt.Println(string(jsonData))
}
