package usecase

import (
	"mahjong/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/k-jun/northpole"
	"github.com/k-jun/northpole/room"
	"github.com/k-jun/northpole/storage"
	"github.com/k-jun/northpole/user"
)

var (
	MaxNumberOfUsers = 4
)

type MatchUsecase interface {
	JoinRandomRoom(user.User) (uuid.UUID, error)
}

type matchUsecaseImpl struct {
	matches  northpole.Match
	write    func(string) error
	read     func([]byte) error
	callback func(uuid.UUID) error
}

func NewMatchUsecase(matches northpole.Match, write func(string) error, read func([]byte) error, callback func(uuid.UUID) error) MatchUsecase {
	return &matchUsecaseImpl{
		matches:  matches,
		read:     read,
		write:    write,
		callback: callback,
	}

}

func (uc *matchUsecaseImpl) JoinRandomRoom(u user.User) (uuid.UUID, error) {
	var room room.Room
	rc, err := uc.matches.JoinRandomRoom(u)
	if err != nil {
		if err == storage.RoomStorageRoomNotFound {
			rc, err = uc.CreateRoom(u)
			if err != nil {
				return uuid.Nil, err
			}
		}
	}

	// dead check
	go func() {
		for {
			if err := uc.write(""); err != nil {
				if room != nil {
					uc.matches.LeaveRoom(u, room)
				}
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for {
		room = <-rc
		if room == nil {
			return uuid.Nil, nil
		}

		if err := uc.write(roomStatus(room)); err != nil {
			return uuid.Nil, err
		}
		if !room.IsOpen() {
			break
		}
	}

	return room.ID(), nil
}

func (uc *matchUsecaseImpl) CreateRoom(u user.User) (chan room.Room, error) {
	newId := utils.NewUUID()
	newRoom := room.New(newId, MaxNumberOfUsers, uc.callback)
	return uc.matches.CreateRoom(u, newRoom)
}

func roomStatus(r room.Room) string {
	message := "current number of users : " + strconv.Itoa(r.CurrentNumberOfUsers()) + "\n"
	message += "max number of users     : " + strconv.Itoa(r.MaxNumberOfUsers()) + "\n"

	return message
}
