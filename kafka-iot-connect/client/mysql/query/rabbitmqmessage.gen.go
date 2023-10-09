// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"kafka-iot-connect/client/mysql/model"
)

func newRabbitMQMessage(db *gorm.DB, opts ...gen.DOOption) rabbitMQMessage {
	_rabbitMQMessage := rabbitMQMessage{}

	_rabbitMQMessage.rabbitMQMessageDo.UseDB(db, opts...)
	_rabbitMQMessage.rabbitMQMessageDo.UseModel(&model.RabbitMQMessage{})

	tableName := _rabbitMQMessage.rabbitMQMessageDo.TableName()
	_rabbitMQMessage.ALL = field.NewAsterisk(tableName)
	_rabbitMQMessage.ID = field.NewInt32(tableName, "id")
	_rabbitMQMessage.Name = field.NewString(tableName, "name")
	_rabbitMQMessage.Type = field.NewString(tableName, "type")
	_rabbitMQMessage.ProjectID = field.NewString(tableName, "projectId")
	_rabbitMQMessage.Status = field.NewInt32(tableName, "status")
	_rabbitMQMessage.ClientID = field.NewInt32(tableName, "clientId")
	_rabbitMQMessage.StreamID = field.NewInt32(tableName, "streamId")
	_rabbitMQMessage.SinkID = field.NewInt32(tableName, "sinkId")
	_rabbitMQMessage.CallbackID = field.NewInt32(tableName, "callbackId")
	_rabbitMQMessage.RClientID = field.NewInt32(tableName, "rClientId")
	_rabbitMQMessage.QueueID = field.NewInt32(tableName, "queueId")
	_rabbitMQMessage.ExchangeID = field.NewInt32(tableName, "exchangeId")
	_rabbitMQMessage.RoutingKey = field.NewString(tableName, "routingKey")
	_rabbitMQMessage.Args = field.NewString(tableName, "args")
	_rabbitMQMessage.Source = field.NewString(tableName, "source")
	_rabbitMQMessage.CreatedTime = field.NewTime(tableName, "createdTime")
	_rabbitMQMessage.UpdatedTime = field.NewTime(tableName, "updatedTime")

	_rabbitMQMessage.fillFieldMap()

	return _rabbitMQMessage
}

type rabbitMQMessage struct {
	rabbitMQMessageDo

	ALL         field.Asterisk
	ID          field.Int32
	Name        field.String
	Type        field.String
	ProjectID   field.String
	Status      field.Int32
	ClientID    field.Int32
	StreamID    field.Int32
	SinkID      field.Int32
	CallbackID  field.Int32
	RClientID   field.Int32
	QueueID     field.Int32
	ExchangeID  field.Int32
	RoutingKey  field.String
	Args        field.String
	Source      field.String
	CreatedTime field.Time
	UpdatedTime field.Time

	fieldMap map[string]field.Expr
}

func (r rabbitMQMessage) Table(newTableName string) *rabbitMQMessage {
	r.rabbitMQMessageDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r rabbitMQMessage) As(alias string) *rabbitMQMessage {
	r.rabbitMQMessageDo.DO = *(r.rabbitMQMessageDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *rabbitMQMessage) updateTableName(table string) *rabbitMQMessage {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewInt32(table, "id")
	r.Name = field.NewString(table, "name")
	r.Type = field.NewString(table, "type")
	r.ProjectID = field.NewString(table, "projectId")
	r.Status = field.NewInt32(table, "status")
	r.ClientID = field.NewInt32(table, "clientId")
	r.StreamID = field.NewInt32(table, "streamId")
	r.SinkID = field.NewInt32(table, "sinkId")
	r.CallbackID = field.NewInt32(table, "callbackId")
	r.RClientID = field.NewInt32(table, "rClientId")
	r.QueueID = field.NewInt32(table, "queueId")
	r.ExchangeID = field.NewInt32(table, "exchangeId")
	r.RoutingKey = field.NewString(table, "routingKey")
	r.Args = field.NewString(table, "args")
	r.Source = field.NewString(table, "source")
	r.CreatedTime = field.NewTime(table, "createdTime")
	r.UpdatedTime = field.NewTime(table, "updatedTime")

	r.fillFieldMap()

	return r
}

func (r *rabbitMQMessage) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *rabbitMQMessage) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 17)
	r.fieldMap["id"] = r.ID
	r.fieldMap["name"] = r.Name
	r.fieldMap["type"] = r.Type
	r.fieldMap["projectId"] = r.ProjectID
	r.fieldMap["status"] = r.Status
	r.fieldMap["clientId"] = r.ClientID
	r.fieldMap["streamId"] = r.StreamID
	r.fieldMap["sinkId"] = r.SinkID
	r.fieldMap["callbackId"] = r.CallbackID
	r.fieldMap["rClientId"] = r.RClientID
	r.fieldMap["queueId"] = r.QueueID
	r.fieldMap["exchangeId"] = r.ExchangeID
	r.fieldMap["routingKey"] = r.RoutingKey
	r.fieldMap["args"] = r.Args
	r.fieldMap["source"] = r.Source
	r.fieldMap["createdTime"] = r.CreatedTime
	r.fieldMap["updatedTime"] = r.UpdatedTime
}

func (r rabbitMQMessage) clone(db *gorm.DB) rabbitMQMessage {
	r.rabbitMQMessageDo.ReplaceConnPool(db.Statement.ConnPool)
	return r
}

func (r rabbitMQMessage) replaceDB(db *gorm.DB) rabbitMQMessage {
	r.rabbitMQMessageDo.ReplaceDB(db)
	return r
}

type rabbitMQMessageDo struct{ gen.DO }

type IRabbitMQMessageDo interface {
	gen.SubQuery
	Debug() IRabbitMQMessageDo
	WithContext(ctx context.Context) IRabbitMQMessageDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IRabbitMQMessageDo
	WriteDB() IRabbitMQMessageDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IRabbitMQMessageDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IRabbitMQMessageDo
	Not(conds ...gen.Condition) IRabbitMQMessageDo
	Or(conds ...gen.Condition) IRabbitMQMessageDo
	Select(conds ...field.Expr) IRabbitMQMessageDo
	Where(conds ...gen.Condition) IRabbitMQMessageDo
	Order(conds ...field.Expr) IRabbitMQMessageDo
	Distinct(cols ...field.Expr) IRabbitMQMessageDo
	Omit(cols ...field.Expr) IRabbitMQMessageDo
	Join(table schema.Tabler, on ...field.Expr) IRabbitMQMessageDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IRabbitMQMessageDo
	RightJoin(table schema.Tabler, on ...field.Expr) IRabbitMQMessageDo
	Group(cols ...field.Expr) IRabbitMQMessageDo
	Having(conds ...gen.Condition) IRabbitMQMessageDo
	Limit(limit int) IRabbitMQMessageDo
	Offset(offset int) IRabbitMQMessageDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IRabbitMQMessageDo
	Unscoped() IRabbitMQMessageDo
	Create(values ...*model.RabbitMQMessage) error
	CreateInBatches(values []*model.RabbitMQMessage, batchSize int) error
	Save(values ...*model.RabbitMQMessage) error
	First() (*model.RabbitMQMessage, error)
	Take() (*model.RabbitMQMessage, error)
	Last() (*model.RabbitMQMessage, error)
	Find() ([]*model.RabbitMQMessage, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RabbitMQMessage, err error)
	FindInBatches(result *[]*model.RabbitMQMessage, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.RabbitMQMessage) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IRabbitMQMessageDo
	Assign(attrs ...field.AssignExpr) IRabbitMQMessageDo
	Joins(fields ...field.RelationField) IRabbitMQMessageDo
	Preload(fields ...field.RelationField) IRabbitMQMessageDo
	FirstOrInit() (*model.RabbitMQMessage, error)
	FirstOrCreate() (*model.RabbitMQMessage, error)
	FindByPage(offset int, limit int) (result []*model.RabbitMQMessage, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IRabbitMQMessageDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	FilterWithStatus(status int) (result []model.RabbitMQMessage, err error)
}

// SELECT * FROM @@table WHERE status = @status{{end}}
func (r rabbitMQMessageDo) FilterWithStatus(status int) (result []model.RabbitMQMessage, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, status)
	generateSQL.WriteString("SELECT * FROM rabbitMQMessage WHERE status = ? ")

	var executeSQL *gorm.DB
	executeSQL = r.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (r rabbitMQMessageDo) Debug() IRabbitMQMessageDo {
	return r.withDO(r.DO.Debug())
}

func (r rabbitMQMessageDo) WithContext(ctx context.Context) IRabbitMQMessageDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r rabbitMQMessageDo) ReadDB() IRabbitMQMessageDo {
	return r.Clauses(dbresolver.Read)
}

func (r rabbitMQMessageDo) WriteDB() IRabbitMQMessageDo {
	return r.Clauses(dbresolver.Write)
}

func (r rabbitMQMessageDo) Session(config *gorm.Session) IRabbitMQMessageDo {
	return r.withDO(r.DO.Session(config))
}

func (r rabbitMQMessageDo) Clauses(conds ...clause.Expression) IRabbitMQMessageDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r rabbitMQMessageDo) Returning(value interface{}, columns ...string) IRabbitMQMessageDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r rabbitMQMessageDo) Not(conds ...gen.Condition) IRabbitMQMessageDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r rabbitMQMessageDo) Or(conds ...gen.Condition) IRabbitMQMessageDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r rabbitMQMessageDo) Select(conds ...field.Expr) IRabbitMQMessageDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r rabbitMQMessageDo) Where(conds ...gen.Condition) IRabbitMQMessageDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r rabbitMQMessageDo) Order(conds ...field.Expr) IRabbitMQMessageDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r rabbitMQMessageDo) Distinct(cols ...field.Expr) IRabbitMQMessageDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r rabbitMQMessageDo) Omit(cols ...field.Expr) IRabbitMQMessageDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r rabbitMQMessageDo) Join(table schema.Tabler, on ...field.Expr) IRabbitMQMessageDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r rabbitMQMessageDo) LeftJoin(table schema.Tabler, on ...field.Expr) IRabbitMQMessageDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r rabbitMQMessageDo) RightJoin(table schema.Tabler, on ...field.Expr) IRabbitMQMessageDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r rabbitMQMessageDo) Group(cols ...field.Expr) IRabbitMQMessageDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r rabbitMQMessageDo) Having(conds ...gen.Condition) IRabbitMQMessageDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r rabbitMQMessageDo) Limit(limit int) IRabbitMQMessageDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r rabbitMQMessageDo) Offset(offset int) IRabbitMQMessageDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r rabbitMQMessageDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IRabbitMQMessageDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r rabbitMQMessageDo) Unscoped() IRabbitMQMessageDo {
	return r.withDO(r.DO.Unscoped())
}

func (r rabbitMQMessageDo) Create(values ...*model.RabbitMQMessage) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r rabbitMQMessageDo) CreateInBatches(values []*model.RabbitMQMessage, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r rabbitMQMessageDo) Save(values ...*model.RabbitMQMessage) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r rabbitMQMessageDo) First() (*model.RabbitMQMessage, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.RabbitMQMessage), nil
	}
}

func (r rabbitMQMessageDo) Take() (*model.RabbitMQMessage, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.RabbitMQMessage), nil
	}
}

func (r rabbitMQMessageDo) Last() (*model.RabbitMQMessage, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.RabbitMQMessage), nil
	}
}

func (r rabbitMQMessageDo) Find() ([]*model.RabbitMQMessage, error) {
	result, err := r.DO.Find()
	return result.([]*model.RabbitMQMessage), err
}

func (r rabbitMQMessageDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RabbitMQMessage, err error) {
	buf := make([]*model.RabbitMQMessage, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r rabbitMQMessageDo) FindInBatches(result *[]*model.RabbitMQMessage, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r rabbitMQMessageDo) Attrs(attrs ...field.AssignExpr) IRabbitMQMessageDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r rabbitMQMessageDo) Assign(attrs ...field.AssignExpr) IRabbitMQMessageDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r rabbitMQMessageDo) Joins(fields ...field.RelationField) IRabbitMQMessageDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r rabbitMQMessageDo) Preload(fields ...field.RelationField) IRabbitMQMessageDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r rabbitMQMessageDo) FirstOrInit() (*model.RabbitMQMessage, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.RabbitMQMessage), nil
	}
}

func (r rabbitMQMessageDo) FirstOrCreate() (*model.RabbitMQMessage, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.RabbitMQMessage), nil
	}
}

func (r rabbitMQMessageDo) FindByPage(offset int, limit int) (result []*model.RabbitMQMessage, count int64, err error) {
	result, err = r.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = r.Offset(-1).Limit(-1).Count()
	return
}

func (r rabbitMQMessageDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r rabbitMQMessageDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r rabbitMQMessageDo) Delete(models ...*model.RabbitMQMessage) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *rabbitMQMessageDo) withDO(do gen.Dao) *rabbitMQMessageDo {
	r.DO = *do.(*gen.DO)
	return r
}
