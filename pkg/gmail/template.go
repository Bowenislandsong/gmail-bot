package gmail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

const EmailKind = "email"

type EmailTemplate struct {
	Kind       string `json:"kind"`
	Cron    string `json:"cron"`
	To         string `json:"to"`
	Cc         string `json:"cc,omitempty"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Auth       string `json:"auth,omitempty"`
}

func LoadEmailTemplate(filePath string) (*EmailTemplate, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var email EmailTemplate
	if err := json.Unmarshal(content, &email); err != nil {
		return nil, err
	}

	//TODO: Load template from gsheets

	if err := email.validateEmailTemplate(); err != nil {
		return nil, fmt.Errorf("failed to pase %s as email template, %v", filePath, err)
	}

	return &email, nil
}

func (e *EmailTemplate) validateEmailTemplate() error {
	if e.Kind != EmailKind {
		return fmt.Errorf("template is not an email template")
	}

	schedule, err := cron.ParseStandard(e.Cron)
	if err != nil {
		return err
	}
	log.Infof("The next activation for %s is %v", e.Title, schedule.Next(time.Now()))

	if _, err := ParseEmailAddress(e.To); err != nil {
		return err
	}

	if e.Title != "" && e.Body != "" {
		return fmt.Errorf("email requires nonempty title and body")
	}

	return nil
}
