<img src=".\images-docker\Docker架构一图解.png" alt="Docker架构一图解" style="zoom:80%;" />



![架构图](.\images-docker\架构图.png)



# 1.容器原理

容器原理简单了解即可，只需要理解什么是容器，什么是镜像。

## 1.1 Linux 基础知识

### 1.1.1 Linux namespaces

**每个命名空间提供了一种特定类型的资源隔离**，如进程、网络、文件系统等。

所有的进程都会属于至少一个`namespace`，可以将一组进程与其他进程隔离开，使它们拥有**各自独立的资源副本**。

而且每个命名空间都有一个**唯一的标识符**，以区分不同的命名空间。

一个命名空间中可以有**多个进程**，进程可以在所属的命名空间内**自由管理和配置这些资源**，而**不会影响其他命名空间中的进程**

![namespace](.\images-docker\namespace.png)

| 命名空间类型                     | 控制内容                        |
| -------------------------------- | ------------------------------- |
| Mount(mnt)                       | 隔离挂载点                      |
| Time                             | 隔离时钟                        |
| UTS                              | 隔离主机名和NIS域名             |
| User ID(user)                    | 隔离用户和用户组 ID             |
| Process ID (pid)                 | 隔离进程ID                      |
| Network(net)                     | 隔离网络设备、栈、端口等        |
| Control group(cgroup)            | 隔离Cgroup根目录                |
| Inter-process Communication(ipc) | 隔离System V IPC, POSIX消息队列 |

可以通过`sudo ll /proc/<pid>/ns`命令查看进程所属的命名空间

* 每个`namespace`都使一个软连接，其名称即为命名空间类型
* 每个软连接指向该进程所属的真正`namespace`对象，该对象用`inode`号码标识，每个号码是唯一的
* `<ns>_for_children`指示了该进程对应`namespace`的对象



### 1.1.2 Linux cgroups

cgroups (Linux Control groups) 是 Linux 内核提供的一种可以限制单个进程或者多个进程所使用资源的机制，可以对 CPU，内存等资源实现精细化的控制。

cgroups 为每个可控制的资源定义一个子系统，当创建一个 cgroup 实例时，必须至少指定一种子系统。

cgroups主要功能包括：

* **资源限制**。可以限制进程组的资源使用，例如可以设置内存限制不超过特定值（包括虚拟内存）。

* **优先级设置**。可以为进程组设置不同的优先级，以确保重要的进程获得更多的资源分配。

* **资源统计与监控**。cgroups可以收集和记录进程组的资源使用情况，如CPU使用时间、内存使用量等。

* **进程组控制**。可以冻结、检查和重启动进程组

| 子系统类型 | 描述                                                         |
| ---------- | ------------------------------------------------------------ |
| cpu        | 主要限制进程的   CPU 使用率                                  |
| cpuacct    | 可以统计   cgroups 中进程的 CPU 使用报告                     |
| cpuset     | 可以为   cgroups 中的进程分配独立的 CPU 节点或内存节点       |
| memory     | 可以限制进程的内存使用量                                     |
| blkio      | 可以限制进程的块设备   I/O                                   |
| devices    | 可以限制进程的块设备   I/O                                   |
| net_cls    | 可以标记   cgroups 中进程的网络数据包，并使用 tc 模块进行控制 |
| freezer    | 可以挂起或恢复   cgroups 中的进程                            |
| ns         | 可以使不同   cgroups 下的进程使用不同的命名空间              |



cgroup结构体内可以组织成**树**的形式，每个结构体组成的树称之为一个 cgroups 层级结构。每个子系统只能被attach到一个cgroups层级结构中。

![cgroup](.\images-docker\cgroup.png)

最下面的**P**代表一个进程。每一个进程的描述符中有一个指针指向了一个**辅助数据结构`css_set（cgroups subsystem set）`**。指向某一个`css_set`的进程会被加入到当前`css_set`的进程链表中。**进程与css_set之间的关系是一对多**，隶属于同一css_set的进程受到同一个css_set所关联的资源限制。

`"M×N Linkage"`是一种辅助数据结构，实现css_set与cgroup节点之间多对多的关联。



可以使用`cat /proc/<pid>/cgroup`命令来查看指定进程属于哪些cgroup

每一行包含用冒号隔开的三列，含义分别是：

* 树的 ID， 和 /proc/cgroups 文件中的 ID 一一对应。

* 树绑定的所有 subsystem，多个 subsystem 之间用逗号隔开（systemd 表示没有和任何 subsystem 绑定）。

* 树中的路径，即进程所属的 cgroup，这个路径是相对于挂载点的相对路径。



挂载点路径通过` mount | grep cgroup `命令来查看



### 1.1.3 Linux Capabilities

Linux 内核将超级用户的权限分解为细颗粒度的单元，这些单元称为 capabilities。

现在权限检查的过程就变成了：在执行特权操作时，如果线程的有效身份不是 root，就去检查其是否具有该特权操作所对应的 capabilities，并以此为依据，决定是否可以执行特权操作。

**Capabilities 可以在进程执行时赋予，也可以直接从父进程继承**。Linux capabilities 分为**进程** capabilities 和**文件**  capabilities。

对于进程来说，capabilities 是细分到线程的，即每个线程可以有自己的 capabilities。对于文件来说，capabilities 保存在文件的扩展属性中。

![capability](.\images-docker\capability.png)

文件的 capabilities 包含了 3 个集合：

* `Permitted`：这个集合中包含的 capabilities，在文件被执行时，会与线程的 Bounding 集合计算交集，然后添加到线程的 Permitted 集合中。

* `AInheritable`：这个集合与线程的 Inheritable 集合的交集，会被添加到执行完 execve() 后的线程的 Permitted 集合中。

* `Effective`：这是一个标志位。如果设置开启，那么在执行完 execve() 后，线程 Permitted 集合中的 capabilities 会自动添加到它的 Effective 集合中。

线程的 capabilities 包含了 5 个集合，其中有 3 个与文件的 capabilities 一致：

* `Permitted`：定义了线程能够使用的 capabilities 的上限。
* `AInheritable`：当执行 exec() 系统调用时，能够被新的可执行文件继承的 capabilities，被包含在 Inheritable 集合中。
* `AEffective`：内核检查线程是否可以进行特权操作时，检查的对象便是 Effective 集合。
* `Bounding`：Bounding 集合是 Inheritable 集合的超集，如果某个 capability 不在 Bounding 集合中，即使它在 Permitted 集合中，该线程也不能将该 capability 添加到它的 Inheritable 集合中。
* `Ambient`：Linux4.3 内核新增了一个 capabilities 集合叫 Ambient ，用来弥补 Inheritable 的不足，便于子进程继承父进程的特权。



## 1.2 容器(Container)

“容器” 目前还没有严格的定义，狭义上可以认为Linux 容器是一个被隔离（`Linux namespaces`）和被限制（`Linux cgroups， Linux Capabilities`）的 Linux 进程。

要在 Linux 中启动一个进程，需要` fork/exec/clone`；但要**启动一个容器，还需要创建 `namespace`，配置 `cgroups`** 等，也就是对进程进行容器化

**容器是一个将应用程序及其必要依赖打包在一起构成的标准化执行单元**

![container](.\images-docker\container.png)

容器工具链自下而上可以分为五层：

* 低级运行时
* 垫片
* 高级运行时
* 容器引擎
* 容器编排器



## 1.3 低级容器运行时

低级容器运行时一般指按照 OCI 规范、能够接收可运行 roofs 文件系统和配置文件并运行隔离进程的实现。低级运行时的功能有限，通常执行运行容器的低级任务，不提供存储实现和网络实现。低级运行时的特点是底层、轻量，限制清晰。

## 1.4 高级运行时

一般认为高级运行时负责管理容器进程的生命周期。经典的高级运行时有著名的 `containerd` 和` CRI-O`

## 1.5 容器镜像（Container Image)

它将应用程序及其依赖打包成一个可移植的单元，实现了应用程序与底层系统的隔离，使得应用程序可以在任何操作系统、云环境或物理机上平稳运行。这消除了环境差异和配置问题的烦恼，提供了一致性和可重复性的运行环境

具体而言，**容器镜像是一个静态文件，包含了可执行代码和所需的依赖项，用于在计算基础设施上运行隔离的进程**。它是不可更改的，可以被部署在不同环境中。容器镜像由系统库、工具和平台设置组成，与宿主机共享操作系统内核。通过构建文件系统层来创建镜像，以实现组件的重用

简单来讲容器镜像类似于文件快照的形式。

## 1.6 容器镜像服务

容器镜像服务是一种 PaaS 级别的云计算服务，用于存储和分发容器镜像的云服务。它允许用户创建、存储和分享容器镜像，以便在容器化环境中使用。

该服务提供了一个中心化的存储库，用户可以将自己创建的容器镜像上传到该存储库，并从中下载已有的镜像。容器镜像服务还提供了版本控制、权限管理和安全性措施，以确保镜像的可靠性和完整性

例如阿里云ACR

## 1.7 容器引擎

容器引擎是一种软件平台，可以在同一个操作系统内核上运行多个独立的容器实例。它们接受用户请求，包括命令行选项，拉取镜像，并按照用户的需求运行容器

通常，容器引擎直接负责：

* 处理来自用户(或容器编排器)的请求
* 从远程仓库拉取容器镜像
* 解包容器镜像
* 准备容器挂载点
* 准备需要传递给容器运行时的元数据
* 调用容器运行时

容器引擎的主要目标是提供一种轻量级的虚拟化技术，允许应用程序在隔离的环境中运行，同时共享同一个操作系统内核

常见的容器引擎有：

* Docker：目前最流行的容器引擎，提供了强大的容器管理和部署功能
* Podman：Podman是一个用于管理和运行OCI容器（包括Docker镜像）的命令行工具，它提供了与Docker兼容的API和功能。
* rkt（Rocket）：rkt是CoreOS团队开发的容器引擎，旨在提供更高的安全性和可移植性。

## 1.8 容器编排 Container Orchestration

容器编排是一种自动化的技术，用于在**容器化应用程序中进行资源的配置、部署、网络连接、扩展和管理**。它为开发人员和运维团队提供了一种**无需关注底层基础架构的方式**来管理容器化应用程序。

通过容器编排，可以**轻松地创建和管理大规模的容器集群**。它自动处理容器的创建、分发和调度，确保应用程序在集群中平稳运行。容器编排还提供了自动化的扩展功能，根据应用程序的需求，它可以动态地增加或减少容器的数量，以适应流量的变化。

容器编排还**负责管理容器之间的网络连接，使它们能够相互通信和协作**。它提供了一种灵活的网络配置方式，可以定义容器之间的通信规则和策略，确保应用程序的网络连接安全和可靠。

此外，容器编排还提供了容器的**生命周期管理功能**。它可以监控容器的状态，自动恢复失败的容器，并处理容器的更新和升级。

常见的容器编排工具有：`Kubernetes`、`Docker Swarm`、`Apache Mesos`



# 2. Docker原理

关于Docker的原理可以不需要深入了解，只需要知道如何使用即可，可以直接学习`3.Docker 操作命令`

## 2.1 Docker简介

Docker是一种开源的容器化平台，用于构建、分享和运行应用程序。它通过将应用程序及其所有依赖项打包到一个轻量级、可移植的容器中，实现了应用程序与底层操作系统的解耦。

## 2.2 Docker引擎架构

### 2.2.1 现阶段Docker引擎架构

阶段Docker引擎架构的完整工作流程：

**由容器镜像生成容器的过程**

1. 由Docker客户端发送操作命令（如docker run）并将其转换为对应的REST API并发送到正确的API端点
2. Docker daemon接收到容器创建命令后向containerd发出gRPC调用
3. containerd将Docker镜像转化为OCI bundle并交给runc基于该bundle创建新容器
4. runc通过与操作系统内核接口进行通信，基于必要工具（如namespace、cgroups等）创建容器

<img src=".\images-docker\Docker.png" alt="Docker" style="zoom: 67%;" />

#### 2.2.1.1 runc

runc是一个功能强大的**命令行交互工具**，它通过对Libcontainer进行二次包装，提供了一种简单而高效的方式来管理容器。

runc的主要目标是创建和运行根据OCI格式（镜像格式规范）打包的应用程序。OCI格式定义了容器镜像的结构，包括文件系统、配置信息和运行时参数等。

#### 2.2.1.2 containerd

containerd是一个由Docker公司开发的容器执行工具，它是将Docker daemon进行拆解后重构的一部分。

虽然container和runc都是容器运行时，但是职责不同、表现形式也不同。containerd作为一个常驻进程，**负责持续地监视和管理容器的状态**。它提供了一组API，**使用户能够与容器进行交互，包括创建、启动、停止和销毁容器等操作**。然而，runc仅是一个命令行工具，**用于创建容器**，但在容器**创建完成之后即结束其生命周期**。

`systemctl status containerd`使用该命令可以查看containerd服务器状态。

#### 2.2.1.3 containerd-shim

containerd-shim是一个用于管理容器进程的工具。

容器进程需要一个附近程来维护状态和处理输入输出，但使用`containerd`作为容器管理器时，它可能会异常终止，导致所有以器为父进程的容器进程也会终止，不够稳定和可靠

所以引入`containerd-shim`作为容器的父进程解决稳定性和可靠性问题。

**运行流程**：当containerd接收到来自Docker daemon的请求时，会创建一个名为containerd-shim的进程。containerd-shim通过调用runc命令行工具来启动容器进程。一旦runc启动完容器进程，它会直接退出，此时containerd-shim接管成为容器的父进程。

运行流程图图如下：

<img src=".\images-docker\containner.png" alt="containner" style="zoom: 67%;" />

containerd-shim的**主要责任是收集容器进程的状态并上报给containerd**。它通过与内核进行交互，监控容器的资源使用情况，例如CPU、内存和文件系统等。一旦收集到这些信息，containerd-shim会将其发送给containerd，以便后者可以及时更新容器的状态并做出相应的管理决策。

> [!note]
>
> 通过引入containerd-shim，容器的生命周期和管理变得更加可靠和稳定。即使containerd进程异常终止，容器进程也能够保持运行，避免了因containerd的异常导致整个容器环境的崩溃。



#### 2.2.1.4 Docker daemon

Docker daemon是**Docker引擎的后台服务**，也被称为Docker守护进程或Docker服务。

在Linux主机上**表现为一个常驻进程，本质上是一个命令`dockerd`的可执行文件**，通常位于`/usr/bin`目录。

负责监听来自Docker Clinet的REST API请求，并管理Docker对象（包括镜像、容器、网络、卷以及插件等）。

Docker daemon可细分为以下几个模块

![DockerDaemon](.\images-docker\DockerDaemon.png)

#### 2.2.1.5 Docker client

Docker Client作为Docker架构中与Docker daemon通信的客户端,在Linux主机上以名为`docker`的可执行文件的形式存在，通常位于`/usr/bin`目录下。

其主要职责是向Docker daemon进程发送对容器、镜像等操作的请求，例如`docker pull`和`docker run`。（`docker`命令）

Docker Client具备访问本地守护进程（Local Docker daemon）和远程守护进程（Remote Docker daemon）的能力。

它的生命周期可描述为**从发起REST API给Docker daemon开始，到接收到Docker daemon的响应结束**。在此期间，Docker Client与Docker daemon进行交互，完成一系列请求与响应的过程

<img src=".\images-docker\DockerClient.png" alt="DockerClient" style="zoom:50%;" />

## 2.3 Docker 容器镜像

典型的 Linux 文件系统由 `bootfs `和` rootfs `两部分组成。

1. `bootfs `是引导文件系统，它包含了引导时所需的相关文件，如 bootloader、内核映像和引导配置文件等。在引导过程中，bootfs 被加载到内存中，并在**启动完成后被卸载**（umount），因此在进入系统后，无法直接访问到 `bootfs`。

2. `rootfs` 是根文件系统，也称为根目录（root directory），它是操作系统的根节点。`rootfs` **包含了操作系统中的所有文件和目录结构**，例如 `/etc`、`/proc`、`/bin` 等标准目录。进入系统后，看到的就是 rootfs，它提供了整个操作系统的文件层次结构。

Dockers镜像可以被形象的看作一个`rootfs`。

具体来说，Docker容器镜像指的是**正在运行的容器所使用的隔离文件系统**。这个隔离的文件系统由容器镜像提供，而容器镜像**必须包含运行应用程序所需的所有内容**，包括所有依赖项、配置、脚本、二进制文件等。此外，镜像还包含容器的其他配置，例如环境变量、默认的运行命令以及其他元数据。

每个Docker镜像包含一个或多个只读镜像层，这些镜像层被按顺序叠加在一起，形成一个完整的镜像。

每个镜像层都是**只读**的，因此原始镜像的内容保持不变，而**后续的镜像层只包含了与前一层之间的差异**。每一层镜像都含有指向父层镜像的指针，只有最底层没有指针。

这种分层的结构使得镜像的存储和传输更加高效，因为可以共享和重复使用已有的镜像层。这些只读镜像层的**存储、读取和写入等操作均由存储驱动来完成**。

<img src=".\images-docker\ContainerImage.png" alt="ContainerImage" style="zoom: 50%;" />



**Docker镜像与Docker容器的关系可以类比于代码和程序的关系。就像代码是用来创建程序的静态文件一样，Docker镜像是用来创建容器的静态文件。**

> 什么是静态文件:
> 静态文件是指**不是由服务器生成的文件，例如脚本，CSS文件，图像等，但是必须在请求时发送给浏览器**]
>
> https://blog.csdn.net/wchasedream/article/details/107381428
>
> 静态文件包括图片、视频、网站中的文件（html、css、js）、软件安装包、apk文件、压缩包文件等
>
> https://help.aliyun.com/zh/cdn/what-are-static-content-and-dynamic-content

每个Docker容器都是通过基于指定的Docker镜像创建而来的。基于Docker镜像创建容器时，Docker会在该镜像层之上创建一个可读写的薄层，通常被称为可写层或容器层。这个容器层允许容器对文件系统进行写操作，并且所有的写操作都将被记录在这个容器层中。容器独享容器层，多个容器间共享只读镜像层

<img src=".\images-docker\ContainerImage2.png" alt="ContainerImage2" style="zoom:67%;" />



简单讲容器和容器镜像的关系：容器 = 容器镜像 + 可读写层（Read-Write layer)

`Running Container` 运行态容器 = 一层读写层+多层只读层+隔离的进程空间和包含其中的进程

![ContainerImage3](.\images-docker\ContainerImage3.png)

## 2.4 Docker存储驱动(Storage Driver)

### 2.4.1 为什么需要存储驱动程序

Docker 容器镜像采用分层结构，因此引入存储驱动程序来处理容器镜像的分层设计。这确保容器的文件系统能够被高效、一致性和可靠地管理（主要是存储和获取）。

### 2.4.2 写时复制(Cow, copy-on-write)策略

CoW流程：

* 读取文件：当容器需要读取一个文件时，按照**自上而下的顺序检查镜像的各个层级**（包括容器层），查找到所需的文件后，直接读取该文件。

* 修改文件：当容器需要修改一个文件时（包括更改文件的元数据，如更改文件权限等），按照自上而下的顺序检查镜像的各个层级（包括容器层）。如果**文件位于容器层，则直接修改容器层文件**。如果**文件位于镜像层中，Docker会采取写时复制策略**，将需要修改的文件从镜像层复制到容器层，修改容器层中的副本文件，而原始镜像层中的文件**不会被修改**。
* 添加文件：当在容器中添加文件时，只需在容器层中添加该文件即可，**不会影响原始镜像层中的文件**
* 删除文件：当在容器中删除文件时，首先为**容器层中文件添加标记，表明文件已删除，但不会立即释放**。如果该文件存在位于镜像层中的原始文件，则仅标记容器层中的文件，而不影响原始文件。当容器关闭或重新创建镜像时标记删除的文件才会被清理，从而释放底层存储空间。



### 2.4.3 OverlayFS存储驱动程序

#### 2.4.3.1 overlay2工作原理

![overlay2](.\images-docker\overlay2.png)

OverlayFS在Linux主机上**以两个目录存在，通过一个目录进行呈现**。这两个目录被称为“层”，而将这两个目录合并显示为一个目录的过程则被称为“联合挂载”。OverlayFS将底层目录称为`lowerdir`，将顶层目录称为`upperdir`。统一视图（unified view）则通过一个名为merged的目录展示出来。

* lowerdir：基础层(通常是一个镜像层)，能被上层目录(upperdir)共享的只读层，包含了文件系统的初始状态。（即**Docker结构中的Image layer镜像层**)
* upperdir：容器的可写层，用于存储容器中的修改。容器的写操作（包括创建、修改、删除文件）时均在该层进行，确保容器的文件系统隔离。（即**Docker结构中的可读写层read-write layer**）
* workdir：为了确保执行写操作期间，容器层的一致性和可靠性，引入了临时处理写操作的层。当容器进行写操作时，首先将文件**从upperdir中复制到workdir**(Cow策略)，在workdir中完成实际的写操作，更改后将文件或目录复制回upperdir，实现容器中的文件系统的修改。
* mergeddir：将lowerdir、upperdir联合挂载到merged目录，提供完整的“统一视图”。

![compareDockerOverlay](.\images-docker\compareDockerOverlay.png)

#### 2.4.3.2 overlay2的Cow策略

* 读取文件：

  * 若文件**不存在于容器层**（upperdir）中，则直接读取镜像层（lowerdir）中的文件。
  * 若文件**只存在于容器层**而不存在于镜像层中，则直接读取容器层中的文件。
  * 若文件**同时存在于容器层和镜像层**，则读取该文件在容器层的版本。

* 修改文件或目录：
  * 首次写入文件：容器首次写入现有文件时，该文件在容器层（upperdir）中不存在。overlay2存储驱动程序将执行copy_up操作，将文件从镜像层（lowerdir）复制到容器层，并更改容器层中文件副本的内容。
  * 删除文件：当在容器层（upperdir）中删除文件时，实际上不会直接从镜像层（lowerdir）中删除文件（因为镜像层是只读的），而是在容器层中创建一个特殊的“白色标记文件”（whiteout file），标识该文件已被删除。虽然镜像层中文件仍然存在，但由于白色标记文件的存在，统一视图中该文件将不再可见。
  * 删除目录：当在容器层（upperdir）中删除目录时，会在容器层中创建一个“不透明目录”（opaque directory)，与“白色标记文件”思想类似。不透明目录会在统一视图中隐藏该目录，即便镜像层中的目录仍然存在。



## 2.5 Docker容器数据存储

> 关于挂载：
>
> 在Linux系统中，所有都是文件，且都起源于根目录，不像windows做好了硬盘分区如C: D:盘。系统对设备的访问不是直接访问，而是将设备挂载到一个目录，再通过该目录访问、读写该设备。windows的硬盘分区实际上就是已挂载的可视化目录。而Linux没有硬盘分区，但每个设备都挂载了相应的目录，这各个目录也可以看作不同分区。

### 2.5.1 数据卷(Volumes)

数据卷是由Docker管理的特殊目录，存储在Docker主机的指定位置，通常是 `/var/lib/docker/volumes`（对于Linux主机）。数据卷的挂载方式会自动将容器内对应的挂载目录内容填充到数据卷中。

这意味着，将一个数据卷挂载到容器时，数据卷将被容器内的文件和目录内容初始化。官方强烈推荐将数据卷作为数据存储的**首选方式**，并且**非Docker进程**不应修改数据卷的内容。这种设计使得数据卷具有更高的可靠性和可移植性，同时确保了与容器的解耦。

数据卷内存的数据是以`json`的形式保存的

#### 2.5.1.1 使用场景

* 当 Docker 宿主机（host）上的目录结构或文件结构不确定时，使用数据卷（volumes）可以将宿主机的配置与容器运行时解耦。
* 当存在需要进行Docker主机备份、恢复或者迁移到另外一台Docker主机时，volumes是更好的选择。只需停止使用该卷的容器，然后备份该卷的目录（如`/var/lib/docker/volumes/<volume-name>`）即可完成。

#### 2.5.1.2 使用步骤

`docker volume create [VOLUME_NAME]`使用该命令创建指定名称的数据卷，若未指定名称，则该卷为匿名卷。

`docker volume inspect [VOLUME_NAME]`使用该命令查看数据卷信息

`-v [VOLUME_NAME]:[target-directory]`，将指定数据卷与目标目录挂载，挂载后数据卷即可同步目标的数据。



### 2.5.2 绑定挂载(Bind Mounts)

绑定挂载是特殊的挂载，一般的挂载是将设备与目录挂载，而绑定挂载可以将一个目录（或文件）挂载到另外一个目录（或文件），挂载后，可从挂载点目录（target）访问源目录。

#### 2.5.2.1 使用场景

* 当需要将文件从主机共享到容器时，例如应用的自定义配配置。
* 当需要在 Docker 主机上的开发环境和容器之间共享源代码或构建工件时。例如，您可以将 Maven  `target/ `目录挂载到容器中，每次在 Docker 主机上构建 Maven 项目时，容器都可以访问重新构建的制品。
* 当保证Docker主机的文件或目录结构与容器所需的绑定挂载一致时

#### 2.5.2.2 使用步骤

`-v [origin-directory]:[target-directory]`使用该命令将源文件挂载到目标目录，也可以使用`--mount bind [origin-directory]:[target-directory]`命令实现同样的效果



## 2.6 Dockerfile

Dockerfile 是一个静态的文本文件，其中按照顺序记录了构建特定镜像所需的所有指令，**Docker引擎通过解析Dockerfile中的指令来自动化构建镜像**。这种自动化构建的方式使得镜像的创建过程变得可重复和可靠

在Dockerfile中，可以指定所需的基础镜像、安装软件包、拷贝文件、设置环境变量等等。**每个指令都会在前一个指令的基础上进行构建，最终形成一个完整的镜像**。**通过编写清晰的Dockerfile，可以定义和管理镜像的构建过程**，确保镜像可以在不同的环境中正确地部署和运行。

当运行`docker build`命令时，Docker会读取并执行Dockerfile中的指令，并将其转化为一个镜像。**这个镜像可以用来创建和运行多个相同配置的容器实例**。因此，Dockerfile提供了一种标准化和可重复的方式来构建容器化应用程序的环境。

Dockerfile遵循特定的格式：

* 指令忽略大小写，推荐使用大写
* 每条指令后至少有一个参数
* #表示注释

Dockerfile 的指令逻辑通常遵循以下模式：

* 选择合适的基础镜像
* 安装基础工具与依赖
* 添加其他应用
* 清理缓存
* 声明镜像端口暴露情况
* 设置默认启动命令



### 2.6.1 Dockerfile 优化

Dockerfile优化是指针对Docker镜像构建过程中的性能、安全性、可维护性等方面进行的改进。主要可分为**镜像大小优化和镜像构建速度**优化两个方向的优化。

镜像大小优化方式主要包括：

* 减少镜像层层数
* 删除非必要的包和缓存
* 使用多阶段构建

镜像构建速度优化主要包括：

* 使用`.dockerignore`对文件进行忽略
* 充分使用缓存镜像层

#### 2.6.1.1 减少镜像层数

Docker镜像采用分层存储方式，当上层需要修改文件时，通过CoW策略复制一份下层文件并进行修改，并且下层文件始终存在，因此当镜像层数越多，镜像体积越大。

在当前Docker版本(24.0.5)中只有RUN、COPY、ADD三个指令会新引入镜像层，其他指令只创建临时中间镜像，并不会增加构建的大小。

最常用的优化方式为将多个RUN指令合并为一个RUN指令。

#### 2.6.1.2 删除非必要的包和缓存

优化镜像大小，可以考虑删除包管理工具（如apt、yum）的缓存。在安装软件包时，包管理工具会将下载的软件包及其相关文件存储在缓存目录中。这些缓存文件不再需要，可以在构建过程中清理掉。

删除通过 COPY、ADD 添加的压缩包：在构建镜像时，如果使用了 COPY 或 ADD 命令将压缩包添加到镜像中，建议在解压后删除这些原始压缩包。



# 3. Docker操作命令

[Docker常用命令大全 - 简书 (jianshu.com)](https://www.jianshu.com/p/a84e8cf33b34)

## 3.1 容器镜像搜索

`docker search TERM`命令用于搜索镜像别难过获取基本信息

* 若不希望描述字段截断输出，可加上--no-trunc参数
* 搜索结果中未包含容器镜像的tag信息，若需要下载指定tag的镜像，可访问(https://hub.docker.com)，查询mysql以获取更多容器镜像信息

## 3.2 容器镜像拉取

`docker pull NAME[:TAG]`命令从Registry下载指定名称的容器镜像

TAG用于指定下载的容器镜像版本。若为空，则默认下载latest版本

使用`docker images`或`docker image ls`命令显示所有的顶层镜像基本信息

## 3.3 基于镜像创建并运行容器

* `dfocker run [OPTIONS] IMAGE` 命令用于基于指定容器镜像创建并运行容器

  * -d ：后台运行容器并打印容器ID
  * --name：指定容器名称

  * -e，--env：为容器设置环境变量

* 查看所有运行中的容器 `docker ps`
  * `-a`参数可以列出所有的容器（包含停止的）
  * `-q`参数标识只列出container id, 而不包含其他信息

## 3.4 停止、启动、进出容器

* 停止容器
  * 使用`docker stop CONTAINER`命令停止一个或多个运行中的容器，并输入对应容器ID（或名称）
  * CONTAINER既可以使用容器ID指定，也可使用容器名称来指定
* 启动容器
  * 维护结束后，需要启动容器以提供服务，可以使用`docker start CONTAINER`命令启动一个或多个已停止的容器
* 进入容器
  * 容器启动后，使用`docker exec`命令进入容器中验证mysql服务是否正常启动
  * docker exec命令实质上是在运行的容器内执行一条命令
    * 当指定命令为shell程序时，例如/bin/bash，配合-it参数可实现进入容器进行交互式操作
    * -i：保持交互模式
    * -t,--tty：分配虚拟终端

* 退出容器
  * 使用`exit`可退出当前mysql交互终端，再次使用exit可退出当前mysql容器

## 3.5 容器日志查看

`docker logs`命令查看容器的日志输出



# 4. 总结

Docker的运行流程主要包括创建与启动容器、进入容器等步骤。

以下是关于Docker运行流程的简要说明以及原理的解释：

1. 创建与启动容器：
   - 使用`docker run`命令来创建并启动一个新的容器。
   - 可以通过`-i`和`-t`参数为容器分配一个伪终端，这样用户可以交互式地与容器进行交互。
   - 使用`--name`参数可以为容器指定一个名称。
   - `-v`参数用于映射宿主机目录和容器目录，实现数据的共享和持久化。
   - 镜像可以被视为容器的原型，当执行`docker run`命令时，Docker会基于指定的镜像来创建容器。
2. 进入容器：
   - 当容器启动后，可以使用`docker exec`命令进入容器的命令行界面。
   - 这允许用户在容器内部执行命令和操作。

Docker的运行原理主要基于以下几个核心概念：

- **容器**：Docker容器是一个轻量级的运行时环境，它包含了应用程序及其依赖项，并与宿主机的其他部分隔离开来。容器是从镜像创建的，并且共享宿主机的内核。
- **镜像**：Docker镜像是容器的静态表示，它包含了应用程序的代码、运行时环境、库、配置文件等。镜像可以被视为一个只读的模板，用于创建容器。
- **Docker Daemon**：Docker守护进程是Docker架构中的核心组件，它负责监听API请求，管理容器的生命周期（创建、启动、停止等），以及管理镜像和网络。
- **文件系统隔离**：Docker使用AUFS（Advanced Multi-Layered Unification Filesystem）或其他类似的文件系统技术来实现容器之间的文件系统隔离。每个容器都有自己的文件系统视图，这保证了容器之间的独立性。
- **网络隔离**：Docker提供了多种网络模式，允许容器之间进行通信，以及与宿主机或其他Docker主机之间的通信。这确保了容器网络环境的隔离性和安全性。

综上所述，Docker通过创建容器来运行应用程序，这些容器是从镜像创建的，并且与宿主机和其他容器隔离开来。Docker守护进程负责管理容器的生命周期和镜像，而用户可以通过Docker客户端与守护进程进行交互，执行各种容器管理操作。
