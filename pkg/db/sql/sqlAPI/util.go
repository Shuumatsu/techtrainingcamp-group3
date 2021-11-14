package sqlAPI

import (
	"fmt"
	"strconv"
	"strings"
	"techtrainingcamp-group3/pkg/db/dbmodels"
)

func ParseEnvelopeList(envelopeList string) ([]dbmodels.EID, error) {
	envelopesID := make([]dbmodels.EID, 0)
	for _, token := range strings.Split(envelopeList, ",") {
		if len(token) == 0 {
			continue
		}
		eid, err := strconv.Atoi(token)
		if err != nil {
			return nil, fmt.Errorf("invaild: the envelope id can not change to number")
		}
		envelopesID = append(envelopesID, dbmodels.EID(eid))
	}
	return envelopesID, nil
}
