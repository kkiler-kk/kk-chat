package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	var nums = []int{2, 7, 11, 15}
	var target = 9
	newNums := twoSum(nums, target)
	fmt.Println(newNums)
}
func twoSum(nums []int, target int) []int {
	var capMap = make(map[int]int)
	size := len(nums)
	for i := 0; i < size; i++ {
		value, ok := capMap[target-nums[i]]
		if ok {
			return []int{i, value}
		}
		capMap[nums[i]] = i
	}
	return []int{-1, -1}
}

// 合并数组去重
func ArrayIntersection(arr []int, arr1 []int) []int {
	var intersection []int
	arr = append(arr, arr1...)
	sameElem := make(map[int]int) // 利用map的key唯一性 来去重

	for _, v := range arr {
		if _, ok := sameElem[v]; !ok {
			sameElem[v] = 1
			intersection = append(intersection, v)
		}
	}
	return intersection
}

func f1() ([]int, error) {
	length := 3
	nums := make([]int, length)
	for i := 0; i < length; i++ {
		nums[i] = rand.Intn(100)
	}
	return nums, nil
}
func f2() ([]int, error) {
	length := 3
	nums := make([]int, length)
	for i := 0; i < length; i++ {
		nums[i] = rand.Intn(100)
	}
	return nums, nil
}
func f3() ([]int, error) {
	length := 3
	nums := make([]int, length)
	for i := 0; i < length; i++ {
		nums[i] = rand.Intn(100)
	}
	return nums, nil
}

// 消费者
func consumer(cname string, ch chan int) {

	//可以循环 for i := range ch 来不断从 channel 接收值，直到它被关闭。

	for {
		select {
		case i, ok := <-ch:
			if ok {
				fmt.Println("consumer--------", cname, ":", i, time.Now().Unix())
			} else {
				break
			}

		}

	}
	fmt.Println("ch closed.")
}

// 生产者
func producer(pname string, ch chan int) {
	for i := 0; i < 4; i++ {
		fmt.Println("producer--", pname, ":", i, time.Now().Unix())
		ch <- i
	}
}

func deferCall() {
	defer func() { fmt.Println("打印1") }()
	defer func() { fmt.Println("打印2") }()
	defer func() { fmt.Println("打印3") }()
	panic("异常")
}

type Item struct {
	ID   int64
	Name string
}

func DeleteSlice(list []Item, n int) []Item {
	for i := 0; i < len(list); i++ {
		if i == n {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	return list
}

// 交集
func intersection(m1, m2 map[string]int) map[string]int {
	result := make(map[string]int)
	for k, v := range m1 {
		if _, ok := m2[k]; ok {
			result[k] = v
		}
	}
	return result
}

// 并集
func union(m1, m2 map[string]int) map[string]int {
	result := make(map[string]int)
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		if _, ok := result[k]; !ok {
			result[k] = v
		}
	}
	return result
}
