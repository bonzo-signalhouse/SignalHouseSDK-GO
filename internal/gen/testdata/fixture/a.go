package fixture

type Zeta struct{ n int }

type Alpha interface{ Do() }

type hidden struct{}

func (z Zeta) Method() int { return z.n }

func NewZeta() Zeta { return Zeta{} }

func helper() {}

const (
	ModeA = iota
	ModeB
	modeHidden
)

var DefaultName, otherName = "x", "y"
