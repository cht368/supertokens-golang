package smtpService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	evsmtpService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/smtpService"
	plesssmtpService "github.com/supertokens/supertokens-golang/recipe/passwordless/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSmtpService(config emaildelivery.SMTPTypeInput) emaildelivery.EmailDeliveryInterface {
	serviceImpl := makeServiceImplementation(config.SMTPSettings)

	if config.Override != nil {
		serviceImpl = config.Override(serviceImpl)
	}

	emailVerificationServiceImpl := evsmtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: config.SMTPSettings,
		Override:     makeEmailverificationServiceImplementation(serviceImpl),
	})
	passwordlessServiceImpl := plesssmtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: config.SMTPSettings,
		Override:     makePasswordlessServiceImplementation(serviceImpl),
	})

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			return (*emailVerificationServiceImpl.SendEmail)(input, userContext)

		} else if input.PasswordlessLogin != nil {
			return (*passwordlessServiceImpl.SendEmail)(input, userContext)

		} else {
			return errors.New("should never come here")
		}
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
