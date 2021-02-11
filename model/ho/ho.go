package ho

import "mahjong/model/hai"

type Ho interface {
	Add(inHai *hai.Hai) error
	Hais() []*hai.Hai
	Last() (*hai.Hai, error)
	RemoveLast() (*hai.Hai, error)
}

type hoImpl struct {
	hais []*hai.Hai
}

func New() Ho {

	return &hoImpl{}

}

func (h *hoImpl) Hais() []*hai.Hai {
	return h.hais
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

func (h *hoImpl) RemoveLast() (*hai.Hai, error) {
	if len(h.hais) == 0 {
		return nil, HoNoHaiError
	}
	outHai := h.hais[len(h.hais)-1]
	h.hais = h.hais[:len(h.hais)-1]

	return outHai, nil
}
