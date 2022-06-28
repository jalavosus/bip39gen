package datakeys

const (
	Address        string = "address"
	Name           string = "name"
	Pubkey         string = "pubkey"
	Privkey        string = "privkey"
	Entropy        string = "entropy"
	Seed           string = "seed"
	Mnemonic       string = "mnemonic"
	WalletIndex    string = "wallet_index"
	DerivationPath string = "derivation_path"
	Hardened       string = "hardened"
)

// FieldOrder manually sets field order for formatted output
// as it's otherwise entirely random
var FieldOrder = []string{
	Address,
	Name,
	Pubkey,
	Privkey,
	Mnemonic,
	Seed,
	Entropy,
	WalletIndex,
	DerivationPath,
	Hardened,
}
