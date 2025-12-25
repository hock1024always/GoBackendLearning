### NOSQL和SQL的区别？

SQL数据库，指关系型数据库 - 主要代表：SQL Server，Oracle，MySQL(开源)，PostgreSQL(开源)。

关系型数据库存储结构化数据。这些**数据逻辑上以行列二维表的形式存在**，每一列代表数据的一种属性，每一行代表一个数据实体。

![image-20240725232218438](https://cdn.xiaolincoding.com//picgo/image-20240725232218438.png)

NoSQL指非关系型数据库 ，主要代表：MongoDB，Redis。NoSQL 数据库逻辑上提供了不同于二维表的存储方式，**存储方式可以是JSON文档、哈希表或者其他方式**。

![image-20240725232206455](https://cdn.xiaolincoding.com//picgo/image-20240725232206455.png)

选择 SQL vs NoSQL，考虑以下因素。

> ACID vs BASE

关系型数据库支持 ACID 即原子性，一致性，隔离性和持续性。相对而言，NoSQL 采用更宽松的模型 BASE ， 即基本可用，软状态和最终一致性。

从实用的角度出发，我们需要考虑对于面对的应用场景，ACID 是否是必须的。比如银行应用就必须保证 ACID，否则一笔钱可能被使用两次；又比如社交软件不必保证 ACID，因为一条状态的更新对于所有用户读取先后时间有数秒不同并不影响使用。

对于需要保证 ACID 的应用，我们可以优先考虑 SQL。反之则可以优先考虑 NoSQL。

> 扩展性对比

NoSQL数据之间无关系，这样就非常容易扩展，也无形之间，在架构的层面上带来了可扩展的能力。比如 redis 自带主从复制模式、哨兵模式、切片集群模式。

相反关系型数据库的数据之间存在关联性，水平扩展较难 ，需要解决跨服务器 JOIN，分布式事务等问题。

### 数据库三大范式是什么？

**第一范式（1NF）：要求数据库表的每一列都是不可分割的原子数据项。**

举例说明：

![img](https://cdn.xiaolincoding.com//picgo/1218459-20180909201651535-1215699096.png)

在上面的表中，“家庭信息”和“学校信息”列均不满足原子性的要求，故不满足第一范式，调整如下：

![img](https://cdn.xiaolincoding.com//picgo/1218459-20180909202243826-1032549277.png)

可见，调整后的每一列都是不可再分的，因此满足第一范式（1NF）；

**第二范式（2NF）：在1NF的基础上，非码属性必须完全依赖于候选码（在1NF基础上消除非主属性对主码的部分函数依赖）**

**第二范式需要确保数据库表中的每一列都和主键相关，而不能只与主键的某一部分相关（主要针对联合主键而言）。**

举例说明：

![img](https://cdn.xiaolincoding.com//picgo/1218459-20180909204750951-639647799.png)

在上图所示的情况中，同一个订单中可能包含不同的产品，因此主键必须是“订单号”和“产品号”联合组成，

但可以发现，产品数量、产品折扣、产品价格与“订单号”和“产品号”都相关，但是订单金额和订单时间仅与“订单号”相关，与“产品号”无关，

这样就不满足第二范式的要求，调整如下，需分成两个表：

![img](https://cdn.xiaolincoding.com//picgo/1218459-20180909210444227-1008056975.png)

![img](https://cdn.xiaolincoding.com//picgo/1218459-20180909210458847-2092897116.png)

**第三范式（3NF）：在2NF基础上，任何非主[属性 (opens new window)](https://baike.baidu.com/item/属性)不依赖于其它非主属性（在2NF基础上消除传递依赖）**

**第三范式需要确保数据表中的每一列数据都和主键直接相关，而不能间接相关。**

举例说明：

![img](https://cdn.xiaolincoding.com//picgo/1218459-20180909211311408-1364899740.png)

上表中，所有属性都完全依赖于学号，所以满足第二范式，但是“班主任性别”和“班主任年龄”直接依赖的是“班主任姓名”，

而不是主键“学号”，所以需做如下调整：

![img](https://cdn.xiaolincoding.com//picgo/1218459-20180909211539242-1391100354.png)

![img](https://cdn.xiaolincoding.com//picgo/1218459-20180909211602202-1069383439.png)

这样以来，就满足了第三范式的要求。

### MySQL 怎么连表查询？

数据库有以下几种联表查询类型：

1. **内连接 (INNER JOIN)**
2. **左外连接 (LEFT JOIN)**
3. **右外连接 (RIGHT JOIN)**
4. **全外连接 (FULL JOIN)**

![img](https://cdn.xiaolincoding.com//picgo/1721710415166-eff24e6c-555c-436c-b1b8-7c6dbb5850d7.webp)

**1. 内连接 (INNER JOIN)**

内连接返回两个表中有匹配关系的行。**示例**:

```sql
SELECT employees.name, departments.name
FROM employees
INNER JOIN departments
ON employees.department_id = departments.id;
```

这个查询返回每个员工及其所在的部门名称。

**2. 左外连接 (LEFT JOIN)**

左外连接返回左表中的所有行，即使在右表中没有匹配的行。未匹配的右表列会包含NULL。**示例**:

```sql
SELECT employees.name, departments.name
FROM employees
LEFT JOIN departments
ON employees.department_id = departments.id;
```

这个查询返回所有员工及其部门名称，包括那些没有分配部门的员工。

**3. 右外连接 (RIGHT JOIN)**

右外连接返回右表中的所有行，即使左表中没有匹配的行。未匹配的左表列会包含NULL。**示例**:

```sql
SELECT employees.name, departments.name
FROM employees
RIGHT JOIN departments
ON employees.department_id = departments.id;
```

这个查询返回所有部门及其员工，包括那些没有分配员工的部门。

**4. 全外连接 (FULL JOIN)**

全外连接返回两个表中所有行，包括非匹配行，在MySQL中，FULL JOIN 需要使用 UNION 来实现，因为 MySQL 不直接支持 FULL JOIN。**示例**:

```sql
SELECT employees.name, departments.name
FROM employees
LEFT JOIN departments
ON employees.department_id = departments.id

UNION

SELECT employees.name, departments.name
FROM employees
RIGHT JOIN departments
ON employees.department_id = departments.id;
```

这个查询返回所有员工和所有部门，包括没有匹配行的记录。

### MySQL如何避免重复插入数据？

**方式一：使用UNIQUE约束**

在表的相关列上添加UNIQUE约束，确保每个值在该列中唯一。例如：

```sql
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) UNIQUE,
    name VARCHAR(255)
);
```

如果尝试插入重复的email，MySQL会返回错误。

**方式二：使用INSERT ... ON DUPLICATE KEY UPDATE**

这种语句允许在插入记录时处理重复键的情况。如果插入的记录与现有记录冲突，可以选择更新现有记录：

```sql
INSERT INTO users (email, name) 
VALUES ('example@example.com', 'John Doe')
ON DUPLICATE KEY UPDATE name = VALUES(name);
```

**方式三：使用INSERT IGNORE**： 该语句会在插入记录时忽略那些因重复键而导致的插入错误。例如：

```sql
INSERT IGNORE INTO users (email, name) 
VALUES ('example@example.com', 'John Doe');
```

如果email已经存在，这条插入语句将被忽略而不会返回错误。

选择哪种方法取决于具体的需求：

- 如果需要保证全局唯一性，使用UNIQUE约束是最佳做法。
- 如果需要插入和更新结合可以使用`ON DUPLICATE KEY UPDATE`。
- 对于快速忽略重复插入，`INSERT IGNORE`是合适的选择。