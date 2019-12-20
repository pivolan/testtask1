package testtask1

type StateType string
type StatusType string
type SourceType string

const STATE_WIN StateType = "win"
const STATE_LOST StateType = "lost"

const STATUS_SUCCESS StatusType = "success"
const STATUS_FAIL StatusType = "fail"

const SOURCE_TYPE_GAME = "game"
const SOURCE_TYPE_SERVER = "server"
const SOURCE_TYPE_PAYMENT = "payment"
