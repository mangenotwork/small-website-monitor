package global

import "time"

var UserToken = "sign"
var TokenExpires = 60 * 60 * 24 * 7
var LastSendMail int64 = 0
var TimeStamp = time.Now().Unix()

const DayLayout = "20060102"

const Version = "v0.1"
const HttpSubassembly = "gathertool(https://github.com/mangenotwork/gathertool)"
