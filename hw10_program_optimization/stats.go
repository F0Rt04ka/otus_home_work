package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := countDomainsInUserEmails(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}

	return u, nil
}

func countDomainsInUserEmails(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	var user User
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		emailSegments := strings.SplitN(user.Email, "@", 2)
		if len(emailSegments) != 2 {
			continue
		}

		if !strings.HasSuffix(emailSegments[1], domain) {
			continue
		}

		result[strings.ToLower(emailSegments[1])]++
	}

	return result, scanner.Err()
}
