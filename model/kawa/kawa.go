package kawa

import "mahjong/model/hai"

type Kawa interface {
	Add(inHai *hai.Hai) error
	Hais() []*hai.Hai
	Last() (*hai.Hai, error)
	RemoveLast() (*hai.Hai, error)
}

type kawaImpl struct {
	hais []*hai.Hai
}

func New() Kawa {

	return &kawaImpl{}

}

func (h *kawaImpl) Hais() []*hai.Hai {
	return h.hais
}

func (h *kawaImpl) Add(inHai *hai.Hai) error {
	h.hais = append(h.hais, inHai)
	return nil
}

func (h *kawaImpl) Last() (*hai.Hai, error) {
	if len(h.hais) == 0 {
		return nil, KawaNoHaiError
	}
	return h.hais[len(h.hais)-1], nil
}

func (h *kawaImpl) RemoveLast() (*hai.Hai, error) {
	if len(h.hais) == 0 {
		return nil, KawaNoHaiError
	}
	outHai := h.hais[len(h.hais)-1]
	h.hais = h.hais[:len(h.hais)-1]

	return outHai, nil
}
