package greetings

import (
	"errors"
	"fmt"
	"math/rand"
)

func Hello(i int, name string) (string, error) {
	if name == "" {
		return "", errors.New("empty name")
	}
	message := fmt.Sprintf(randomFormat(), i, name)
	return message, nil
}

func Hellos(names []string) (map[string]string, error) {
	messages := make(map[string]string)

	for i, name := range names {
		message, err := Hello(i, name)
		if err != nil {
			return nil, err
		}

		messages[name] = message
	}

	return messages, nil
}

func randomFormat() string {
	formats := []string{
		"[%v] Hi, %v. Welcome!",
		"[%v] Great to see you, %v!",
		"[%v] Hail, %v! Well met!",
	}

	return formats[rand.Intn(len(formats))]
}
