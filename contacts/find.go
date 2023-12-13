package contacts

import "strings"

func FindMissing[O Contact, C Contact](oracle []O, check []C) []Contact {
	result := make([]Contact, 0, len(check))
	for _, checking := range check {
		found := false
		for _, contact := range oracle {
			if strings.EqualFold(checking.EmailAddress(), contact.EmailAddress()) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, checking)
		}
	}
	return result
}
