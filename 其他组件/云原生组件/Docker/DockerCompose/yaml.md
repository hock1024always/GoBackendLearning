# yaml 语法

[YAML 入门教程 | 菜鸟教程](https://www.runoob.com/w3cnote/yaml-intro.html)，直接看菜鸟教程比较方便

## 1. 概念

> [!note]
>
> ​	YAML（YAML Ain't Markup Language）是一种人类可读的数据序列化语言。它的设计目标是使数据在不同编程语言之间交换和共享变得简单。YAML采用了一种简洁、直观的语法，以易于阅读和编写的方式表示数据结构。
>
> ​	YAML广泛应用于配置文件、数据序列化、API设计和许多其他领域。它被许多编程语言和框架所支持，包括Python、Java、Ruby等。在Python中，可以使用PyYAML库来读取和写入YAML文件。
>
> ​	YAML的优点包括易读性高、易于理解、与多种编程语言兼容以及支持丰富的数据结构。它的简洁语法使得配置文件变得更加直观和可维护。无论是作为配置文件格式还是数据交换格式，YAML都是一个强大而受欢迎的选择。

## 2. 语法规范

YAML的语法特点包括：

1. 使用缩进表示层级关系，不使用大括号或者其他符号。
2. 使用冒号来表示键值对。
3. 支持列表和嵌套结构。
4. 使用注释以 "#" 开头。
5. 支持引用和锚点，可以在文档中引用其他部分的数据。



```yaml
# 1. 字符串
name: "John"
addr: "长沙"

# 2. 数字
age: 30

# 3. boor
isStudent: true
isTeacher: false

# 4. 列表,使用短横线（-）表示列表项，列表项之间使用换行进行分隔
fruits:
	- apple
	- banana
	- orange
# 另外一种表示
fruits: [apple, banana, orange]

# 5. 字典 
person:
	name: "John"
	age: 30
	
# 6. 空值
status: null

```

