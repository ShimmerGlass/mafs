package num

const (
	TypeSignedInt = "sint"
	TypeFloat     = "float"
)

type Number interface {
	Add(Number) (Number, error)
	Sub(Number) (Number, error)
	Div(Number) (Number, error)
	Mul(Number) (Number, error)
	Mod(Number) (Number, error)

	Inverse() (Number, error)
}
