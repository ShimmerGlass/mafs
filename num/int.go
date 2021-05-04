package num

type SignedInt int64

func (i SignedInt) Add(o Number) (Number, error) {
	return SignedInt(i + o.(SignedInt)), nil
}
func (i SignedInt) Sub(o Number) (Number, error) {
	return SignedInt(i - o.(SignedInt)), nil
}
func (i SignedInt) Div(o Number) (Number, error) {
	return SignedInt(i / o.(SignedInt)), nil
}
func (i SignedInt) Mul(o Number) (Number, error) {
	return SignedInt(i * o.(SignedInt)), nil
}
func (i SignedInt) Mod(o Number) (Number, error) {
	return SignedInt(i % o.(SignedInt)), nil
}

func (i SignedInt) Inverse() (Number, error) {
	return -i, nil
}
