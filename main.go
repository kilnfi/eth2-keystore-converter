package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
)

// Keystore json file representation as a Go struct.
// See https://eips.ethereum.org/EIPS/eip-2335
type Keystore struct {
	Crypto      map[string]interface{} `json:"crypto"`
	Description string                 `json:"description"`
	Pubkey      string                 `json:"pubkey"`
	Path        string                 `json:"path"`
	ID          string                 `json:"uuid"`
	Version     uint                   `json:"version"`
}

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "cipher",
		Value:   "pbkdf2",
		Aliases: []string{"c"},
		Usage:   "Cipher to use to encrypt the keystore (scrypt|pbkdf2)",
	},
	&cli.StringFlag{
		Name:    "input",
		Aliases: []string{"f"},
		Usage:   "Keystore file (if empty read stdin)",
	},
	&cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "Keystore file (if empty read stdin)",
	},
	&cli.StringFlag{
		Name:    "password",
		Value:   "",
		Aliases: []string{"p"},
		Usage:   "Keystore decrypting password",
	},
	&cli.StringFlag{
		Name:  "new-password",
		Value: "",
		Usage: "Set another password for generated keystore",
	},
	&cli.BoolFlag{
		Name:  "raw",
		Value: false,
		Usage: "Print raw key without encryption",
	},
}

func main() {
	app := &cli.App{
		Name:   "eth2-keystore-converter",
		Usage:  "Decrypt or recrypt ethereum keystore files",
		Flags:  Flags,
		Action: Run,
	}

	if err := app.Run(os.Args); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func Run(cCtx *cli.Context) error {
	var (
		err error

		inputFile   = cCtx.String("input")
		outputFile  = cCtx.String("output")
		password    = cCtx.String("password")
		newPassword = cCtx.String("new-password")
		cipher      = cCtx.String("cipher")
		raw         = cCtx.Bool("raw")
	)

	encryptor := keystorev4.New(keystorev4.WithCipher(cipher))
	keystore := &Keystore{}

	if inputFile != "" && inputFile == outputFile {
		return errors.New("input and output files must be different")
	}

	input := os.Stdin
	if inputFile != "" {
		input, err = os.Open(inputFile)
		if err != nil {
			return errors.New("could not read keystore file")
		}
	}

	output := os.Stdout
	if outputFile != "" {
		output, err = os.Create(outputFile)
		if err != nil {
			return errors.New("could not write to file")
		}
	}

	err = json.NewDecoder(input).Decode(&keystore)
	if err != nil {
		return errors.New("could not decode keystore json")
	}

	if keystore.Pubkey == "" {
		return errors.New("could not decode keystore json")
	}

	secret, err := encryptor.Decrypt(keystore.Crypto, password)
	if err != nil {
		return fmt.Errorf("could not decrypt keystore: %w", err)
	}

	if raw {
		fmt.Fprintf(output, "0x%x", secret)
		return nil
	}

	if newPassword != "" {
		password = newPassword
	}

	crypto, err := encryptor.Encrypt(secret, password)
	if err != nil {
		return fmt.Errorf("could not encrypt keystore: %w", err)
	}

	keystore.Crypto = crypto

	keystore2, err := json.Marshal(keystore)
	if err != nil {
		return fmt.Errorf("could not encode keystore to json: %w", err)
	}

	fmt.Fprint(output, string(keystore2))

	return nil
}
