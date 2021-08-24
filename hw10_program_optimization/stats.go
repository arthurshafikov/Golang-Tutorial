package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

var ErrEmptyDomain = fmt.Errorf("domain cannot be empty")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, ErrEmptyDomain
	}
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")

		if strings.Contains(email, "."+domain) {
			result[strings.ToLower(strings.SplitN(email, "@", 2)[1])]++
		}
	}

	return result, nil
}
