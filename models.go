package main

import "time"

type ErrorCode uint64

type SnaptchReq struct {
	Uid uint64 `json:uid`
}
type SnaptchData struct {
	EnvelopeId uint64 `json:envelope_id`
	MaxCount   uint64 `json:max_count`
	CurCount   uint64 `json:cur_count`
}
type SnaptchResp struct {
	Code ErrorCode   `json:code`
	Msg  string      `json:msg`
	Data SnaptchData `json:data`
}

type OpenReq struct {
	Uid        uint64 `json:uid`
	EnvelopeId uint64 `json:envelope_id`
}
type OpenData struct {
	Value uint64 `json:envelope_id`
}
type OpenResp struct {
	Code ErrorCode `json:code`
	Msg  string    `json:msg`
	Data OpenData  `json:data`
}

type Envelope struct {
	EnvelopeId  uint64    `json:envelope_id`
	Value       uint64    `json:value`
	Opened      bool      `json:opened`
	SnaptchTime time.Time `json:snatch_time`
}

type WalletListReq struct {
	Uid uint64 `json:uid`
}
type WalletListData struct {
	Amount       uint64     `json:amount`
	EnvelopeList []Envelope `json:envelope_list`
}
type WalletListResp struct {
	Code ErrorCode      `json:code`
	Msg  string         `json:msg`
	Data WalletListData `json:data`
}
