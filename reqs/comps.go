package reqs

func GT(value, check interface{}) (bool, err) {
	comp := value.(float64) > check.(float64)
	return comp, nil
}

func GTE(value, check interface{}) (bool, err) {
	comp := value.(float64) >= check.(float64)
	return comp, nil
}
func LT(value, check interface{}) (bool, err) {
	comp := value.(float64) < check.(float64)
	return comp, nil
}
func LTE(value, check interface{}) (bool, err) {
	comp := value.(float64) <= check.(float64)
	return comp, nil
}

func EQ(value, check interface{}) (bool, err) {
	comp := value == check
	return comp, nil
}

func NEQ(value, check interface{}) (bool, err) {
	comp := value == check
	return comp, nil
}
