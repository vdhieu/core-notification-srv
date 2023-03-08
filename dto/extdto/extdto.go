package extdto

import "github.com/Neutronpay/lib-go-common/definitions/txndef"

type WebhookCallbackResp struct {
	TxnId     string          `json:"txnId"`
	ExtRefId  string          `json:"extRefId"`
	TxnState  txndef.TxnState `json:"txnState"`
	Msg       string          `json:"msg"`
	UpdatedAt int64           `json:"updatedAt"`
}

// place any EXTERNAL (NOT NEUTRONPAY, either downstream of upstream) dto objects used only by this service here
type WebhookReqBody struct {
	Data []WebhookReqData `json:"data"`
}

type WebhookReqData struct {
	Callback  string `json:"callback"`
	CreatedAt int64  `json:"createdAt"`
}
