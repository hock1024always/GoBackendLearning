# 语法基础

## 标识符

### 命名规则

1. 标识符对大小写敏感，count 和 Count 是不同的标识符
2. 下划线可以放开头，但是数字不能放开头

### 关键字及含义

| **类别**     | **关键字**     | **说明**                               |
| :----------- | :------------- | :------------------------------------- |
| **逻辑值**   | `True`         | 布尔真值                               |
|              | `False`        | 布尔假值                               |
|              | `None`         | 表示空值或无值                         |
| **逻辑运算** | `and`          | 逻辑与运算                             |
|              | `or`           | 逻辑或运算                             |
|              | `not`          | 逻辑非运算                             |
| **条件控制** | `if`           | 条件判断语句                           |
|              | `elif`         | 否则如果（else if 的缩写）             |
|              | `else`         | 否则分支                               |
| **循环控制** | `for`          | 迭代循环                               |
|              | `while`        | 条件循环                               |
|              | `break`        | 跳出循环                               |
|              | `continue`     | 跳过当前循环的剩余部分，进入下一次迭代 |
| **异常处理** | **`try`**      | **尝试执行代码块**                     |
|              | **`except`**   | **捕获异常**                           |
|              | **`finally`**  | **无论是否发生异常都会执行的代码块**   |
|              | **`raise`**    | **抛出异常**                           |
| **函数定义** | **`def`**      | **定义函数**                           |
|              | `return`       | 从函数返回值                           |
|              | **`lambda`**   | **创建匿名函数**                       |
| **类与对象** | `class`        | 定义类                                 |
|              | **`del`**      | **删除对象引用**                       |
| **模块导入** | `import`       | 导入模块                               |
|              | `from`         | 从模块导入特定部分                     |
|              | **`as`**       | **为导入的模块或对象创建别名**         |
| **作用域**   | `**global`**   | **声明全局变量**                       |
|              | **`nonlocal`** | **声明非局部变量（用于嵌套函数）**     |
| **异步编程** | **`async`**    | **声明异步函数**                       |
|              | **`await`**    | **等待异步操作完成**                   |
| **其他**     | **`assert`**   | **断言，用于测试条件是否为真**         |
|              | **`in`**       | **检查成员关系**                       |
|              | **`is`**       | **检查对象身份（是否是同一个对象）**   |
|              | **`pass`**     | **空语句，用于占位**                   |
|              | **`with`**     | **上下文管理器，用于资源管理**         |
|              | **`yield`**    | **从生成器函数返回值**                 |

### 注释

```python
# 第一个注释

'''
第二注释
第三注释
'''
```

### 行与缩进（star）

1. python最具特色的就是使用缩进来表示代码块，不需要使用大括号 **{}** 。

2. 缩进的空格数是可变的，但是同一个代码块的语句必须包含相同的缩进空格数

   ```python
   if True:
       print ("Answer")
       print ("True")
   else:
       print ("Answer")
     print ("False")    # 缩进不一致，会导致运行错误
   ```

### 多行语句

Python 通常是一行写完一条语句，但如果语句很长，我们可以使用反斜杠  \ 来实现多行语，情况如下：

```python
item_one = 1
item_two = 2
item_three = 3
total = item_one + \
        item_two + \
        item_three
print(total) # 输出为 6
```

## 数据类型

### 数字(Number)类型

python中数字有四种类型：整数、布尔型、浮点数和复数。

- **int** (整数), 如 1, 只有一种整数类型 int，表示为长整型，没有 python2 中的 Long。
- **bool** (布尔), 如 True。
- **float** (浮点数), 如 1.23、3E-2
- **complex** (复数) - 复数由实部和虚部组成，形式为 a + bj，其中 a 是实部，b 是虚部，j 表示虚数单位。如 1 + 2j、 1.1 + 2.2j

### 字符串(String)

- Python 中单引号 **'** 和双引号 **"** 使用完全相同。
- 使用三引号(**'''** 或 **"""**)可以指定一个多行字符串。
- 转义符 \。
- 反斜杠可以用来转义，使用 **r** 可以让反斜杠不发生转义。 如 **r"this is a line with \n"** 则 **\n** 会显示，并不是换行。
- 按字面意义级联字符串，如 **"this " "is " "string"** 会被自动转换为 **this is string**。
- 字符串可以用 **+** 运算符连接在一起，用 ***** 运算符重复。
- Python 中的字符串有两种索引方式，从左往右以 **0** 开始，从右往左以 **-1** 开始。
- Python 中的字符串不能改变。
- Python 没有单独的字符类型，一个字符就是长度为 1 的字符串。
- 字符串切片 **str[start:end]**，其中 start（包含）是切片开始的索引，end（不包含）是切片结束的索引。
- 字符串的切片可以加上步长参数 step，语法格式如下：**str[start : end : step]**

```python
#!/usr/bin/python3
 
str='123456789'
 
print(str)                 # 输出字符串
print(str[0:-1])           # 输出第一个到倒数第二个的所有字符
print(str[0])              # 输出字符串第一个字符
print(str[2:5])            # 输出从第三个开始到第六个的字符（不包含）
print(str[2:])             # 输出从第三个开始后的所有字符
print(str[1:5:2])          # 输出从第二个开始到第五个且每隔一个的字符（步长为2）
print(str * 2)             # 输出字符串两次
print(str + '你好')         # 连接字符串
 
print('hello\nrunoob')      # 使用反斜杠(\)+n转义特殊字符
print(r'hello\nrunoob')     # 在字符串前面添加一个 r，表示原始字符串，不会发生转义


>>> print('\n')       # 输出空行
>>> print(r'\n')      # 输出 \n
# 输出结果
123456789
12345678
1
345
3456789
24
123456789123456789
123456789你好
------------------------------
hello
runoob
hello\nrunoob
```

## 特殊语法与显示

### 空行

空行的作用在于分隔两段不同功能或含义的代码，便于日后代码的维护或重构.

### 同一行显示多条语句

Python 可以在同一行中使用多条语句，语句之间使用分号 **;** 分割

1. 程序代码

   ```python
   import sys; x = 'runoob'; sys.stdout.write(x + '\n') #输出结果：runoob
   ```

2. 命令行交互

   ```python
   >>> import sys; x = 'runoob'; sys.stdout.write(x + '\n')
   runoob
   7 #输出字符串长度，含\n长度 
   
   >>> import sys
   >>> sys.stdout.write(" hi ")    # hi 前后各有 1 个空格
    hi 4 #前后都有空格
   ```

### 多个语句构成代码组

1. 缩进相同的一组语句构成一个代码块，我们称之代码组。
2. 像if、while、def和class这样的复合语句，首行以关键字开始，以冒号( : )结束，该行之后的一行或多行代码构成代码组。
3. 我们将首行及后面的代码组称为一个子句(clause)

```python
if expression : 
   suite
elif expression : 
   suite 
else : 
   suite
```

### print 输出

**print** 默认输出是换行的，如果要实现不换行需要在变量末尾加上 **end=""**：

```python
#!/usr/bin/python3
 
x="a"
y="b"
# 换行输出
print( x )
print( y )
 
print('---------')
# 不换行输出
print( x, end=" " )
print( y, end=" " )
print()

a
b
---------
a b

```

## import 与 from...import

在 python 用 **import** 或者 **from...import** 来导入相应的模块。

1. 将整个模块(somemodule)导入，格式为： **import somemodule**
2. 从某个模块中导入某个函数,格式为： **from somemodule import somefunction**
3. 从某个模块中导入多个函数,格式为： **from somemodule import firstfunc, secondfunc, thirdfunc**
4. 将某个模块中的全部函数导入，格式为： **from somemodule import \***

```python
# 整模块导入
import sys
print('================Python import mode==========================')
print ('命令行参数为:')
for i in sys.argv:
    print (i)
print ('\n python 路径为',sys.path)

# 导入特定成员
from sys import argv,path  #  导入特定的成员
 
print('================python from import===================================')
print('path:',path) # 因为已经导入path成员，所以此处引用时不需要加sys.path
```

## 基本数据类型与转换

