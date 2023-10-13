package centre

import (
	"kafka-iot-connect/client/mqtt"
	"kafka-iot-connect/errorx"
	"kafka-iot-connect/tool"

	"github.com/gofrs/uuid"
)

func (r *MysqlEkuiperCentre) GetClientByName(name string) (*MqttClient, error) {
	for _, ele := range r.Clients {
		if ele.Name == name {
			return ele, nil
		}
	}
	return nil, errorx.ErrMqttClientNotFound.New().ErrByName(name).ToError()
}

func (r *MysqlEkuiperCentre) GetClientById(id int) (*MqttClient, error) {
	for _, ele := range r.Clients {
		if ele.Id == id {
			return ele, nil
		}
	}
	return nil, errorx.ErrMqttClientNotFound.New().ErrById(id).ToError()
}

func (r *MysqlEkuiperCentre) GetClientIndexByName(name string) int {
	r.RLock()
	defer r.RUnlock()
	for ind, ele := range r.Clients {
		if ele.Name == name {
			return ind
		}
	}
	return -1
}

func (r *MysqlEkuiperCentre) UpdateClient(client *mqtt.MqttConfig, name string) (*MqttClient, error) {
	ind := r.GetClientIndexByName(name)
	if ind < 0 {
		return nil, errorx.ErrMqttClientNotFound.New().ErrByName(name).ToError()
	}
	cl := &MqttClient{
		Id:     len(r.Clients) + 1,
		Uuid:   uuid.Must(uuid.NewV4()),
		Name:   name,
		Client: client,
	}
	r.Lock()
	defer r.Unlock()
	r.Clients[ind].Client.Client.Disconnect(10)
	r.Clients[ind] = cl
	return cl, nil
}

func (r *MysqlEkuiperCentre) CreateClient(client *mqtt.MqttConfig, name string, id int) (*MqttClient, error) {
	if r == nil {
		return nil, errorx.ErrMysqlEkuiperCentreNotActivated
	}
	_, err := r.GetClientByName(name)
	if err == nil {
		return nil, errorx.ErrClientExists.New().ErrByName(name)
	}
	cl := &MqttClient{
		Id:     id,
		Uuid:   uuid.Must(uuid.NewV4()),
		Name:   name,
		Client: client,
	}
	r.Lock()
	defer r.Unlock()
	r.Clients = append(r.Clients, cl)
	return cl, nil
}

func (r *MysqlEkuiperCentre) RemoveClient(name string) error {
	ind := r.GetClientIndexByName(name)
	if ind < 0 {
		return errorx.ErrMysqlEkuiperCentreClientNotFound.New().ErrByName(name)
	}
	r.Lock()
	defer r.Unlock()
	r.Clients = tool.RemoveElementFromSlice(r.Clients, ind)
	return nil
}

func (r *MysqlEkuiperCentre) AllTopics() []string {
	ls := []string{}
	for _, ele := range r.Clients {
		for _, t := range ele.Client.AllCh.Chs {
			for _, top := range t.Topics {
				var exists bool
				for _, s := range ls {
					if s == top.Name {
						exists = true
						break
					}
				}
				if !exists {
					ls = append(ls, top.Name)
				}
			}
		}
	}
	return ls
}

func (r *MysqlEkuiperCentre) TopicContains(topic string) bool {
	for _, ele := range r.Clients {
		for _, t := range ele.Client.AllCh.Chs {
			for _, top := range t.Topics {
				if topic == top.Name {
					return true
				}
			}
		}
	}
	return false
}

func (r *MysqlEkuiperCentre) GetClientIndexById(id int) int {
	r.RLock()
	defer r.RUnlock()
	for ind, ele := range r.Clients {
		if ele.Id == id {
			return ind
		}
	}
	return -1
}

func (r *MysqlEkuiperCentre) RemoveClientById(id int) error {
	ind := r.GetClientIndexById(id)
	if ind < 0 {
		return errorx.ErrMysqlEkuiperCentreClientNotFound.New().ErrById(id)
	}
	r.Lock()
	defer r.Unlock()
	r.Clients = tool.RemoveElementFromSlice(r.Clients, ind)
	return nil
}

func (r *MysqlEkuiperCentre) DisconnectClientById(id int) error {
	ind := r.GetClientIndexById(id)
	if ind < 0 {
		return errorx.ErrMysqlEkuiperCentreClientNotFound.New().ErrById(id)
	}
	r.Lock()
	defer r.Unlock()
	r.Clients[ind].Client.Client.Disconnect(250)
	return nil
}
