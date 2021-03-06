package storage

import "mahjong/model/board"

type BoardStorage interface {
	Add(string, board.Board) error
	Remove(string) error
	Find(string) (board.Board, error)
}

type boardStorageImpl struct {
	boards map[string]board.Board
}

func NewBoardStorage() BoardStorage {
	return &boardStorageImpl{boards: make(map[string]board.Board)}
}

func (ts *boardStorageImpl) Add(id string, t board.Board) error {
	if ts.boards[id] != nil {
		return BoardStorageAlreadyExistErr
	}

	ts.boards[id] = t
	return nil
}

func (ts *boardStorageImpl) Remove(id string) error {
	if ts.boards[id] == nil {
		return BoardStorageNotExistErr
	}

	ts.boards[id] = nil
	return nil
}

func (ts *boardStorageImpl) Find(id string) (board.Board, error) {
	if ts.boards[id] == nil {
		return nil, BoardStorageNotExistErr
	}

	return ts.boards[id], nil
}
