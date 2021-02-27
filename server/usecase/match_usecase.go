package usecase

import (
	"mahjong/model/taku"
	"mahjong/utils"
	"strconv"

	"github.com/k-jun/northpole"
	"github.com/k-jun/northpole/room"
	"github.com/k-jun/northpole/storage"
	"github.com/k-jun/northpole/user"
)

type MatchUsecase interface {
	JoinRandomRoom(user.User) (string, error)
}

type matchUsecaseImpl struct {
	matches  northpole.Match
	write    func(string) error
	read     func([]byte) error
	callback func(string) error
}

func NewMatchUsecase(matches northpole.Match, write func(string) error, read func([]byte) error, callback func(string) error) MatchUsecase {
	return &matchUsecaseImpl{
		matches:  matches,
		read:     read,
		write:    write,
		callback: callback,
	}

}

func (uc *matchUsecaseImpl) JoinRandomRoom(u user.User) (string, error) {
	rc, err := uc.matches.JoinRandomRoom(u)
	if err != nil {
		if err == storage.RoomStorageRoomNotFound {
			rc, err = uc.CreateRoom(u)
			if err != nil {
				return "", err
			}
		}
	}

	room, _ := <-rc
	go uc.deadCheck(u, room)
	if err := uc.write(roomStatus(room)); err != nil {
		return "", err
	}

	for {
		_, isOpen := <-rc
		if !isOpen {
			return room.ID(), nil
		}
		if err := uc.write(roomStatus(room)); err != nil {
			return "", err
		}
	}
}

func (uc *matchUsecaseImpl) deadCheck(u user.User, room room.Room) {
	for {
		if err := uc.write(""); err != nil {
			if room != nil {
				// connection end
				uc.matches.LeaveRoom(u, room)
			}
			break
		}
		if room != nil && !room.IsOpen() {
			break
		}
	}
}

func (uc *matchUsecaseImpl) CreateRoom(u user.User) (chan room.Room, error) {
	newId := utils.NewUUID()
	newRoom := room.New(newId.String(), taku.MaxNumberOfUsers, uc.callback)
	return uc.matches.CreateRoom(u, newRoom)
}

func roomStatus(r room.Room) string {
	message := "current number of users : " + strconv.Itoa(r.CurrentNumberOfUsers()) + "\n"
	message += "max number of users     : " + strconv.Itoa(r.MaxNumberOfUsers()) + "\n"

	return message
}
