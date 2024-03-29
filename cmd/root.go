package cmd

import (
	"fmt"
	"github.com/nwtgck/go-piping-duplex"
	"github.com/nwtgck/go-piping-duplex/util"
	"github.com/nwtgck/go-piping-duplex/version"
	"github.com/spf13/cobra"
	"io"
	"os"
)

const (
	ServerUrlEnvName = "PIPING_SERVER"
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
	RootCmd.Flags().StringVarP(&server, "server", "s", defaultServer, "Piping Server URL")
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
			return fmt.Errorf("upload path and download are required")
		}
		var passphrase string
		var err error
		if usesPassphrase {
			passphrase, err = util.InputPassphrase()
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "[INFO] End-to-end encrypted")
		}
		uploadPath := args[0]
		downloadPath := args[1]
		fmt.Fprintf(os.Stderr, "[INFO] Server: %s\n", server)
		fmt.Fprintf(os.Stderr, "[INFO] Your upload path: '%s', your download path: '%s'\n", uploadPath, downloadPath)
		var input io.Reader = os.Stdin
		if usesPassphrase {
			input = util.OpenpgpSymmetricallyEncrypt(input, []byte(passphrase))
		}
		output := os.Stdout
		r, uploadFinishErrCh, err := piping_duplex.DuplexReader(server, uploadPath, downloadPath, input)
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
		_, err = io.Copy(output, r)
		if err != nil {
			return err
		}
		// Wait for uploading
		err = <-uploadFinishErrCh
		if err != nil {
			return err
		}
		return nil
	},
}
