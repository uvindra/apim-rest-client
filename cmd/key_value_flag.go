package cmd

import (
	"strings"
	//"apim-rest-client/constants"
	)

type FlagMap map[string]string

func (this *FlagMap) String() string {
	var values string

	for key, value := range *this {
		if values != "" {
			values += ","
		}

		values += key + ":" + value
	}
	return values
}

func (this *FlagMap) Set(value string) error {
	splits := strings.Split(value, ":")

	if len(splits) != 2 {
		panic("Map format violated")
	}

	(*this)[splits[0]] = splits[1]

	return nil
}
