// Package memstorage provides the necessary methods to use valkey
package memstorage

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/valkey-io/valkey-go"
)

func clientInit() (valkey.Client, error) {
	if os.Getenv("VALKEY_SERVER") == "" {
		err := errors.New("won't connect to valkey, please set VALKEY_SERVER envar")
		return nil, err
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{os.Getenv("VALKEY_SERVER")},
	})

	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetKey tries to find a given key and returns it's value
func GetKey(k string) (string, error) {
	client, err := clientInit()

	if err != nil {
		log.Println(err)
		return "", err
	}

	value, err := client.Do(context.Background(), client.B().Get().Key(k).Build()).ToString()

	if err != nil {
		log.Println(err)
		return "", err
	}

	return value, nil
}

// SetValue tries to create a key with a value on valkey
func SetValue(k string, val string) error {
	client, err := clientInit()
	if err != nil {
		log.Println(err)
		return err
	}

	err = client.Do(context.Background(), client.B().Set().Key(k).Value(val).Build()).Error()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
