package cmd

import (
	"fmt"
	"github.com/nwtgck/go-piping-duplex"
	"github.com/spf13/cobra"
	"os"
)

const (
	ServerUrlEnvName = "PIPING_SERVER_URL"
)

var server string

func init() {
	cobra.OnInitialize()
	defaultServer, ok := os.LookupEnv(ServerUrlEnvName)
	if !ok {
		defaultServer = "https://ppng.io"
	}
	RootCmd.Flags().StringVarP(&server,  "server",  "s", defaultServer, "Piping Server URL")
}

var RootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "piping-duplex",
	Long:  "Duplex communication over Piping Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("Your ID and peer ID are required\n")
		}
		selfId := args[0]
		peerId := args[1]
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Server: %s\n", server)
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Establishing between '%s' and '%s'...\n", selfId, peerId)
		err := piping_duplex.Wait(server, selfId, peerId)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(os.Stderr, "[INFO] Established!")
		input := os.Stdin
		output := os.Stdout
		return piping_duplex.Duplex(server, selfId, peerId, input, output)
	},
}
