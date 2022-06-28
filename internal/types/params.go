package types

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"

	"github.com/jalavosus/bip39gen/internal/outformat"
)

type CLIParams struct {
	OutFormat         outformat.OutFormat
	OutfilePath       string
	ExcludeFromOutput map[string]bool
	Num               int
	Passphrase        string
	Mnemonic          string
	MnemonicLength    int
	ValidMnemonic     bool
	GenName           bool
	Hardened          bool
	SequentialIndex   bool
}

func (p CLIParams) GeneratorParams() *GeneratorParams {
	return &GeneratorParams{
		OutFormat:       p.OutFormat,
		Excludes:        p.ExcludeFromOutput,
		Passphrase:      p.Passphrase,
		Mnemonic:        p.Mnemonic,
		MnemonicLen:     p.MnemonicLength,
		GenName:         p.GenName,
		Hardened:        p.Hardened,
		SequentialIndex: p.SequentialIndex,
	}
}

type GeneratorParams struct {
	OutFormat       outformat.OutFormat
	Excludes        map[string]bool
	Passphrase      string
	Mnemonic        string
	MnemonicLen     int
	Entropy         []byte
	GenName         bool
	Hardened        bool
	SequentialIndex bool
}

func ValidateMnemonic(mnemonic string) (mnemonicErr error) {
	mnemonicLen := len(strings.Split(mnemonic, " "))

	if mnemonicLen < 12 || mnemonicLen > 24 {
		mnemonicErr = errors.Errorf(
			"provided mnemonic must be at least 12 words and no more than 24 words; got %d words",
			mnemonicLen,
		)

		return
	}

	if mnemonicLen%3 != 0 {
		mnemonicErr = errors.Errorf(
			"provided mnemonic must be 12, 15, 18, 21, or 24 words long (divisible by 3); got %d words",
			mnemonicLen,
		)

		return
	}

	_, mnemonicErr = bip39.EntropyFromMnemonic(mnemonic)

	return
}
