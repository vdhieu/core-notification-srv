package webhook

type CreateWebHookReq struct {
	CallbackURL string `json:"callback_url"`
	Secret      string `json:"secret"`
}

type GetWebhookRes struct {
	ID          string `json:"id"`
	CallbackURL string `json:"callback_url"`
	CreatedAt   string `json:"created_at"`
	Secret      string `json:"secret"`
}

type UpdateWebHookReq struct {
	CallbackURL string `json:"callback_url"`
	Secret      string `json:"secret"`
}
