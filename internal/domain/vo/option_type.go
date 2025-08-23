package vo

type OptionType string

const (
	OptionTypeSingle OptionType = "single"
	OptionTypeMulti  OptionType = "multi"
)

func (o OptionType) String() string {
	return string(o)
}
