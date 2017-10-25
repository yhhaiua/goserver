package common

//Max 比较大小
func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

//Min 比较大小
func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

//Alignment 字节对齐
func Alignment(value, num int) int {
	newlen := value
	surplus := value % num
	if surplus > 0 {
		newlen += num - surplus
	}
	return newlen
}
