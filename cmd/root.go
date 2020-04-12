package cmd

import (
	"fmt"
	"github.com/mattn/go-tty"
	"github.com/nwtgck/go-piping-duplex"
	"github.com/nwtgck/go-piping-duplex/util"
	"github.com/nwtgck/go-piping-duplex/version"
	"github.com/spf13/cobra"
	"io"
	"os"
)

const (
	ServerUrlEnvName = "PIPING_SERVER_URL"
)

var server string
var showsVersion bool
var usesPassphrase bool

func init() {
	cobra.OnInitialize()
	defaultServer, ok := os.LookupEnv(ServerUrlEnvName)
	if !ok {
		defaultServer = "https://ppng.io"
	}
	RootCmd.Flags().StringVarP(&server,  "server",  "s", defaultServer, "Piping Server URL")
	RootCmd.Flags().BoolVarP(&usesPassphrase, "symmetric", "c", false, "use symmetric passphrase protection")
	RootCmd.Flags().BoolVarP(&showsVersion, "version", "v", false, "show version")
}

var RootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "piping-duplex",
	Long:  "Duplex communication over Piping Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showsVersion {
			fmt.Println(version.Version)
			return nil
		}
		if len(args) != 2 {
			return fmt.Errorf("Your ID and peer ID are required\n")
		}
		var passphrase string
		if usesPassphrase {
			tty, err := tty.Open()
			if err != nil {
				return err
			}
			defer tty.Close()
			fmt.Fprint(tty.Output(), "Passphrase: ")
			passphrase, err = tty.ReadPasswordNoEcho()
			if err != nil {
				return err
			}
			fmt.Fprintln(tty.Output(), "[INFO] End-to-end encrypted")
		}
		var _ = passphrase
		selfId := args[0]
		peerId := args[1]
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Server: %s\n", server)
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Establishing between '%s' and '%s'...\n", selfId, peerId)
		err := piping_duplex.Wait(server, selfId, peerId)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(os.Stderr, "[INFO] Established!")
		var input io.Reader = os.Stdin
		if usesPassphrase {
			input = util.OpenpgpSymmetricallyEncrypt(input, []byte(passphrase))
		}
		output := os.Stdout
		r, err := piping_duplex.Duplex(server, selfId, peerId, input)
		if err != nil {
			return err
		}
		if usesPassphrase {
			var decrypted, err = util.OpenpgpSymmetricallyDecrypt(r, []byte(passphrase))
			if err != nil {
				return err
			}
			r = decrypted
		}
		io.Copy(output, r)
		return nil
	},
}
