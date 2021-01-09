package usecase

import (
	"mahjong/cha"
	"mahjong/storage"
	"time"

	"github.com/google/uuid"
)

type GameUsecase interface {
	JoinTaku(uuid.UUID, cha.Cha) error
}

type gameUsecaseImpl struct {
	takuStorage storage.TakuStorage
	read        func([]byte) error
	write       func(string) error
}

func NewGameUsecase(ts storage.TakuStorage, write func(string) error, read func([]byte) error) GameUsecase {
	return &gameUsecaseImpl{
		takuStorage: ts,
		read:        read,
		write:       write,
	}

}

func (gu *gameUsecaseImpl) JoinTaku(id uuid.UUID, c cha.Cha) error {
	taku, err := gu.takuStorage.Find(id)
	if err != nil {
		return err
	}

	ct, err := taku.JoinCha(c)
	if err != nil {
		return err
	}

	// dead check
	go func() {
		for {
			if err := gu.write(""); err != nil {
				taku.LeaveCha(c)
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for {
		isClose := <-ct
		if isClose == nil {
			break
		}
		// TODO shell art
	}
	return nil
}
