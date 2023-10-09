package mqtt

import (
	"fmt"
	"kafka-iot-connect/tool"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (a *AllMessageChannels) FindChannelById(id string) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for ind, ele := range a.Chs {
		if ele.Id == id {
			return ind
		}
	}
	return -1
}

func (a *AllMessageChannels) AddChannel(mch MessageChannel) int {
	if ind := a.FindChannelById(mch.Id); ind > -1 {
		return -1
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Chs = append(a.Chs, mch)
	return -1
}

func (a *AllMessageChannels) UnsubscribeUnuseTopic(topic string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var stillInUse bool
	for _, ele := range a.Chs {
		for _, el := range ele.Topics {
			if topic == el.Name {
				stillInUse = true
			}
		}
	}
	return !stillInUse
}

func (a *AllMessageChannels) RemoveChannelById(id string) error {
	if ind := a.FindChannelById(id); ind > -1 {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.Chs = tool.RemoveElementFromSlice(a.Chs, ind)
		return nil
	} else {
		return fmt.Errorf("cannot find channel by id: %s", id)
	}
}

func (a *AllMessageChannels) RemoveAllChannels() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Chs = make([]MessageChannel, 0)
}

func (a *AllMessageChannels) FindChannelsByTopic(topic string) []int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	ls := []int{}
	for ind, ele := range a.Chs {
		for _, top := range ele.Topics {
			if top.Name == topic {
				ls = append(ls, ind)
			}
		}
	}
	return ls
}

func (a *AllMessageChannels) UpdateCallbackById(id string, callback func(mqtt.Client, [][]byte, []Topic, int)) error {
	ind := a.FindChannelById(id)
	if ind < 0 {
		return fmt.Errorf("cannot find channel by id: %s", id)
	} else {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.Chs[ind].CallbackAction = callback
		return nil
	}
}
