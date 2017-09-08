package packer

import (
	"errors"

	"github.com/RaniSputnik/lovepac/target"
)

const (
	FormatStarling = "starling"
	FormatLove     = "love"
)

var ErrFormatIsInvalid = errors.New("Format is not valid")

var formatLookup = map[string]target.Format{
	FormatStarling: target.Starling,
	FormatLove:     target.Love,
}

func FormatIsValid(format string) bool {
	_, ok := formatLookup[format]
	return ok
}

func GetFormatNamed(format string) target.Format {
	if !FormatIsValid(format) {
		return target.Unknown
	}
	return formatLookup[format]
}
