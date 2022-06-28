package bip39gen

import (
	"bytes"
	"fmt"

	"github.com/jalavosus/bip39gen/internal/datakeys"
	"github.com/jalavosus/bip39gen/internal/outformat"
)

type ads struct {
	Addresses AddressDataOutputSlice `json:"addresses" yaml:"addresses" toml:"addresses"`
}

func (a AddressDataOutput) FormatJSON() []byte {
	return outformat.JSON.Marshal(a)
}

func (a AddressDataOutput) FormatYAML() []byte {
	return outformat.YAML.Marshal(a)
}

func (a AddressDataOutput) FormatTOML() []byte {
	return outformat.TOML.Marshal(a)
}

func (a AddressDataOutput) FormatText() []byte {
	return outformat.Text.Marshal(a, func(data any, buf *bytes.Buffer) {
		addrData := data.(AddressDataOutput)

		for _, field := range datakeys.FieldOrder {
			var dataVal *string

			switch field {
			case datakeys.Address:
				dataVal = addrData.Address
			case datakeys.Name:
				dataVal = addrData.Name
			case datakeys.Entropy:
				dataVal = addrData.Entropy
			case datakeys.Mnemonic:
				dataVal = addrData.mnemonic()
			case datakeys.Seed:
				dataVal = addrData.Seed
			case datakeys.Pubkey:
				dataVal = addrData.Pubkey
			case datakeys.Privkey:
				dataVal = addrData.Privkey
			case datakeys.WalletIndex:
				dataVal = addrData.walletIndex()
			case datakeys.DerivationPath:
				dataVal = addrData.DerivationPath
			case datakeys.Hardened:
				dataVal = addrData.hardened()
			default:
				continue
			}

			numTabs := 1

			if dataVal != nil {
				switch {
				case len(field) < 7:
					numTabs = 2
				case len(field) > 12:
					numTabs = 0
				}

				writeWithTabs(field, dataVal, numTabs, buf)
			}
		}
	})
}

func (s AddressDataOutputSlice) FormatJSON() []byte {
	return outformat.JSON.Marshal(s)
}

func (s AddressDataOutputSlice) FormatYAML() []byte {
	return outformat.YAML.Marshal(s.ads())
}

func (s AddressDataOutputSlice) FormatTOML() []byte {
	return outformat.TOML.Marshal(s.ads())
}

func (s AddressDataOutputSlice) FormatText() []byte {
	return outformat.Text.Marshal(s, func(data any, buf *bytes.Buffer) {
		for i, addr := range s {
			formatted := addr.FormatText()
			buf.Write(formatted)
			if i < len(s)-1 {
				buf.WriteString("\n")
			}
		}
	})
}

func writeWithTabs(field string, val *string, numTabs int, buf *bytes.Buffer) {
	buf.WriteString(field + ":")
	for i := 0; i < numTabs; i++ {
		buf.WriteString("\t")
	}
	buf.WriteString("  ")
	buf.WriteString(fmt.Sprintf("%v", *val))
	buf.WriteString("\n")
}
