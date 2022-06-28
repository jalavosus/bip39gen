package types

type OutputFormatter interface {
	FormatJSON() []byte
	FormatYAML() []byte
	FormatTOML() []byte
	FormatText() []byte
}
