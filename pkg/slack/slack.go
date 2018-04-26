package slack

import (
	slackClient "github.com/nlopes/slack"
	"fmt"
)

type (
	Slack struct {
		token string
		api   *slackClient.Client
		rtm   *slackClient.RTM
		eventList chan Event
	}

	Message interface {
		Text() string
		Channel() string
		UserID() string
	}

	Event struct {
		name EventName
		data interface{}
	}

	EventName string

	UserInfo struct {
		id string
	}
)

const (
	eventBufferSize = 1024
)

var (
	paramsMessageSlack = slackClient.PostMessageParameters{AsUser: true, EscapeText: false}
)

func New(token string) *Slack {
	return &Slack{
		token: token,
	}
}

func (s *Slack) Listen() chan Event {
	s.api = slackClient.New(s.token)
	s.eventList = make(chan Event, eventBufferSize)

	go s.manageMessages()

	return s.eventList
}

func (s *Slack) manageMessages() {
	s.rtm = s.api.NewRTM()
	go s.rtm.ManageConnection()
	for rtmEvent := range s.rtm.IncomingEvents {
		s.eventList <- Event{
			name: EventName(rtmEvent.Type),
			data: rtmEvent.Data,
		}
	}
}

func (s *Slack) Send(msg Message) error {
	_, _, err := s.api.PostMessage(msg.Channel(), msg.Text(), paramsMessageSlack)
	return err
}

func (s *Slack) UserInfo() (*UserInfo, error) {
	info, err := s.api.AuthTest()
	if err != nil {
		return nil, fmt.Errorf("AuthTest err: %s", err.Error())
	}
	userInfo := &UserInfo{
		id:info.UserID,
	}
	return userInfo, nil
}

func (e *Event) Name() EventName {
	return e.name
}

func (e *Event) Data() interface{} {
	return e.data
}

func (u UserInfo) ID() string {
	return u.id
}