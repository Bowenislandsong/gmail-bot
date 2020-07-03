package gmail

import (
	"fmt"
	"strings"
)

func ParseEmailAddress(emailAddr string) ([]string, error) {
	if emailAddr == "" || strings.Count(emailAddr, "@") == 0 {
		return nil, fmt.Errorf("no email address found in %s", emailAddr)
	}

	var emails []string
	for _, e := range strings.Split(emailAddr, ",") {
		if strings.Count(e, "@") == 1 {
			emails = append(emails, e)
		}
	}

	return emails, nil
}
