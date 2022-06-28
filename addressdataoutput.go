package bip39gen

import (
	"strconv"
	"strings"

	"github.com/jalavosus/bip39gen/internal/utils"
)

// AddressDataOutput is only used for data output,
// and is almost always created by AddressData.BuildOutput.
type AddressDataOutput struct {
	Address        *string  `json:"address" yaml:"address" toml:"address"`
	Name           *string  `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Pubkey         *string  `json:"pubkey,omitempty" yaml:"pubkey,omitempty" toml:"pubkey,omitempty"`
	Privkey        *string  `json:"privkey,omitempty" yaml:"privkey,omitempty" toml:"privkey,omitempty"`
	Mnemonic       []string `json:"mnemonic,omitempty" yaml:"mnemonic,omitempty" toml:"mnemonic,omitempty"`
	Seed           *string  `json:"seed,omitempty" yaml:"seed,omitempty" toml:"seed,omitempty"`
	Entropy        *string  `json:"entropy,omitempty" yaml:"entropy,omitempty" toml:"entropy,omitempty"`
	WalletIndex    *int     `json:"wallet_index,omitempty" yaml:"wallet_index,omitempty" toml:"wallet_index,omitempty"`
	DerivationPath *string  `json:"derivation_path,omitempty" yaml:"derivation_path,omitempty" toml:"derivation_path,omitempty"`
	Hardened       *bool    `json:"hardened,omitempty" yaml:"hardened,omitempty" toml:"hardened,omitempty"`
}

func (a AddressDataOutput) mnemonic() *string {
	if a.Mnemonic != nil {
		return utils.ToPointer(strings.Join(a.Mnemonic, " "))
	}

	return nil
}

func (a AddressDataOutput) walletIndex() *string {
	if a.WalletIndex != nil {
		return utils.ToPointer(strconv.Itoa(*a.WalletIndex))
	}

	return nil
}

func (a AddressDataOutput) hardened() *string {
	if a.Hardened != nil {
		return utils.ToPointer(strconv.FormatBool(*a.Hardened))
	}

	return nil
}

type AddressDataOutputSlice []AddressDataOutput

func (s AddressDataOutputSlice) ads() ads {
	return ads{s}
}
