package main

import "fmt"

// 定义接口
type Speaker interface {
	Speak()
}

// 父结构体
type Animal struct {
	Name string
}

// 实现接口方法
func (a *Animal) Speak() {
	fmt.Println(a.Name, "says hello!")
}

// 子结构体
type Dog struct {
	Animal
	Breed string
}

func main() {
	dog := Dog{
		Animal: Animal{Name: "Buddy"},
		Breed:  "Golden Retriever",
	}
	dog.Speak() // 通过接口调用方法
}
