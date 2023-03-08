package entity

// put any entity internal to this service here, esp persistable ones with gorm capabilities

type WebhookRecord interface {
	AccountId() string
	UrlStr() string
	Secret() string
}

type webhookRecord struct {
	accountId string
	urlStr    string
	secret    string
}

func (w *webhookRecord) AccountId() string {
	return w.accountId
}

func (w *webhookRecord) UrlStr() string {
	return w.urlStr
}

func (w *webhookRecord) Secret() string {
	return w.secret
}

// GetWebhookRecordForAccount returns nil, nil if no webhook found, since it's not an error state if so
func GetWebhookRecordForAccount(accountId string) (record WebhookRecord, err error) {

	record = &webhookRecord{
		accountId: accountId,
		urlStr:    "http://192.168.3.145:8234/callback/platform/",
		secret:    "somerandomtsignaturestring",
	}

	return
}

func UpdateWebhookRecord(record WebhookRecord) (err error) {

	return
}

func DeleteWebhookRecord(accountId string, record WebhookRecord) (err error) {

	return
}
