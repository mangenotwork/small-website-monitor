package global

import "time"

var UserToken = "sign"
var TokenExpires = 60 * 60 * 24 * 7
var LastSendMail int64 = 0
var TimeStamp = time.Now().Unix()

const DayLayout = "20060102"
