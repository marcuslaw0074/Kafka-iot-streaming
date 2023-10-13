package errorx

var (
	ErrMqttClientNotFound Errorx = New("cannot find any mqtt clients")
	ErrMysqlEkuiperCentreNotActivated Errorx = New("mysqlEkuiperCentre is not activated yet")
	ErrClientExists          Errorx = New("client already exists")
	ErrMysqlEkuiperCentreClientNotFound    Errorx = New("cannot find any client from mysqlEkuiperCentre")

)
