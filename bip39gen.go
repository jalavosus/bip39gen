package bip39gen

import (
	"math/rand"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jalavosus/hdwallet-go"
	"github.com/tyler-smith/go-bip39"

	"github.com/jalavosus/bip39gen/internal/types"
)

var (
	rando *rand.Rand
)

func init() {
	seed := time.Now().UnixNano()

	// call this so external libraries utilizing
	// math/rand have a decent seed source.
	rand.Seed(seed)

	src := rand.NewSource(seed)
	rando = rand.New(src)
}

// GenerateAddressesSingleMnemonic generates a slice of AddressData structs, all of which
// are initialized using the same entropy, seed, and mnemonic.
func GenerateAddressesSingleMnemonic(params *types.GeneratorParams, num int) (addrs []AddressData) {
	addrs = make([]AddressData, num)

	for i := 0; i < num; i++ {
		var idx = i

		if !params.SequentialIndex {
			randInt := rando.Intn(1000) + 1

			if params.Hardened {
				randInt += hdkeychain.HardenedKeyStart
			}

			idx = randInt
		}

		addrs[i] = generateAddress(params, idx)
	}

	return
}

// GenerateAddress returns a AddressData struct initialized
// using random values for entropy, seed, mnemonic, and wallet index.
func GenerateAddress(params *types.GeneratorParams) AddressData {
	randInt := rando.Intn(1000) + 1
	if params.Hardened {
		randInt += hdkeychain.HardenedKeyStart
	}

	return generateAddress(params, randInt)
}

func generateAddress(params *types.GeneratorParams, idx int) AddressData {
	var newWalletParams = []hdwallet.NewWalletOpt{
		hdwallet.WithPassphrase(params.Passphrase),
	}

	if params.Mnemonic != "" {
		newWalletParams = append(newWalletParams, hdwallet.WithMnemonic(params.Mnemonic))
	}

	if params.Entropy != nil {
		newWalletParams = append(newWalletParams, hdwallet.WithEntropy(params.Entropy))
	}

	return makeDerivedAddress(idx, params.GenName, newWalletParams...)
}

func GenerateMnemonic(length int) []string {
	var entropyLen int
	switch length {
	case 12:
		entropyLen = hdwallet.Entropy128Bit
	case 15:
		entropyLen = hdwallet.Entropy160Bit
	case 18:
		entropyLen = hdwallet.Entropy192Bit
	case 21:
		entropyLen = hdwallet.Entropy224Bit
	case 24:
		entropyLen = hdwallet.Entropy256Bit
	}

	entropy, _ := bip39.NewEntropy(entropyLen)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	split := strings.Split(mnemonic, " ")
	if len(split) != length {
		panic("wtf")
	}

	return split
}

func GenerateMnemonicAndEntropy(length int) (mnemonic string, entropy []byte) {
	var entropyLen int
	switch length {
	case 12:
		entropyLen = hdwallet.Entropy128Bit
	case 15:
		entropyLen = hdwallet.Entropy160Bit
	case 18:
		entropyLen = hdwallet.Entropy192Bit
	case 21:
		entropyLen = hdwallet.Entropy224Bit
	case 24:
		entropyLen = hdwallet.Entropy256Bit
	}

	entropy, _ = bip39.NewEntropy(entropyLen)
	mnemonic, _ = bip39.NewMnemonic(entropy)

	split := strings.Split(mnemonic, " ")
	if len(split) != length {
		panic("wtf")
	}

	return
}

func makeDerivedAddress(pathIdx int, genName bool, params ...hdwallet.NewWalletOpt) AddressData {
	wallet, err := hdwallet.NewHDWallet(params...)
	if err != nil {
		panic(err)
	}

	derivedAccount, err := wallet.DeriveAddressFromIndex(pathIdx)
	if err != nil {
		panic(err)
	}

	rawAddrData := AddressData{
		Address:        derivedAccount.Address().String(),
		Entropy:        common.Bytes2Hex(wallet.Entropy()),
		Seed:           common.Bytes2Hex(wallet.Seed()),
		Mnemonic:       wallet.Mnemonic(),
		PubKey:         derivedAccount.PublicKeyHex(),
		PrivKey:        derivedAccount.PrivateKeyHex(),
		WalletIndex:    derivedAccount.DerivationIndex(),
		DerivationPath: derivedAccount.DerivationPath(),
		Hardened:       derivedAccount.Hardened(),
	}

	if genName {
		rawAddrData.Name = genNameFromMnemonic(rawAddrData.Mnemonic)
	}

	return rawAddrData
}

func genNameFromMnemonic(mnemonic string) string {
	split := strings.Split(mnemonic, " ")

	getWord := func() string {
		return split[rand.Intn(len(split))]
	}

	var (
		n1 = getWord()
		n2 string
	)

	for n2 = getWord(); (n2 == n1) || len(n2) < 5; {
		n2 = getWord()
	}

	return n1 + " " + n2
}
