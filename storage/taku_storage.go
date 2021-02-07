package storage

import (
	"mahjong/model/taku"
)

type TakuStorage interface {
	Add(string, taku.Taku) error
	Remove(string) error
	Find(string) (taku.Taku, error)
}

type takuStorageImpl struct {
	takus map[string]taku.Taku
}

func NewTakuStorage() TakuStorage {
	return &takuStorageImpl{takus: make(map[string]taku.Taku)}
}

func (ts *takuStorageImpl) Add(id string, t taku.Taku) error {
	if ts.takus[id] != nil {
		return TakuStorageAlreadyExistErr
	}

	ts.takus[id] = t
	return nil
}

func (ts *takuStorageImpl) Remove(id string) error {
	if ts.takus[id] == nil {
		return TakuStorageNotExistErr
	}

	ts.takus[id] = nil
	return nil
}

func (ts *takuStorageImpl) Find(id string) (taku.Taku, error) {
	if ts.takus[id] == nil {
		return nil, TakuStorageNotExistErr
	}

	return ts.takus[id], nil
}
