package conv

import "strconv"

var letterNames = [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func IntDefault(str string, d int) int {
	if v, err := strconv.Atoi(str); err != nil {
		return d
	} else {
		return v
	}
}

func Int64Default(str string, d int64) int64 {
	if v, err := strconv.ParseInt(str, 10, 64); err != nil {
		return d
	} else {
		return v
	}
}

// 正数十进制转 26 进制
func Int2Letter(num int) string {
	if num <= 0 {
		return ""
	}
	var axis string
	for num > 0 {
		if num%26 == 0 {
			axis = letterNames[25] + axis
			num -= 26
		} else {
			axis = letterNames[(num%26)-1] + axis
		}
		num /= 26
	}
	return axis
}
