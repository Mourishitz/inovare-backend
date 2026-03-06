package services

import "log"

type EmailService interface {
	SendCatalogChangesNotification(hostEmail, hostName string, catalogID uint) error
}

type emailService struct{}

func NewEmailService() EmailService {
	return &emailService{}
}

// SendCatalogChangesNotification notifies the host that their catalog has been updated.
// TODO: replace this stub with a real SMTP/SES implementation.
func (s *emailService) SendCatalogChangesNotification(hostEmail, hostName string, catalogID uint) error {
	log.Printf("[email] catalog %d changes notification → %s <%s>", catalogID, hostName, hostEmail)
	return nil
}
