package main

import (
	"fmt"
	"time"
)

func main() {
	//arr := []int{123,3,2,12,31,23,23}
	//fmt.Println(add(arr,0))
	//str := "hello world"
	//fmt.Println(reverseRe(str, 0))
	now := time.Now()
	fmt.Println(fib(50))
	milliseconds := time.Since(now).Milliseconds()
	fmt.Println(milliseconds)
	//fmt.Println(gcd1(100, 25))
	//arr := []int{123,3,2,12,31,23,23}
	//insertSort(arr,len(arr)-1)
	//fmt.Println("--------------")
	//fmt.Println(arr)
	//printHanoiTower(3,"A","B","C")

}

func printHanoiTower(n int,from,to,help string)  {
	if n == 1 {
		fmt.Printf("move %d from %s to %s  help is %s\n",n,from,to,help)
		return
	}
	printHanoiTower(n-1,from,help,to) //把n-1个盘子移动到辅助盘
	fmt.Printf("move  %d from %s to %s  help is %s\n",n,from,to,help)//n就能顺利到达目标
	printHanoiTower(n-1,help,to,from)//把n-1个盘子移动到目标盘

}
//递归版插入排序
func insertSort(arr []int, k int) {
	if k == 0{
		return
	}
	//对前k-1个元素排序
	insertSort(arr, k-1)
	//把位置k的元素插入到前面的部分，执行下面的代码时，0到K-1已经有序了
	x := arr[k]
	index := k - 1

	for index >= 0 && x < arr[index] {
		//把数往后挪
		arr[index+1] = arr[index]
		index--
	}
	//退出for循环时，index为-1,所以要+1
	arr[index+1] = x

}

// 求最大公约数，展转相除
func gcd(m, n int) int {
	if n == 0 {
		return m
	}
	return gcd(n, m%n)
}

// 更相减损术
func gcd1(m, n int) int {
	if m == n {
		return n
	}
	if m <= n {
		m, n = n, m
	}
	return gcd1(n, m-n)
}

//斐波拉
func fib(n int) int {
	if n == 1 || n == 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

//迭代反转
func reverseFor(s string) (str string) {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < len(runes) && i <= j; i++ {
		runes[i], runes[j] = runes[j], runes[i]
		j--
	}
	return string(runes)
}

//递归反转
func reverseRe(s string, index int) (result string) {
	if index == len(s) {
		return result
	}
	result = string(s[index])
	tmp := reverseRe(s, index+1)
	return tmp + result
}

func t(s string) {
	//通过索引，得到的类型是uint8-->byte
	for i := 0; i < 1; i++ {
		fmt.Printf("%T\n", s[i])
	}
	//通过range，到得的类型是int32-->rune
	for _, v := range s {
		fmt.Printf("%T\n", v)
		break
	}
}

//递归求和
func add(a []int, index int) (sum int) {
	if index >= len(a) {
		return sum
	}

	sum = a[index] + add(a, index+1)
	return sum
}
