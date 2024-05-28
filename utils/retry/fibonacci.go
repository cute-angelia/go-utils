package retry

// Fibonacci 斐波那契数列
// 十次 0,1,1,2,3,5,8,13,21,34
func Fibonacci(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}
