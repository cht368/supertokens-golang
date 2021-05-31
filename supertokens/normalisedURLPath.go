package supertokens

import (
	"errors"
	"net/url"
	"strings"
)

type NormalisedURLPath struct {
	Value string
}

func NewNormalisedURLPath(url string) (*NormalisedURLPath, error) {
	val, err := NormaliseURLPathOrThrowError(url)
	if err != nil {
		return nil, err
	}
	return &NormalisedURLPath{
		Value: val,
	}, nil
}

func (n *NormalisedURLPath) StartsWith(other NormalisedURLPath) bool {
	return strings.HasPrefix(n.Value, other.Value)
}

func (n *NormalisedURLPath) AppendPath(other NormalisedURLPath) NormalisedURLPath {
	return NormalisedURLPath{Value: n.Value + other.Value}
}

func (n *NormalisedURLPath) Equals(other NormalisedURLPath) bool {
	return n.Value == other.Value
}

func (n *NormalisedURLPath) IsARecipePath() bool {
	return n.Value == "/recipe" || strings.HasPrefix(n.Value, "/recipe/")
}

func NormaliseURLPathOrThrowError(input string) (string, error) {
	input = strings.ToLower(strings.Trim(input, ""))
	if strings.HasPrefix(input, "http://") != true && strings.HasPrefix(input, "https://") != true {
		return "", errors.New("converting to proper URL")
	}

	urlObj, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	input = urlObj.Path

	if input[len(input)-1:] == "/" {
		return input[:len(input)-1], nil
	}

	return input, nil
}
