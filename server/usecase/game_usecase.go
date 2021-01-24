package usecase

import (
	"mahjong/cha"
	"mahjong/storage"

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

	// input
	go func() {
		for {
			buffer := make([]byte, 2048)
			err := gu.read(buffer)
			if err != nil {
				taku.LeaveCha(c)
				break
			}
			if string(buffer) != "" {
				if taku.IsYourTurn(c) {
					c.Dahai(c.TumoHai())
					// TODO huro check
					taku.TurnChange(taku.NextTurn())
					taku.Broadcast()
				}

			}
		}
	}()

	// output
	for {
		isClose := <-ct
		if isClose == nil {
			break
		}

		if taku.IsYourTurn(c) && c.TumoHai() == nil {
			err := c.Tumo()
			if err != nil {
				// game end
				gu.write("thank you for playing!" + "\n")
				taku.LeaveCha(c)
			}
		}

		// TODO shell art
		tehaistr := ""
		for _, h := range c.Tehai().Hais() {
			tehaistr += h.Name()
		}

		if c.TumoHai() != nil {
			tehaistr += " " + c.TumoHai().Name()
		}
		gu.write(tehaistr + "\n")
	}
	return nil
}
