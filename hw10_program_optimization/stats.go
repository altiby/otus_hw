package hw10programoptimization

import (
	"bufio"
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

//easyjson:json
type users [100000]User

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func getUsers(r io.Reader) (result users, err error) {
	buffer := bufio.NewReader(r)
	counter := 0
	var data []byte
	for {
		data, _ = buffer.ReadSlice('\n')
		if len(data) == 0 {
			break
		}
		if marshalErr := easyjson.Unmarshal(data, &result[counter]); marshalErr != nil {
			err = marshalErr
			return
		}
		counter++
	}
	return
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
