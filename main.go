package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

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

type Options struct {
	Cipher     string
	InputFile  string
	OutputFile string
	Password   string
}

func main() {
	opt := Options{}

	flag.StringVar(&opt.Cipher, "c", "pbkdf2", "Cipher (scrypt|pbkdf2)")
	flag.StringVar(&opt.InputFile, "f", "", "Keystore file (if empty read stdin)")
	flag.StringVar(&opt.OutputFile, "o", "", "Keystore file (if empty read stdin)")
	flag.StringVar(&opt.Password, "p", "", "Keystore password")
	flag.Parse()

	if err := run(opt); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run(opt Options) error {
	var err error

	encryptor := keystorev4.New(keystorev4.WithCipher(opt.Cipher))
	keystore := &Keystore{}

	if opt.InputFile != "" && opt.InputFile == opt.OutputFile {
		return errors.New("input and output files must be different")
	}

	input := os.Stdin
	if opt.InputFile != "" {
		input, err = os.Open(opt.InputFile)
		if err != nil {
			return errors.New("could not read keystore file")
		}
	}

	output := os.Stdout
	if opt.OutputFile != "" {
		output, err = os.Create(opt.OutputFile)
		if err != nil {
			return errors.New("could not write to file")
		}
	}

	err = json.NewDecoder(input).Decode(&keystore)
	if err != nil {
		return err
	}

	if keystore.Pubkey == "" {
		return errors.New("could not decode keystore json")
	}

	secret, err := encryptor.Decrypt(keystore.Crypto, opt.Password)
	if err != nil {
		return err
	}

	crypto, err := encryptor.Encrypt(secret, opt.Password)
	if err != nil {
		return err
	}

	keystore.Crypto = crypto

	keystore2, err := json.Marshal(keystore)
	if err != nil {
		return err
	}

	fmt.Fprint(output, string(keystore2))

	return nil
}
