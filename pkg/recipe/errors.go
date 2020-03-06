package recipe

import "errors"

var ErrUnmarshalResponse = errors.New("failed to unmarshal response")
var ErrNoResults = errors.New("failed to retrieve results")
