package helpers

type MonitorDecission byte

const (
	UN MonitorDecission = iota
	TOP
	BOTTOM
)

func (d MonitorDecission) String() string {
	switch d {
	case UN:
		return "?"
	case TOP:
		return "\u22A4"
	case BOTTOM:
		return "\u22A5"
	}
	panic("MonitorDecission: Unhandled case in String()")
}
