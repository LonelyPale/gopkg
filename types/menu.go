package types

type Enum struct {
	value int
	name  string
}

type FormatType int

const (
	FTByte FormatType = iota + 1
	FTArray
	FTDefine
)

func (ft FormatType) String() string {
	switch ft {
	case FTByte:
		return "byte"
	case FTArray:
		return "array"
	case FTDefine:
		return "define"
	default:
		return ""
	}
}
