package chat

import (
	"container/list"
	"encoding/json"
	"time"

	"github.com/MessageDream/webIM/models"
	"github.com/MessageDream/webIM/modules/log"
	"github.com/MessageDream/webIM/modules/safemap"

	"github.com/gorilla/websocket"
)

const TIMEOUT = 200

type Subscription struct {
	Archive []models.Event      // All the events from the archive.
	New     <-chan models.Event // New events coming in.
}

type Subscriber struct {
	Name string
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

type ChatRoom struct {
	//Room ID
	ID string
	// Channel for new join users.
	subscribe chan Subscriber
	// Channel for exit users.
	unsubscribe chan string
	// Send events here to publish them.
	publish chan models.Event
	// Long polling waiting list.
	waitingList *list.List
	subscribers *list.List
}

func (this *ChatRoom) Join(user string, ws *websocket.Conn) {
	this.subscribe <- Subscriber{Name: user, Conn: ws}
}

func (this *ChatRoom) Leave(user string) {
	this.unsubscribe <- user
}

func (this *ChatRoom) Close() {
	close(this.subscribe)
	close(this.unsubscribe)
	close(this.publish)
	this.waitingList = nil
	this.subscribers = nil
}

func (this *ChatRoom) EnPublish(ev models.Event) {
	this.publish <- ev
}

func (this *ChatRoom) PushBack(ch chan bool) {
	this.waitingList.PushBack(ch)
}

// This function handles all incoming chan messages.
func (this *ChatRoom) chatroom() {
	for {
		select {
		case sub := <-this.subscribe:
			if !isUserExist(this.subscribers, sub.Name) {
				this.subscribers.PushBack(sub) // Add user to the end of list.
				// Publish a JOIN event.
				this.publish <- newEvent(models.EVENT_JOIN, models.User{Name: sub.Name}, "")
				log.Info("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				log.Info("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-this.publish:
			// Notify waiting list.
			for ch := this.waitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				this.waitingList.Remove(ch)
			}

			this.broadcastWebSocket(event)
			models.NewArchive(this.ID, event)

			if event.Type == models.EVENT_MESSAGE {
				log.Info("Message from ", event.User.Name, ";Content:", event.Content)
			}
		case unsub := <-this.unsubscribe:
			for sub := this.subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == unsub {
					this.subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						log.Error("WebSocket closed:", unsub)
					}
					this.publish <- newEvent(models.EVENT_LEAVE, models.User{Name: unsub}, "") // Publish a LEAVE event.
					break
				}
			}
		case <-time.After(time.Second * TIMEOUT):
			if this.subscribers.Len() == 0 {
				break
			}
		}
	}

	CloseRoom(this)
}

func (this *ChatRoom) broadcastWebSocket(event models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Error("Fail to marshal event:", err)
		return
	}

	for sub := this.subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				this.unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}

var (
	roomContainer *safemap.SafeMap
)

func init() {
	roomContainer = safemap.NewSafeMap()
}

func NewChatRoom(roomid string) *ChatRoom {
	room := &ChatRoom{
		ID:          roomid,
		subscribe:   make(chan Subscriber, 10),
		unsubscribe: make(chan string, 10),
		publish:     make(chan models.Event, 10),
		waitingList: list.New(),
		subscribers: list.New(),
	}
	go room.chatroom()

	addRoom(roomid, room)

	return room
}

func NewEvent(ep models.EventType, user models.User, msg string) models.Event {
	return newEvent(ep, user, msg)
}

func newEvent(ep models.EventType, user models.User, msg string) models.Event {
	return models.Event{ep, user, int(time.Now().Unix()), msg}
}

func CheckRoom(roomid string) bool {
	return checkRoom(roomid)
}

func checkRoom(roomid string) bool {
	return roomContainer.Check(roomid)
}

func addRoom(roomid string, room *ChatRoom) bool {
	return roomContainer.Add(roomid, room)
}

func RemoveRoom(roomid string) {
	roomContainer.Delete(roomid)
}

func CloseRoom(room *ChatRoom) {
	room.Close()
	roomContainer.Delete(room.ID)
}

func GetRoom(roomid string) *ChatRoom {
	return roomContainer.Get(roomid).(*ChatRoom)
}
func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}
