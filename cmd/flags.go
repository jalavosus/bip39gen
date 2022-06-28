package main

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/jalavosus/bip39gen/internal/datakeys"
	"github.com/jalavosus/bip39gen/internal/outformat"
	"github.com/jalavosus/bip39gen/internal/types"
)

const (
	defaultOutFormat     = outformat.JSON
	defaultNum       int = 3
)

const (
	categoryOutput    string = "output options"
	categoryGenParams string = "generator options"
)

var defaultOutExcludes = map[string]bool{
	datakeys.Mnemonic:       false,
	datakeys.Entropy:        false,
	datakeys.Privkey:        false,
	datakeys.Pubkey:         false,
	datakeys.Seed:           false,
	datakeys.WalletIndex:    false,
	datakeys.DerivationPath: false,
}

var (
	outFileFlag = cli.PathFlag{
		Name:     "outfile",
		Aliases:  []string{"o"},
		Usage:    "`path` to file to write output to. If not provided, data is printed to stdout.",
		Required: false,
		Category: categoryOutput,
	}

	outFormatFlag = AllowedStringValuesFlag{
		StringFlag: &cli.StringFlag{
			Name:     "format",
			Aliases:  []string{"f"},
			Usage:    "`format` to use for outputting data to. Allowed values: json,yaml,text",
			Required: false,
			Value:    defaultOutFormat.String(),
			Category: categoryOutput,
		},
		AllowedValues: []string{
			outformat.JSON.String(),
			outformat.YAML.String(),
			outformat.TOML.String(),
			outformat.Text.String(),
		},
	}

	outExcludesFlag = AllowedStringValuesFlag{
		StringFlag: &cli.StringFlag{
			Name:     "exclude",
			Aliases:  []string{"e"},
			Usage:    "comma-separated list of `values` to exclude from output. Allowed values: mnemonic,seed,entropy,privkey,pubkey,wallet_index,derivation_path",
			Required: false,
			Value:    "",
			Category: categoryOutput,
		},
		AllowEmpty: true,
		AllowedValues: []string{
			datakeys.Mnemonic,
			datakeys.Seed,
			datakeys.Entropy,
			datakeys.Privkey,
			datakeys.Pubkey,
			datakeys.WalletIndex,
			datakeys.DerivationPath,
		},
	}

	numAddressesFlag = cli.IntFlag{
		Name:     "num",
		Aliases:  []string{"n"},
		Usage:    "`num`ber of addresses to generate.",
		Required: false,
		Value:    defaultNum,
		Category: categoryGenParams,
	}

	mnemonicFlag = cli.StringFlag{
		Name:     "mnemonic",
		Aliases:  []string{"m"},
		Usage:    "Use the provided `mnemonic` for wallet generation instead of randomly generating one.",
		Required: false,
		Value:    "",
		Category: categoryGenParams,
	}

	mnemonicLenFlag = cli.IntFlag{
		Name:     "mnemonic-len",
		Aliases:  []string{"l"},
		Usage:    "Length of mnemonic. Allowed values: 12 15 18 21 24",
		Required: false,
		Value:    24,
		Category: categoryGenParams,
	}

	passphraseFlag = cli.StringFlag{
		Name:     "passphrase",
		Aliases:  []string{"p"},
		Usage:    "pass`phrase` for seed generation.",
		Required: false,
		Value:    "",
		Category: categoryGenParams,
	}

	oneMnemonicFlag = cli.BoolFlag{
		Name:     "single-mnemonic",
		Aliases:  []string{"s"},
		Usage:    "If passed, all generated addresses use the same mnemonic. Not recommended.",
		Required: false,
		Value:    false,
		Category: categoryGenParams,
	}

	genNameFlag = cli.BoolFlag{
		Name:     "gen-name",
		Aliases:  []string{"g"},
		Usage:    "[Optional] generate a sort of \"name\" for addresses from their mnemonic",
		Required: false,
		Value:    false,
		Category: categoryGenParams,
	}

	hardendedFlag = cli.BoolFlag{
		Name:     "hardended",
		Aliases:  []string{"d"},
		Usage:    "[Optional] use `hardended` derivation path indices",
		Required: false,
		Value:    false,
		Category: categoryGenParams,
	}

	sequentialIndexFlag = cli.BoolFlag{
		Name:     "sequential-index",
		Usage:    "[Optional] If true and --single-mnemonic is true, addresses are generated with sequential wallet indices.",
		Required: false,
		Value:    false,
		Category: categoryGenParams,
	}
)

func parseFlags(c *cli.Context) (params types.CLIParams, err error) {
	rawExcludes := strings.Split(
		strings.Replace(outExcludesFlag.Get(c), " ", "", -1),
		",",
	)

	var excludesMap = make(map[string]bool)

	for k, v := range defaultOutExcludes {
		excludesMap[k] = v
	}

	for _, e := range rawExcludes {
		excludesMap[e] = true
	}

	var (
		mnemonic      = mnemonicFlag.Get(c)
		validMnemonic = false
	)

	if mnemonic != "" {
		if err = types.ValidateMnemonic(mnemonic); err != nil {
			err = errors.Wrap(err, "error validating provided mnemonic")
			return
		} else {
			validMnemonic = true
		}
	}

	if !validMnemonic {
		err = validateMnemonicLength(c)
		if err != nil {
			return
		}
	}

	num := numAddressesFlag.Get(c)

	outFile := outFileFlag.Get(c)
	if outFile != "" && !filepath.IsAbs(outFile) {
		outFile, err = filepath.Abs(outFile)
		if err != nil {
			return
		}
	}

	sequentialIndices := sequentialIndexFlag.Get(c) && oneMnemonicFlag.Get(c)

	params = types.CLIParams{
		Num:               num,
		OutfilePath:       outFile,
		OutFormat:         outformat.FromString(outFormatFlag.Get(c)),
		ExcludeFromOutput: excludesMap,
		Passphrase:        passphraseFlag.Get(c),
		Mnemonic:          mnemonic,
		MnemonicLength:    mnemonicLenFlag.Get(c),
		ValidMnemonic:     validMnemonic,
		GenName:           genNameFlag.Get(c),
		Hardened:          hardendedFlag.Get(c),
		SequentialIndex:   sequentialIndices,
	}

	return
}

func validateMnemonicLength(c *cli.Context) (err error) {
	ml := mnemonicLenFlag.Get(c)

	switch ml {
	case 12:
	case 15:
	case 18:
	case 21:
	case 24:
	default:
		err = errors.Errorf("invalid mnemonic length %[1]d. Allowed values: 12 15 18 21 24", ml)
	}

	return

}
