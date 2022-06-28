package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/jalavosus/bip39gen"
	"github.com/jalavosus/bip39gen/internal/outformat"
	"github.com/jalavosus/bip39gen/internal/types"
)

const (
	genCmdName string = "gen"
)

var genCmd = cli.Command{
	Name:                   genCmdName,
	Usage:                  "Generate an address",
	UseShortOptionHandling: true,
	Flags: []cli.Flag{
		&numAddressesFlag,
		&outFormatFlag,
		&outExcludesFlag,
		&outFileFlag,
		&mnemonicFlag,
		&mnemonicLenFlag,
		&passphraseFlag,
		&oneMnemonicFlag,
		&genNameFlag,
		&hardendedFlag,
		&sequentialIndexFlag,
	},
	Action: genCmdAction,
}

func genCmdAction(c *cli.Context) error {
	params, paramsErr := parseFlags(c)
	if paramsErr != nil {
		return paramsErr
	}

	var (
		addrs      = make([]bip39gen.AddressData, params.Num)
		addrsCheck = make(map[string]bool)
	)

	genParams := params.GeneratorParams()

	if oneMnemonicFlag.Get(c) {
		genParams.Mnemonic, genParams.Entropy = bip39gen.GenerateMnemonicAndEntropy(genParams.MnemonicLen)
		addrs = bip39gen.GenerateAddressesSingleMnemonic(genParams, params.Num)
	} else {
		for i := 0; i < params.Num; i++ {
			addr := bip39gen.GenerateAddress(genParams)

			// make sure we never get the same address
			_, ok := addrsCheck[addr.Address]
			for ok {
				addr = bip39gen.GenerateAddress(genParams)
				_, ok = addrsCheck[addr.Address]
			}

			addrs[i] = addr
			addrsCheck[addr.Address] = true
		}
	}

	formattedAddrs := make(bip39gen.AddressDataOutputSlice, len(addrs))

	for i, addr := range addrs {
		formattedAddrs[i] = addr.FormatOutput(genParams.Excludes)
	}

	return writeDataOut(formattedAddrs, params)
}

func writeDataOut(data bip39gen.AddressDataOutputSlice, params types.CLIParams) error {
	var (
		marshaled []byte
		formatter types.OutputFormatter
	)

	if len(data) == 1 {
		formatter = data[0]
	} else {
		formatter = data
	}

	switch params.OutFormat {
	case outformat.JSON:
		marshaled = formatter.FormatJSON()
	case outformat.YAML:
		marshaled = formatter.FormatYAML()
	case outformat.TOML:
		marshaled = formatter.FormatTOML()
	case outformat.Text:
		marshaled = formatter.FormatText()
	}

	var out *os.File
	if params.OutfilePath == "" {
		out = os.Stdout
	} else {
		var err error

		out, err = os.OpenFile(params.OutfilePath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}

		defer func() {
			_ = out.Close()
		}()
	}

	_, err := out.WriteString(string(marshaled))

	return err
}
