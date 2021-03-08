package tehai

import "errors"

var (
	TehaiReachMaxHaiErr = errors.New("reached to the max hai number")
	TehaiHaiNotFoundErr = errors.New("the hai not found in the tehai")
	TehaiHaiIsNilErr    = errors.New("the hai is nil")
)
