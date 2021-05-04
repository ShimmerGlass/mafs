package main

import (
	"fmt"
	"strconv"
)

type cmdFunc func(args ...string) error

func (u *UI) cmdSetBase(args ...string) error {
	if u.ctx.Type == typeFloat {
		return fmt.Errorf("cannot change base in float mode")
	}
	if len(args) != 1 {
		return fmt.Errorf("this command accepts 1 argument")
	}

	b, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	err = validateBase(b)
	if err != nil {
		return err
	}

	u.ctx.Base = b
	return nil
}

func (u *UI) cmdSetDisplayedBases(args ...string) error {
	if u.ctx.Type == typeFloat {
		return fmt.Errorf("cannot change displayed bases in float mode")
	}

	bases := []int{}
	for _, a := range args {
		b, err := strconv.Atoi(a)
		if err != nil {
			return err
		}
		err = validateBase(b)
		if err != nil {
			return err
		}

		bases = append(bases, b)
	}
	u.bases = bases
	return nil
}

func (u *UI) cmdSetType(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("this command accepts 1 argument")
	}

	if args[0] != typeFloat && args[0] != typeSignedInt {
		return fmt.Errorf("bad type")
	}

	u.ctx.Type = args[0]
	return nil
}

func validateBase(b int) error {
	if b < 2 {
		return fmt.Errorf("base must be >= 2")
	}
	if b > 36 {
		return fmt.Errorf("base must be <= 36")
	}

	return nil
}
