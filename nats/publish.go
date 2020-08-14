package nats

import "drops-backend/models"

func PublishToken(m *models.MailInfo) {
	Nats.Publish("mail.token", m)
}
