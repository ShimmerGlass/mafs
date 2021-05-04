package num

import "math"

type Float float64

func (i Float) Add(o Number) (Number, error) {
	return Float(i + o.(Float)), nil
}
func (i Float) Sub(o Number) (Number, error) {
	return Float(i - o.(Float)), nil
}
func (i Float) Div(o Number) (Number, error) {
	return Float(i / o.(Float)), nil
}
func (i Float) Mul(o Number) (Number, error) {
	return Float(i * o.(Float)), nil
}
func (i Float) Mod(o Number) (Number, error) {
	return Float(math.Mod(float64(i), float64(o.(Float)))), nil
}

func (i Float) Inverse() (Number, error) {
	return -i, nil
}
