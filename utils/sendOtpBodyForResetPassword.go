package utils

import (
	"context"
	"careville_backend/database"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

var subject1 = aws.String("OTP to Reset Your Password")

func SendEmailForPassword(to string, link string) (*ses.SendEmailOutput, error) {
	sesClient := database.GetSesClient()

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				"careville@yopmail.com",
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data:    ResetPasswordOtpEmailBodyHtml(link),
					Charset: charSet,
				},
				Text: &types.Content{
					Data:    ResetPasswordOtpEmailBodyText(link),
					Charset: charSet,
				},
			},
			Subject: &types.Content{
				Data:    subject1,
				Charset: charSet,
			},
		},
		Source: sender,
	}

	return sesClient.SendEmail(context.Background(), input)
}
