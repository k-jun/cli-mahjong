package storage

import (
	"mahjong/taku"

	"github.com/google/uuid"
)

type TakuStorage interface {
	Add(uuid.UUID, taku.Taku) error
	Remove(uuid.UUID) error
	Find(uuid.UUID) (taku.Taku, error)
}

type takuStorageImpl struct {
	takus map[uuid.UUID]taku.Taku
}

func NewTakuStorage() TakuStorage {
	return &takuStorageImpl{takus: make(map[uuid.UUID]taku.Taku)}
}

func (ts *takuStorageImpl) Add(id uuid.UUID, t taku.Taku) error {
	if ts.takus[id] != nil {
		return TakuStorageAlreadyExistErr
	}

	ts.takus[id] = t
	return nil
}

func (ts *takuStorageImpl) Remove(id uuid.UUID) error {
	if ts.takus[id] == nil {
		return TakuStorageNotExistErr
	}

	ts.takus[id] = nil
	return nil
}

func (ts *takuStorageImpl) Find(id uuid.UUID) (taku.Taku, error) {
	if ts.takus[id] == nil {
		return nil, TakuStorageNotExistErr
	}

	return ts.takus[id], nil
}
