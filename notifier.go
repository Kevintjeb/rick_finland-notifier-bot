package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"reflect"
	"os"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Listeners struct {
	Listeners []User
	Report    Report
}

var listeners Listeners

var fileName = os.Getenv("LISTENER_FILE")
const siteName = "https://rickvanfessem.nl/finland"
const messageFormat = "🎉🎉🎉🎉🎉🎉 Rick just uploaded a new post! 🎉🎉🎉🎉🎉🎉\n\nThe title is '%s'.\nThis is post %d.\n\nUse this link to read it : %s"

func Init() {
	listeners, e := getLocalListenersFile()

	if e != nil {
		listeners.saveListeners()
	}
}

func UpdateStats(report Report) bool {
	equal := reflect.DeepEqual(listeners.Report, report)

	if equal {
		fmt.Println("States equal, nothing new!")
		return false
	} else {
		fmt.Println("States not equal, new Data!")
		listeners.Report = report
		listeners.saveListeners()
		return true
	}
}

func NotifySubscribers(bot *tb.Bot) {
	for _, user := range listeners.Listeners {
		message := fmt.Sprintf(messageFormat, listeners.Report.LatestTopic, listeners.Report.TotalTopics, siteName)

		fmt.Println("Sending :\n\n" + message + "\n\nto: " + user.Name)
		user.SendMessage(message, bot)
	}
}

func GetLatestTopic() string {
	return listeners.Report.LatestTopic
}

func GetTotalTopicCount() int {
	return listeners.Report.TotalTopics
}

func AddUser(user User) {
	if !listeners.containsUser(user) {
		listeners.Listeners = append(listeners.Listeners, user)
		listeners.saveListeners()
	}
}

func RemoveUser(user User) bool {
	if listeners.containsUser(user) {
		if i, e := indexOf(user, listeners.Listeners); e == nil {
			listeners.Listeners = append(listeners.Listeners[:i], listeners.Listeners[i+1:]...)
			listeners.saveListeners()
			return true
		}
	}

	return false
}

func (listeners Listeners) containsUser(user User) bool {
	for _, v := range listeners.Listeners {
		if v.Id == user.Id {
			return true
		}
	}
	return false
}

func (listeners Listeners) saveListeners() error {
	bytes, e := json.MarshalIndent(listeners, "", "  ")
	if e != nil {
		return e
	}

	if err := ioutil.WriteFile(fileName, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func getLocalListenersFile() (Listeners, error) {
	bytes, e := ioutil.ReadFile(fileName)

	if e != nil {
		return listeners, e
	}

	json.Unmarshal(bytes, &listeners)
	return listeners, nil
}

func indexOf(element User, list []User) (int, error) {
	for k, v := range list {
		if v.Id == element.Id {
			return k, nil
		}
	}

	return -1, errors.New("No matching element found!")
}
