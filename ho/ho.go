package ho

import "mahjong/hai"

type Ho interface {
	Add(inHai *hai.Hai) error
	Last() (*hai.Hai, error)
}

type hoImpl struct {
	hais []*hai.Hai
}

func New() Ho {

	return &hoImpl{}

}

func (h *hoImpl) Add(inHai *hai.Hai) error {
	h.hais = append(h.hais, inHai)
	return nil
}

func (h *hoImpl) Last() (*hai.Hai, error) {
	if len(h.hais) == 0 {
		return nil, HoNoHaiError

	}

	return h.hais[len(h.hais)-1], nil
}
