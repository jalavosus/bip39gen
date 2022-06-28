package bip39gen

import (
	"strings"

	"github.com/jalavosus/bip39gen/internal/datakeys"
	"github.com/jalavosus/bip39gen/internal/utils"
)

// AddressData contains generated data for an address derived from a
// hdwallet.
type AddressData struct {
	Address        string
	Name           string
	Entropy        string
	PubKey         string
	PrivKey        string
	Mnemonic       string
	Seed           string
	WalletIndex    int
	DerivationPath string
	Hardened       bool
}

func (a AddressData) BuildOutput() (ad AddressDataOutput) {
	ad = AddressDataOutput{
		Address:     utils.ToPointer(a.Address),
		Hardened:    utils.ToPointer(a.Hardened),
		WalletIndex: utils.ToPointer(a.WalletIndex),
	}

	if !checkZeroVal(a.Name) {
		ad.Name = utils.ToPointer(a.Name)
	}

	if !checkZeroVal(a.Entropy) {
		ad.Entropy = utils.ToPointer(a.Entropy)
	}

	if !checkZeroVal(a.PubKey) {
		ad.Pubkey = utils.ToPointer(a.PubKey)
	}

	if !checkZeroVal(a.PrivKey) {
		ad.Privkey = utils.ToPointer(a.PrivKey)
	}

	if !checkZeroVal(a.Mnemonic) {
		ad.Mnemonic = strings.Split(a.Mnemonic, " ")
	}

	if !checkZeroVal(a.Seed) {
		ad.Seed = utils.ToPointer(a.Seed)
	}

	if !checkZeroVal(a.DerivationPath) {
		ad.DerivationPath = utils.ToPointer(a.DerivationPath)
	}

	return
}

func (a AddressData) FormatOutput(excludes map[string]bool) AddressDataOutput {
	if excludes == nil {
		excludes = make(map[string]bool)
	}

	out := a.BuildOutput()

	for _, field := range datakeys.FieldOrder {
		switch field {
		case datakeys.Name:
			checkExclude(excludes, field, utils.ToPointer(out.Name))
		case datakeys.Entropy:
			checkExclude(excludes, field, utils.ToPointer(out.Entropy))
		case datakeys.Mnemonic:
			checkExclude(excludes, field, utils.ToPointer(out.Mnemonic))
		case datakeys.Seed:
			checkExclude(excludes, field, utils.ToPointer(out.Seed))
		case datakeys.Pubkey:
			checkExclude(excludes, field, utils.ToPointer(out.Pubkey))
		case datakeys.Privkey:
			checkExclude(excludes, field, utils.ToPointer(out.Privkey))
		case datakeys.WalletIndex:
			checkExclude(excludes, field, utils.ToPointer(out.WalletIndex))
		case datakeys.DerivationPath:
			checkExclude(excludes, field, utils.ToPointer(out.DerivationPath))
		case datakeys.Hardened:
			checkExclude(excludes, field, utils.ToPointer(out.Hardened))
		}
	}

	return out
}

func checkZeroVal[T comparable](val T) bool {
	var zeroVal T

	return val == zeroVal
}

func checkExclude[T any](excludes map[string]bool, field string, data *T) {
	var zeroVal T

	val, ok := excludes[field]
	if ok && val {
		*data = zeroVal
	}
}
