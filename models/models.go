package models

import (
	"encoding/json"
)

type ErrorCode uint64

type UID uint64

func (id UID) String() string {
	return int2str(uint64(id))
}

type EID uint64

func (id EID) String() string {
	return int2str(uint64(id))
}

type SnatchReq struct {
	Uid UID `json:"uid"`
}
type SnatchData struct {
	EnvelopeId EID    `json:"envelope_id"`
	MaxCount   uint64 `json:"max_count"`
	CurCount   uint64 `json:"cur_count"`
}
type SnatchResp struct {
	Code ErrorCode  `json:"code"`
	Msg  string     `json:"msg"`
	Data SnatchData `json:"data"`
}

type OpenReq struct {
	Uid        UID    `json:"uid"`
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
	EnvelopeId EID    `json:"envelope_id" bson:"envelope_id"`
	Opened     bool   `json:"opened" bson:"opened"`
	Value      uint64 `json:"value,omitempty" bson:"value"`
	SnatchTime int64  `json:"snatch_time" bson:"snatch_time"`
}

type WalletListReq struct {
	Uid UID `json:"uid"`
}
type WalletListData struct {
	Amount       uint64     `json:"amount" bson:"amount"`
	EnvelopeList []Envelope `json:"envelope_list" bson:"envelope_list"`
}
type WalletListResp struct {
	Code ErrorCode      `json:"code"`
	Msg  string         `json:"msg"`
	Data WalletListData `json:"data"`
}

type User struct {
	Uid    UID            `json:"uid" bson:"uid"`
	Wallet WalletListData `json:"wallet" bson:"wallet"`
}

func (u *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (User) CollectionName() string {
	return "user"
}

func (w *WalletListData) Size() int {
	return len(w.EnvelopeList)
}

func int2str(num uint64) string {
	if num == 0 {
		return "0"
	}
	var ret []byte
	for num != 0 {
		ret = append(ret, byte(num%10)+'0')
		num /= 10
	}
	for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = ret[j], ret[i]
	}
	return string(ret)
}
