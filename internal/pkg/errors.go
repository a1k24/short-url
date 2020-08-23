package pkg

import "fmt"

type UrlAlreadyExistsError string

func (e UrlAlreadyExistsError) Error() string {
	return fmt.Sprintf("URL with name: %s already exists!", string(e))
}
