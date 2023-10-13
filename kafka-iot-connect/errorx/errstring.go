package errorx

import (
	"errors"
	"fmt"
)

type Errorx interface {
	ToError() error
	New() Errorx
	Error() string
	ErrByName(name string) Errorx
	ErrById(id int) Errorx
	ErrByClientId(id int) Errorx
	ErrByScheduleId(id int) Errorx
	ErrByBMSId(bms_id string) Errorx
	ErrByRuleId(ruleId int) Errorx
	ErrByFileId(ruleId int) Errorx
	ErrByField(field string) Errorx
	ErrByStreamId(streamId int) Errorx
	ErrBySinkId(sinkId int) Errorx
	ErrByCallbackId(callbackId int) Errorx
	ErrByTemplateId(templateId int) Errorx
	ErrByDataSourceId(templateId int) Errorx
	ErrByProjectId(projectId string) Errorx
	ErrByValue(value float64) Errorx
	ErrByType(ty string, value interface{}) Errorx
	ErrCombine(err error) Errorx
	Expect(value interface{}) Errorx
	Got(value interface{}) Errorx
}

type ErrorxString struct {
	s string
}

func New(text string) Errorx {
	return &ErrorxString{text}
}

func FromError(err error) Errorx {
	return &ErrorxString{err.Error()}
}

func CombineErrors(split string, err ...error) error {
	st := ""
	for _, ele := range err {
		st = st + ele.Error() + split
	}
	return fmt.Errorf(st)
}
func CombineErrorx(split string, err ...Errorx) Errorx {
	st := ""
	for _, ele := range err {
		st = st + ele.Error() + split
	}
	return &ErrorxString{st}
}

var (
	ErrByType = func(ty string, value interface{}) error {
		return fmt.Errorf(", %s: %v", ty, value)
	}
	ErrByValue = func(value float64) error {
		return fmt.Errorf(", value: %f", value)
	}
	ErrByName = func(name string) error {
		return fmt.Errorf(", name: %s", name)
	}
	ErrById = func(id int) error {
		return fmt.Errorf(", id: %d", id)
	}
	ErrByScheduleId = func(id int) error {
		return fmt.Errorf(", id: %d", id)
	}
	ErrByRuleId = func(ruleId int) error {
		return fmt.Errorf(", ruleId: %d", ruleId)
	}
	ErrByStreamId = func(streamId int) error {
		return fmt.Errorf(", streamId: %d", streamId)
	}
	ErrByClientId = func(clientId int) error {
		return fmt.Errorf(", clientId: %d", clientId)
	}
	ErrBySinkId = func(sinkId int) error {
		return fmt.Errorf(", sinkId: %d", sinkId)
	}
	ErrByCallbackId = func(callbackId int) error {
		return fmt.Errorf(", callbackId: %d", callbackId)
	}
	ErrByTemplateId = func(templateId int) error {
		return fmt.Errorf(", templateId: %d", templateId)
	}
	ErrByProjectId = func(projectId string) error {
		return fmt.Errorf(", projectId: %s", projectId)
	}
	ErrByDataSourceId = func(datasourceId int) error {
		return fmt.Errorf(", datasourceId: %d", datasourceId)
	}
	ErrByBmsId = func(bmsId string) error {
		return fmt.Errorf(", bmsId: %s", bmsId)
	}
	ErrByField = func(field string) error {
		return fmt.Errorf(", field: %s", field)
	}
	ErrByFileId = func(fileId int) error {
		return fmt.Errorf(", fileId: %d", fileId)
	}
	ErrCombine = func(err error) error {
		return fmt.Errorf(", error: %s", err.Error())
	}
	Expect = func(value interface{}) error {
		return fmt.Errorf(", expect: %v", value)
	}
	Got = func(value interface{}) error {
		return fmt.Errorf(", got: %v", value)
	}
)

func (e *ErrorxString) Error() string {
	return e.s
}

func (e *ErrorxString) ToError() error {
	return errors.New(e.s)
}

func (e *ErrorxString) New() Errorx {
	return &ErrorxString{e.s}
}

func (e *ErrorxString) ErrByName(name string) Errorx {
	e.s = e.s + ErrByName(name).Error()
	return e
}

func (e *ErrorxString) ErrByBMSId(bmsId string) Errorx {
	e.s = e.s + ErrByBmsId(bmsId).Error()
	return e
}

func (e *ErrorxString) ErrByField(field string) Errorx {
	e.s = e.s + ErrByField(field).Error()
	return e
}

func (e *ErrorxString) ErrByFileId(fileId int) Errorx {
	e.s = e.s + ErrByFileId(fileId).Error()
	return e
}

func (e *ErrorxString) ErrByScheduleId(scheduleId int) Errorx {
	e.s = e.s + ErrByScheduleId(scheduleId).Error()
	return e
}

func (e *ErrorxString) ErrByClientId(clientId int) Errorx {
	e.s = e.s + ErrByClientId(clientId).Error()
	return e
}

func (e *ErrorxString) ErrById(id int) Errorx {
	e.s = e.s + ErrById(id).Error()
	return e
}

func (e *ErrorxString) ErrByDataSourceId(id int) Errorx {
	e.s = e.s + ErrByDataSourceId(id).Error()
	return e
}

func (e *ErrorxString) ErrByProjectId(projectId string) Errorx {
	e.s = e.s + ErrByProjectId(projectId).Error()
	return e
}

func (e *ErrorxString) ErrByRuleId(ruleId int) Errorx {
	e.s = e.s + ErrByRuleId(ruleId).Error()
	return e
}

func (e *ErrorxString) ErrByStreamId(streamId int) Errorx {
	e.s = e.s + ErrByStreamId(streamId).Error()
	return e
}

func (e *ErrorxString) ErrBySinkId(sinkId int) Errorx {
	e.s = e.s + ErrBySinkId(sinkId).Error()
	return e
}

func (e *ErrorxString) ErrByCallbackId(callbackId int) Errorx {
	e.s = e.s + ErrByCallbackId(callbackId).Error()
	return e
}

func (e *ErrorxString) ErrByTemplateId(templateId int) Errorx {
	e.s = e.s + ErrByTemplateId(templateId).Error()
	return e
}

func (e *ErrorxString) ErrByType(ty string, value interface{}) Errorx {
	e.s = e.s + ErrByType(ty, value).Error()
	return e
}

func (e *ErrorxString) ErrByValue(value float64) Errorx {
	e.s = e.s + ErrByValue(value).Error()
	return e
}

func (e *ErrorxString) ErrCombine(err error) Errorx {
	e.s = e.s + ErrCombine(err).Error()
	return e
}

func (e *ErrorxString) Expect(value interface{}) Errorx {
	e.s = e.s + Expect(value).Error()
	return e
}

func (e *ErrorxString) Got(value interface{}) Errorx {
	e.s = e.s + Got(value).Error()
	return e
}
