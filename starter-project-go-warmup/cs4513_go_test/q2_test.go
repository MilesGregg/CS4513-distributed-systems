package cs4513_go_test

import (
	"fmt"
	"testing"

	"project-go-warmup/cs4513_go_impl"
)

func test(t *testing.T, fileName string, num int, expected int) {
	result := cs4513_go_impl.Sum(num, fileName)
	if result != expected {
		t.Fatal(fmt.Sprintf(
			"Sum of %s failed: got %d, expected %d\n", fileName, result, expected))
	}
}

func Test1(t *testing.T) {
	test(t, "q2_test1.txt", 1, 499500)
}

func Test2(t *testing.T) {
	test(t, "q2_test1.txt", 10, 499500)
}

func Test3(t *testing.T) {
	test(t, "q2_test2.txt", 1, 117652)
}

func Test4(t *testing.T) {
	test(t, "q2_test2.txt", 10, 117652)
}

func Test5(t *testing.T) {
	test(t, "q2_test3.txt", 1, 617152)
}

func Test6(t *testing.T) {
	test(t, "q2_test3.txt", 10, 617152)
}

func Test7(t *testing.T) {
	test(t, "q2_test4.txt", 1, 4995000)
}

func Test8(t *testing.T) {
	test(t, "q2_test4.txt", 10, 4995000)
}

func Test9(t *testing.T) {
	test(t, "q2_test5.txt", 1, 49950000)
}

func Test10(t *testing.T) {
	test(t, "q2_test5.txt", 10, 49950000)
}
