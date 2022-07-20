package registrationform

import (
	"regexp"
	"strings"
)

type Form struct {
	Firstname string
	Lastname  string
	Email     string
	School    string
	Class     string
	Phones    []Phone
	Comment   string
}

type Phone struct {
	Country string
	Number  string
}

func (p Phone) String() string {
	return p.Country + p.Number
}

// https://regex101.com/r/WaNv4J
var regex_phone = regexp.MustCompile(`(?P<country>\+[0-9]{1,3})?(?: |-)?(?P<numer>(?:[0-9](?: |-)?){7}[0-9])`)

// https://regex101.com/r/po7lKn
var regex_email = regexp.MustCompile(`^\s*([^\s@]+)@([^\s@]+.[^\s@]+)\s*$`)

func extractDigits(s string) string {
	res := []rune{}
	for _, n := range s {
		if '0' <= n && n <= '9' {
			res = append(res, n)
		}
	}
	return string(res)
}

// ParsePhones parses all Phones contained in s and fill Phones field
// ok will be false if number of parsed Phones is less than 1
func (r *Form) ParsePhones(s string) (ok bool) {
	res := regex_phone.FindAllStringSubmatch(s, -1)
	for _, n := range res {
		p := Phone{}
		p.Country = n[1]
		p.Number = extractDigits(n[2])
		if p.Country == "" {
			p.Country = "+371"
		}
		r.Phones = append(r.Phones, p)
	}
	return len(r.Phones) > 0
}

// ParseEmail parses single email from string, trimming spaces if neccesery
// ok will be false if email is not valid or was not found
func (r *Form) ParseEmail(s string) (ok bool) {
	if regex_email.MatchString(s) {
		s = strings.TrimPrefix(s, " ")
		s = strings.TrimSuffix(s, " ")
		r.Email = s
		ok = true
	}
	return
}

func (r Form) GetPhones() string {
	res := make([]string, len(r.Phones))
	for i, n := range r.Phones {
		res[i] = n.String()
	}
	return strings.Join(res, "; ")
}

func (r Form) IsValid() (valid bool, reason string) {
	if r.Firstname == "" {
		return false, "Firstname"
	} else if r.Lastname == "" {
		return false, "Lastname"
	} else if r.School == "" {
		return false, "School"
	} else if r.Class == "" {
		return false, "Class"
	} else if len(r.Phones) == 0 {
		return false, "Phones"
	}
	return true, ""
}
