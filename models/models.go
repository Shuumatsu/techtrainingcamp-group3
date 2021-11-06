package models

type ErrorCode uint64

const (
	Success ErrorCode = iota
	NotFound
	ParseError
	EnvelopeAlreadyOpen
	ErrorEnvelopeOwner
	DataBaseError
	SnatchLimit
	NotDefined
)

func (e ErrorCode) Message() string {
	switch e {
	case Success:
		return "success"
	case NotFound:
		return "not found"
	case ParseError:
		return "the envelope list parse error"
	case EnvelopeAlreadyOpen:
		return "the envelope already open"
	case ErrorEnvelopeOwner:
		return "the envelope not belong to this user"
	case DataBaseError:
		return "the database error"
	case SnatchLimit:
		return "snatch limit"
	default:
		return "not defined"
	}
}

type SnatchReq struct {
	Uid uint64 `json:"uid"`
}
type SnatchData struct {
	EnvelopeId uint64 `json:"envelope_id"`
	MaxCount   uint64 `json:"max_count"`
	CurCount   uint64 `json:"cur_count"`
}
type SnatchResp struct {
	Code ErrorCode  `json:"code"`
	Msg  string     `json:"msg"`
	Data SnatchData `json:"data"`
}

type OpenReq struct {
	Uid        uint64 `json:"uid"`
	EnvelopeId uint64 `json:"envelope_id"`
}
type OpenData struct {
	Value uint64 `json:"envelope_id"`
}
type OpenResp struct {
	Code ErrorCode `json:"code"`
	Msg  string    `json:"msg"`
	Data OpenData  `json:"data"`
}

type Envelope struct {
	EnvelopeId uint64 `json:"envelope_id"`
	Opened     bool   `json:"opened"`
	Value      uint64 `json:"value,omitempty"`
	SnatchTime int64  `json:"snatch_time"`
}

type WalletListReq struct {
	Uid uint64 `json:"uid"`
}
type WalletListData struct {
	Amount       uint64     `json:"amount"`
	EnvelopeList []Envelope `json:"envelope_list"`
}
type WalletListResp struct {
	Code ErrorCode      `json:"code"`
	Msg  string         `json:"msg"`
	Data WalletListData `json:"data"`
}

func (w *WalletListData) Size() int {
	return len(w.EnvelopeList)
}
