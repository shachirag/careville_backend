package utils

import (
	"careville_backend/database"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

var (
	charSet = aws.String("UTF-8")
	sender  = aws.String("careville@yopmail.com")
	subject = aws.String("otp for user signup")
)

func SendEmail(to string, link string) (*ses.SendEmailOutput, error) {
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
					Data:    SignupUserOtpEmailBodyHtml(link),
					Charset: charSet,
				},
				Text: &types.Content{
					Data:    SignupUserOtpEmailBodyText(link),
					Charset: charSet,
				},
			},
			Subject: &types.Content{
				Data:    subject,
				Charset: charSet,
			},
		},
		Source: sender,
	}

	return sesClient.SendEmail(context.Background(), input)
}
