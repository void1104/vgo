package main

import "fmt"

func main() {
	//fmt.Println("before panic")
	// 比较常见的错误方法是返回error，由调用者决定后续如何处理。但是如果是无法恢复的错误
	// 可以手动触发panic，当然如果在程序运行过程中出现了类似与数组越界的错误，panic也会被触发
	//panic("crash")
	//fmt.Println("after panic")
	//

	// panic会导致程序被中止，但是在退出前，会先处理完当前协程上已经defer的任务，
	// 执行完成后再推出。效果类似于java语言的try...catch
	//defer func(){
	//	fmt.Println("-----------------defer func------------------")
	//}()
	//
	//arr := []int{1,2,3}
	//fmt.Println(arr[4])

	testRecover()
	fmt.Println("after recover")
}

// Go还提供了recover函数，可以避免因为panic发生而导致整个程序终止，recover函数只在defer中生效
func testRecover() {
	defer func() {
		fmt.Println("defer func")
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
	fmt.Println("after panic")
}
