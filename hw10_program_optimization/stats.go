package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100000]User

func getUsers(r io.Reader) (result users, err error) {
	buffer := bufio.NewReader(r)
	delim := byte('\n')
	counter := 0
	var line []byte
	for {
		line, err = buffer.ReadBytes(delim)
		if err != nil && !errors.Is(err, io.EOF) {
			return
		}

		if marshalErr := easyjson.Unmarshal(line, &result[counter]); marshalErr != nil {
			err = marshalErr
			return
		}
		if errors.Is(err, io.EOF) {
			err = nil
			return
		}
		counter++
	}
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	exp, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		matched := exp.Match([]byte(user.Email))

		if matched {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
