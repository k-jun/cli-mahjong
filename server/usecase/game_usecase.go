package usecase

import (
	"bytes"
	"fmt"
	"mahjong/cha"
	"mahjong/hai"
	"mahjong/storage"
	"mahjong/taku"
	"strings"
)

type GameUsecase interface {
	JoinTaku(string, cha.Cha) (chan taku.Taku, error)
	InputController(string, cha.Cha) error
	OutputController(string, cha.Cha, chan taku.Taku) error
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

func (gu *gameUsecaseImpl) InputController(id string, c cha.Cha) error {
	taku, err := gu.takuStorage.Find(id)
	if err != nil {
		return err

	}
	for {
		buffer := make([]byte, 1024)
		err := gu.read(buffer)
		if err != nil {
			// dead check
			taku.LeaveCha(c)
			break
		}
		if string(buffer) != "" {
			if taku.IsYourTurn(c) {
				buffer = bytes.Trim(buffer, "\x00")
				buffer = bytes.Trim(buffer, "\x10")
				haiName := strings.TrimSpace(string(buffer))
				outHai := hai.AtoHai(haiName)
				fmt.Println("outhai:", outHai.Name())
				if outHai == nil {
					fmt.Println("unknown hai: ", haiName)
					continue
				}
				c.Dahai(outHai)
				// TODO huro check
				taku.TurnChange(taku.NextTurn())
				taku.Broadcast()
			}
		}
	}

	return nil
}

func (gu *gameUsecaseImpl) OutputController(id string, c cha.Cha, channel chan taku.Taku) error {
	taku, err := gu.takuStorage.Find(id)
	if err != nil {
		return err
	}
	for {
		isClose := <-channel
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
			tehaistr += " " + h.Name()
		}

		if c.TumoHai() != nil {
			tehaistr += " | " + c.TumoHai().Name()
		}
		gu.write(tehaistr + "\n")
	}

	return nil
}

func (gu *gameUsecaseImpl) JoinTaku(id string, c cha.Cha) (chan taku.Taku, error) {
	taku, err := gu.takuStorage.Find(id)
	if err != nil {
		return nil, err
	}

	return taku.JoinCha(c)
}
