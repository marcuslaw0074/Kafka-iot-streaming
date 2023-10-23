package mqtt

import (
	"fmt"
	logging "kafka-iot-connect/log"

	"time"
)

func (m *MqttConfig) OnKaflaChannelHandler(channelId string, props interface{}) {
	go func(m *MqttConfig, chId string, i interface{}) {
		timer := time.NewTimer(1 * time.Millisecond)
		cid := m.AllCh.FindChannelById(channelId)
		m.Log.Log(logging.LogInfo, pkgName, fmt.Sprintf("finding message channel id: %s ", channelId))
		for cid > -1 {
			m.AllCh.Chs[cid].IsActive = true
			select {
			case val := <-m.AllCh.Chs[cid].Ch:
				m.AllCh.Chs[cid].Data = append(m.AllCh.Chs[cid].Data, val)
				timer.Reset(m.AllCh.Chs[cid].WaitingTime)
			case <-timer.C:
				for len(m.AllCh.Chs[cid].Data) > 0 {
					m.Log.Log(logging.LogInfo, pkgName, fmt.Sprintf("message length: %d", len(m.AllCh.Chs[cid].Data)))
					go m.AllCh.Chs[cid].CallbackAction(m.Client, m.AllCh.Chs[cid].Data, m.AllCh.Chs[cid].Topics, 0)
					m.AllCh.Chs[cid].Data = [][]byte{}
				}
			}
			cid = m.AllCh.FindChannelById(channelId)
		}
		m.AllCh.Chs[cid].IsActive = false
		m.Log.Log(logging.LogInfo, pkgName, fmt.Sprintf("message channel id: %s ends!", channelId))
	}(m, channelId, props)
}
