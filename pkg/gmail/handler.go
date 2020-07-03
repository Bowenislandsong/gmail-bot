package gmail

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/gmail/v1"
	goption "google.golang.org/api/option"
)

var gAuthFilePath = os.Getenv("GOOGLE_CREDS")

type EmailContent struct {
	fromAddr string
	toAddr   []string
	ccAddr   []string
	title    string
	msg      string
	service  *gmail.Service
}

func (e *EmailContent) apply(options []EmailOption) {
	for _, option := range options {
		option(e)
	}
}

func (e *EmailContent) complete() error {
	if e.fromAddr == "" {
		return fmt.Errorf("from address can not be empty")
	}

	if e.toAddr == nil {
		return fmt.Errorf("to address can not be empty")
	}

	if e.msg == "" {
		return fmt.Errorf("message can not be empty")
	}
	return nil
}

type EmailOption func(config *EmailContent)

func From(from string) EmailOption {
	return func(config *EmailContent) {
		config.fromAddr = from
	}
}

func To(to ...string) EmailOption {
	return func(config *EmailContent) {
		config.toAddr = to
	}
}

func Cc(cc ...string) EmailOption {
	return func(config *EmailContent) {
		config.ccAddr = cc
	}
}
func Title(title string) EmailOption {
	return func(config *EmailContent) {
		config.title = title
	}
}

func MessageBody(msg string) EmailOption {
	return func(config *EmailContent) {
		config.msg = msg
	}
}

func NewEmail() (*EmailContent, error) {
	ctx := context.Background()
	service, err := gmail.NewService(ctx, goption.WithCredentialsFile(gAuthFilePath), goption.WithScopes(gmail.MailGoogleComScope))
	if err != nil {
		return nil, err
	}

	return &EmailContent{
		service: service,
	}, nil
}

func (e *EmailContent) Send(option ...EmailOption) (*gmail.Message, error) {
	e.apply(option)
	if err := e.complete(); err != nil {
		return nil, err
	}

	var headers []*gmail.MessagePartHeader
	for _, to := range e.toAddr {
		headers = append(headers, &gmail.MessagePartHeader{
			Name:  "<code>To</code>",
			Value: "<code>" + to + "</code>",
		})
	}

	for _, cc := range e.ccAddr {
		headers = append(headers, &gmail.MessagePartHeader{
			Name:  "<code>cc</code>",
			Value: "<code>" + cc + "</code>",
		})
	}

	headers = append(headers, &gmail.MessagePartHeader{
		Name:  "<code>From</code>",
		Value: "<code>" + e.fromAddr + "</code>",
	})

	headers = append(headers, &gmail.MessagePartHeader{
		Name:  "<code>Subject</code>",
		Value: "<code>" + e.title + "</code>",
	})

	headers = append(headers, &gmail.MessagePartHeader{
		Name:  "<code>Body</code>",
		Value: "<code>" + e.msg + "</code>",
	})

	msg := gmail.Message{
		Payload: &gmail.MessagePart{
			//Body: &gmail.MessagePartBody{
			//	Data: e.msg,
			//},
			Headers: headers,
		},
	}

	return e.service.Users.Messages.Send("", &msg).Do()
}
