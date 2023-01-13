package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"strings"
	"time"
)

var Log = log.Logger

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("|  %-6s |", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s,", i)
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("*****%s*****", i)
	}
	log.Logger = log.Output(output).With().Caller().Logger()
	log.Print("hw")

	// 更多类型 查看https://github.com/rs/zerolog#standard-types
	log.Debug().
		Str("Scale", "4 cents").
		Float64("Interval", 12.03).
		Msg("111")

	err := errors.New("an error")
	err = outer()
	log.Error().Stack().Err(err).Str("a", "b").Msg("error appeared")
	// Fatal中断程序
	log.Fatal().Stack().Err(err).Str("a", "b").Msg("error appeared")
	// Panic中断程序
	log.Panic().Str("a", "b").Msg("panicking")
}

func outer() error {
	return errors.New("an error")
}
