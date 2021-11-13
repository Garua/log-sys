package main

import "fmt"

func main() {
	//a := []int{123, 3, 2, 12, 31, 23, 23}
	//insertSort(a, len(a)-1)
	//fmt.Println(a)
	//target := binarySearch(a, 0, len(a)-1, 31)
	//fmt.Println(target)
	a := []int{123, 3, 2, 12, 31, 23, 23}
	shellSort(a)
	fmt.Println("shell short")
	fmt.Println(a)

}

func shellSort(arr []int) {
	//不断缩小增量
	for interval := len(arr) / 2; interval > 0; interval /= 2 {
		//增量为1的插入排序
		//for i := 1;i < len(arr);i++{
		//增量为interval的插入排序,这里不是i = i+interval,
		for i := interval; i < len(arr); i ++ {
			target := arr[i]
			j := i - interval
			for j >= 0 && target < arr[j] {
					arr[j+interval] = arr[j]
					j-=interval
			}
			arr[j+interval] = target
		}
	}
}

//插入排序——递归
func insertSort(arr []int, index int) {
	if index <= 0 {
		return
	}
	insertSort(arr, index-1)
	v := arr[index]
	i := index - 1 //
	for i >= 0 && arr[i] > v {
		arr[i+1] = arr[i]
		i--
	}
	arr[i+1] = v
}

//先排序再调用
func binarySearch(arr []int, l, r, k int) int {
	if r < l {
		return -1
	}
	mid := l + ((r - l) >> 1)
	midV := arr[mid]
	if midV < k {
		return binarySearch(arr, mid+1, r, k)
	} else if midV > k {
		return binarySearch(arr, l, mid-1, k)
	} else {
		return midV
	}
}
