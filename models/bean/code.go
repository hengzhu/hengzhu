package bean

const (
	CODE_Success               = 200
	CODE_Created               = 201
	CODE_Bad_Request           = 400 // 请求错误
	CODE_Unauthorized          = 401 // 没有登录
	CODE_Not_Found             = 404 // not found
	CODE_Forbidden             = 403 // 没有权限
	CODE_Method_Not_Allowed    = 405 // 方法不对 (POST,PUT,GET)
	CODE_Not_Acceptable        = 406 // 不能通过
	CODE_Internal_Server_Error = 500 // 服务错误

	CODE_Params_Err = 430 // 参数错误
)

func CodeString(code int) string {
	s := map[int]string{
		CODE_Success:               "OK",
		CODE_Created:               "Created",
		CODE_Bad_Request:           "Bad_Request",
		CODE_Unauthorized:          "Unauthorized",
		CODE_Not_Found:             "Not_Found",
		CODE_Forbidden:             "Forbidden",
		CODE_Method_Not_Allowed:    "Method_Not_Allowed",
		CODE_Not_Acceptable:        "Not_Acceptable",
		CODE_Internal_Server_Error: "Server_Error",
		CODE_Params_Err:            "Params_Error",
	}[code]
	return s
}

type RabbitMqMessage struct {
	CabinetId string `json:"cabinet_id"`
	Door      int    `json:"door"`
	Heartbeat int    `json:"heartbeat"`
	UserId    string `json:"user_id"`
	DoorState string `json:"door_state"`
}

type CabinetInfo struct {
	CabinetID string      `json:"cabinet_id" description:"柜子id"`
	Number    string      `json:"number" description:"编号"`
	Desc      string      `json:"desc" description:"备注"`
	Doors     int         `json:"doors" description:"门数"`
	DoorStat  []DoorState `json:"door_stat" descrpition:"门状态"`
}

type DoorState struct {
	Door          int  `json:"door"`
	Locked        bool `json:"locked"`
	WireConnected bool `json:"wire_connected"`
}
