package tempno

import (
	"bufio"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
)

type TempNo struct {
	blacklist map[string]struct{}
}

func Load(mailsfile string) (*TempNo, error) {
	tempno := TempNo{}
	tempno.blacklist = make(map[string]struct{})
	f, err := os.Open(mailsfile)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tempno.blacklist[strings.ToLower(scanner.Text())] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &tempno, nil

}

func (mb *TempNo) IsValid(email string) bool {
	email = strings.ToLower(email)
	if !govalidator.IsEmail(email) {
		return false
	}
	// First get the domain, then split the domain into parts
	domainComponents := strings.Split(
		strings.Split(email, "@")[1],
		".",
	)
	// Iterate over the domain components
	for i := range domainComponents[:len(domainComponents)-1] {
		// Join the domain's components
		structuredDomain := strings.Join(domainComponents[i:], ".")
		// Check for the existance of the created domain inside the blacklist
		if _, isInvalid := mb.blacklist[structuredDomain]; isInvalid {
			// Invalid - The email domain is a bad domain
			return false
		}
	}
	// Valid - The domain is allowed
	return true
}
