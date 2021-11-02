package handler

import (
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"

	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"strings"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/db"
	"techtrainingcamp-group3/models"
	"techtrainingcamp-group3/tools"
	"time"
)

func SnatchHandler(c *gin.Context) {
	var req models.SnatchReq
	c.Bind(&req)

	if rand.Float32() > config.SnatchProb {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "fail",
		})
	} else {
		var user db.User
		max_count := config.MaxSnatchAmount
		cur_count := 0
		if err := db.DB.Table(db.User{}.TableName()).First(
			&user, req.Uid).Error; err == nil {
			cur_count = strings.Count(user.EnvelopeList, ",") + 1
		}
		if cur_count >= max_count {
			c.JSON(200, gin.H{
				"code": 2,
				"msg":  "too many envelopes",
			})
		} else {
			envelope := tools.REPool.Snatch()
			eid := strconv.FormatUint(envelope.Eid, 10)
			if cur_count == 0 {
				db.DB.Table(db.User{}.TableName()).Create(db.User{Uid: req.Uid, EnvelopeList: eid, Amount: 0})
			} else {
				db.DB.Table(db.User{}.TableName()).Where("uid", req.Uid).Update("envelope_list", user.EnvelopeList+","+eid)
			}
			db.DB.Table(db.Envelope{}.TableName()).Create(db.Envelope{EnvelopeId: envelope.Eid, Opened: false, Value: envelope.Money, SnatchTime: uint64(time.Now().UTC().UnixNano())})
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "success",
				"data": gin.H{
					"envelope_id": envelope.Eid,
					"max_count":   max_count,
					"cur_count":   cur_count + 1,
				},
			})
		}
	}
}
