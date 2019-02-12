//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/arangodb-managed/apis/common/auth"
	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// RootCmd is the root (and only) command of this service
	RootCmd = &cobra.Command{
		Use:   "oasis",
		Short: "ArangoDB Oasis",
		Long:  "ArangoDB Oasis. The Managed Cloud for ArangoDB",
		Run:   showUsage,
	}

	cliLog   = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	rootArgs struct {
		token    string
		endpoint string
		format   format.Options
	}
)

const (
	// Prefix of all environment variables
	envKeyPrefix  = "OASIS_"
	apiPortSuffix = ":443"
)

func init() {
	f := RootCmd.PersistentFlags()
	// Persistent flags
	defaultToken := envOrDefault("TOKEN", "")
	defaultEndpoint := envOrDefault("ENDPOINT", "cloud.adbtest.xyz")
	f.StringVar(&rootArgs.token, "token", defaultToken, "Token used to authenticate at ArangoDB Oasis")
	f.StringVar(&rootArgs.endpoint, "endpoint", defaultEndpoint, "API endpoint of the ArangoDB Oasis")
}

// Show usage of the given command
func showUsage(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

// envOrDefault returns the value from an environment value with given key
// or if no such environment variable exists, the given default value.
func envOrDefault(envKeySuffix string, defaultValue string) string {
	if v := os.Getenv(envKeyPrefix + envKeySuffix); v != "" {
		return v
	}
	return defaultValue
}

// mustDialAPI dials the ArangoDB Oasis API
func mustDialAPI() *grpc.ClientConn {
	// Set up a connection to the server.
	tc := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.Dial(rootArgs.endpoint+apiPortSuffix, grpc.WithTransportCredentials(tc))
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to connect to ArangoDB Oasis API")
	}
	return conn
}

// contextWithToken returns a context with access token in it.
func contextWithToken() context.Context {
	if rootArgs.token == "" {
		cliLog.Fatal().Msg("--token missing")
	}
	return auth.WithAccessToken(context.Background(), rootArgs.token)
}

// reqOption returns given value if not empty.
// Fails with clear error message when not set.
// Returns: option-value, number-of-args-used(0|argIndex+1)
func reqOption(key, value string, args []string, argIndex int) (string, int) {
	if value != "" {
		return value,0
	}
	if len(args) > argIndex {
		return args[argIndex], argIndex+1
	}
	cliLog.Fatal().Msgf("--%s missing", key)
	return "",0
}

// mustCheckNumberOfArgs compares the number of arguments with the expected
// number of arguments. 
// If there is a difference a fatal error is raised.
func mustCheckNumberOfArgs(args []string, expectedNumberOfArgs int) {
	if  len(args) > expectedNumberOfArgs{
		cliLog.Fatal().Msg("Too many arguments")
	}
	if  len(args) < expectedNumberOfArgs{
		cliLog.Fatal().Msg("Too few arguments")
	}
}
