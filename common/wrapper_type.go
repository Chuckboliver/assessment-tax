package common

import "strconv"

type Float64 float64

func (f Float64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(f), 'f', 2, 32)), nil
}
