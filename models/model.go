package models

import (
	"container/list"

	"github.com/MessageDream/webIM/modules/safemap"
)

type User struct {
	Id      int64
	Name    string
	IsAdmin bool
}

type EventType int

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
)

type Event struct {
	Type      EventType // JOIN, LEAVE, MESSAGE
	User      User
	Timestamp int // Unix timestamp (secs)
	Content   string
	Room      string
}

const archiveSize = 20

// Event archives.
var roomArchive = safemap.NewSafeMap()

// NewArchive saves new event to archive list.
func NewArchive(roomid string, event Event) {
	var archive *list.List
	if roomArchive.Check(roomid) {
		archive = roomArchive.Get(roomid).(*list.List)
		if archive.Len() >= archiveSize {
			archive.Remove(archive.Front())
		}
	} else {
		archive = list.New()
		roomArchive.Add(roomid, archive)
	}
	archive.PushBack(event)
}

// GetEvents returns all events after lastReceived.
func GetEvents(roomid string, lastReceived int) []Event {
	var archive *list.List
	if roomArchive.Check(roomid) {
		archive = roomArchive.Get(roomid).(*list.List)
		events := make([]Event, 0, archive.Len())
		for event := archive.Front(); event != nil; event = event.Next() {
			e := event.Value.(Event)
			if e.Timestamp > int(lastReceived) {
				events = append(events, e)
			}
		}
		return events
	}
	return nil
}
