package usecase

import (
	"mahjong/taku"
	"mahjong/utils"
	"strconv"
	"time"

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
	var room room.Room
	rc, err := uc.matches.JoinRandomRoom(u)
	if err != nil {
		if err == storage.RoomStorageRoomNotFound {
			rc, err = uc.CreateRoom(u)
			if err != nil {
				return "", err
			}
		}
	}

	// dead check
	go func() {
		for {
			if err := uc.write(""); err != nil {
				if room != nil {
					// connection end
					uc.matches.LeaveRoom(u, room)
				}
				break
			}
			time.Sleep(100 * time.Millisecond)
			if !room.IsOpen() {
				break
			}
		}
	}()

	for {
		room, isClose := <-rc
		if isClose || uc.write(roomStatus(room)) != nil {
			return "", err
		}

		if !room.IsOpen() {
			break
		}
	}

	return room.ID(), nil
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
