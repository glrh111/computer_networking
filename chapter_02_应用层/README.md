网络应用是计算机网络存在的理由。

## 2.1 应用层协议原理

### 2.1.1 网络应用程序体系结构

应用程序的体系结构，明显不同于网络的体系结构。

从应用程序研发者的角度来看，网络体系结构是固定的，并为应用程序提供了特定的服务集合。
另一方面，应用程序体系结构application architecture 由应用程序研发者设计，
规定了如何在各种端系统上组织该应用程序。

在选择应用程序体系结构时，应用程序研发者很有可能利用现代网络程序中所使用的两种主流体系结构之一，
客户-服务器体系结构，P2P体系结构（适合流量密集型应用）。

P2PAPP设计的三个挑战
+ ISP友好：起源于对ISP提供的下行带宽比上行带宽的不对等
+ 安全性：高度分布和开放特定带来的挑战
+ 激励：怎么说服用户自愿向应用提供带宽、存储和计算资源。


### 2.1.2 进程通信

在不同端系统end system上的进程，通过跨越计算机网络交换message而互相通信。

+ 客户和服务器进程：通信的两端，发起通信的进程被标识为client，等待联系的进程是server。
+ 进程与计算机网络之间的接口：进程通过被称为socket套接字的软件接口向网络发送和接受message。

![socket通信示意](http://o9hjg7h8u.bkt.clouddn.com/socket%E9%80%9A%E4%BF%A1.png)

socket是同一台主机内应用层和传输层的接口，由于该应用程序是建立网络应用程序的可编程接口，
因此socket也称为应用程序和网络之间的API。应用程序开发者可以控制socket在应用层的一切，但是对该socket的
传输层端几乎没有控制权。对传输层的控制权仅限于：
+ 选择传输层协议
+ 设定最大缓存，最大segment长度等参数

为了标识接收目的进程，需要定义：
+ 主机的地址
+ 定义在目的主机中的接收进程的标识符 `https://www.iana.org/`

主机由IP地址标识，目的端口号port number可以标识目的进程

### 2.1.3 可供应用程序使用的传输服务

+ 可靠数据传输：packet可能在网络可能丢失：在路由器的缓存中溢出，bit损坏后被丢失等。
如果一个协议提供了能够将数据正确、完全地交付给另一端的服务，就认为提供了可靠数据传输reliable data transfer.
+ 吞吐量：具有吞吐量要求的应用程序，被称为带宽敏感应用bandwidth-sensitive application. 许多多媒体应用是这样的；与之相反的是
弹性应用elastic application, 如电子邮件，文件传输
+ 定时：这样的例子比如：发送方注入socket的每个byte到达接收方的socket不迟于100ms。
+ 安全性：传输协议能够为app提供一种或者多种安全性服务。比如加密数据，完整性，端点鉴别等。

### 2.1.4 因特网提供的传输服务

#### TCP服务

TCP服务模型包括面向连接服务和可靠数据传输服务
+ 面向连接的服务：当应用层message流通之前，TCP让C/S交换运输层控制信息，成为握手。之后，一个TCP连接在两个进程的socket之间建立了。
这条连接是全双工的，即双方进程可以在此连接上同时进行message收发。当app结束发送message之后，必须删除该连接。
+ 可靠的数据传输服务：通信进程能够依靠TCP，无差错，按适当顺序交付所有发送的数据。

TCP还具有拥塞控制机制，不一定能为通信进程带来直接好处，但是对整个互联网非常有好处。
+ 当网络出现拥塞时，会抑制发送程序
+ 试图限制每一个TCP连接，使他们达到公平共享网络带宽的目的。

SSL Secure Sockets Layer 是一种在应用层上对TCP的加强，提供了关键的进程到进程的安全性服务，包括加密，数据完整性，端点鉴别等。

#### UDP服务

不提供不必要服务的轻量级传输协议。无连接的，没有握手过程。不能保证全部message可达，也不能保证顺序到达。

#### 传输协议不能提供的服务

没有提供对定时，和安全性的保障。

### 2.1.5 应用层协议

如何构造这些报文？各个字段的含义是神马？进程合适发送这些报文？这是应用层的范畴。
+ 交换的message类型：例如请求报文和响应报文
+ 各种message类型的语法：报文中的各个字段，和这些字段如何描述
+ 字段的语义：一个进程何时以及如何发送message，对message进行相应的规则

应用层协议只是应用层的一部分。

### 2.1.6 本书设计的网络应用

+ Web
+ 文件传输
+ 电子邮件
+ 目录服务
+ P2P

## 2.2 Web和HTTP

Web是一个引起公众注意的因特网应用。也将因特网从只是很多数据网之一的地位提升为仅有的一个数据网。

Web的按需操作，不同于电视播放的节目。

### 2.2.1 HTTP概况

可以看到分层体系带来的最大的优点：HTTP协议不用担心数据丢失，也不关注TCP从网络的数据丢失和乱序中恢复的细节。

HTTP是一个无状态协议。stateless protocol

### 2.2.2 非持续连接和持续连接

#### 采用非持续连接的HTTP non-persistent connection

每个TCP连接在服务器发送一个对象后关闭，该连接并不为其他的对象而持续下来。

RTT Round-Trip Time 往返时间 一个短packet从客户端到服务器再回到客户端所花费的时间。

请求一个文件所需要的时间，是 2 * RTT + 文件传输时间, 第一个RTT内是握手的前两次，第三次握手会发送HTTP请求，然后就是传输对象文件。

![HTTP响应](http://o9hjg7h8u.bkt.clouddn.com/http%E7%9B%B8%E5%BA%94rtt.png)

#### 采用持续连接的HTTP persistent connection

非持续连接的缺点
+ 每个请求需要一个全新的连接，为客户端和服务器增加了开销
+ 每个对象经受两倍RTT的时延，一个用于创建TCP，另一个用于请求接收对象

### 2.2.3 HTTP报文格式

请求报文和响应报文

请求行 request line GET /wocao/nidaye HTTP/1.1
首部行 header line ...
实体   entity body

![HTTP请求报文](http://o9hjg7h8u.bkt.clouddn.com/HTTP%E8%AF%B7%E6%B1%82%E6%8A%A5%E6%96%87.png)

状态行 status line 
首部行 header line 
实体   entity body

![HTTP响应报文](http://o9hjg7h8u.bkt.clouddn.com/HTTP%E5%93%8D%E5%BA%94.png)

## 2.3 文件传输协议 FTP

FTP使用了两条并行的TCP连接来传输文件，控制连接control connection 和数据连接data connection。
+ 控制连接用于在两主机间传输控制信息：用户标识，口令，改变远程目录的命令等 使用21端口
+ 实际发送一个文件 使用20端口

所以说，FTP的控制信息是带外传送的，out-of-band; HTTP是带内in-band发送控制信息的。

## 2.4 电子邮件

电子邮件系统包含：用户代理 user agent，邮件服务器 mail server，简单邮件传输协议 Simple Mail Transfer Protocol

![SMTP](http://o9hjg7h8u.bkt.clouddn.com/SMTP.png)

两台SMTP服务器直连，没有中间服务器。 使用25端口

与HTTP的对比：
+ HTTP是一个pull protocol，即用户主动从服务器获取信息；SMTP是一个push protocol, 发送服务器将邮件push给接收服务器。
+ SMTP要求报文使用ASCII编码
+ 处理即包含文本又包含图形的文档，SMTP将他们放到一个报文里边。

### 2.4.3 邮件报文格式和MIME

### 2.4.4 邮件访问协议

POP3 Post Office Protocol Version 3, IMAP Internet Mail Access Protocol 因特网邮件访问协议

为什么需要这些东西？SMTP是一个push portocol，需要发送数据的一方建立连接，接收方取邮件是一个pull操作，所以需要这些协议。

```bash
glrh11@glrh11-ThinkPad-T460:~/Downloads$ telnet pop.163.com 110
Trying 123.125.50.29...
Connected to pop3.163.idns.yeah.net.
Escape character is '^]'.
+OK Welcome to coremail Mail Pop3 Server (163coms[b62aaa251425b4be4eaec4ab4744cf47s])
user ====
+OK core mail
pass ---=
+OK 1798 message(s) [430209247 byte(s)]
list
+OK 1798 430209247
1 12676
2 4264071
3 2371162
4 4442
5 70281522

客户端可以使用list, retr, dele, quit四个命令
```

IMAP可以满足用户的更多需求。比如在远程服务器上简历文件夹等。

基于web的电子邮件：用户代理是浏览器，发送邮件给服务器时，使用HTTP。

## 2.5 DNS 因特网的目录服务

路由器处理定长的标识符比较擅长，但是对人类不友好，所以需要这项服务。

### 2.5.1 DNS提供的服务 Domain Name System

DNS将主机名转换为背后的IP地址：
+ DNS是一个由分层的DNS服务器实现的分布式数据库
+ 一个使得主机能够查询分布式数据库的应用层协议

DNS服务器通常是运行BIND(Berkeley Internet Name Domain)软件的Unix机器。DNS协议运行在UDP上，使用53端口。

还提供以下重要的服务：
+ 主机别名host aliasing：应用程序可以调用DNS来获得主机别名对应的规范主机名，以及主机的IP地址。
+ 邮件服务器别名mail server aliasing: 
+ 负载分配load distribution：可以使用一个IP集合循环来响应DNS请求

### 2.5.1 DNS的工作机理概述

gethostbyname() 获取ip地址

集中DNS设计的问题：
+ 单点故障a single point of failure 将导致整个因特网崩溃
+ 通信容量traffic volume 单个服务器得处理所有查询
+ 远距离的集中式数据库distant centralized database 有些地方查询距离太远
+ 维护maintenance 导致中央数据库特别庞大，更新也太频繁

#### 1 分布式，层次数据库

为了处理扩展性问题，DNS使用了大量的DNS服务器，他们以层次方式组织，并且分布在全世界范围内。没有一台DNS服务器拥有因特网上的所有主机映射。

![DNS服务器部分层次结构](http://o9hjg7h8u.bkt.clouddn.com/dns%E6%9C%8D%E5%8A%A1%E5%99%A8%E5%B1%82%E6%AC%A1%E7%BB%93%E6%9E%84.png)

![DNS根服务器](http://o9hjg7h8u.bkt.clouddn.com/dns.png)

有三种类型的DNS服务器
+ 根DNS服务器：因特网上有13个根DNS服务器，从A到M，大部分位于北美。每台根服务器，是一个冗余的服务器网络，到2011年秋，共有247个根服务器。
+ 顶级域DNS服务器Top-Level Domain TLD：这些服务器负责域名如com，org，net，edu，gov以及所有国家的顶级域名uk，fr等。参见IANA TLD
+ 权威DNS服务器：在因特网上具有公共可访问的每个组织机构，必须提供公共可访问的DNS记录。可以选择自己维护一台权威DNS服务器，
或者付费存储在服务提供商中。

另外还有本地DNS服务器，local DNS server, 如一个大学，居民区通常有一台这样的服务器。通常起DNS代理作用。

![DNS服务器交互](http://o9hjg7h8u.bkt.clouddn.com/dns%E6%9C%8D%E5%8A%A1%E5%99%A8%E4%BA%A4%E4%BA%92.png)

但是TLD服务器并不总是知道主机对应的权威服务器的ip，只是知道某个中间DNS的IP地址。

上述例子采用了递归查询recursive query和迭代查询iterative query.

从理论上讲，任何DNS查询可以是递归的，也可以是迭代的。但通常情况下，从本机到本地DNS服务器是递归查询，本地DNS与其他各层DNS是迭代查询。

#### 2 DNS缓存 caching

在一个请求链中，当某DNS服务器接收一个DNS应答，它能够缓存包含在该应答中的所有信息。
DNS在一段时间后，将丢弃缓存的信息。

### 2.5.3 DNS记录和报文

共同实现DNS分布式数据库的所有DNS服务器存储了资源记录(resource record RR), RR提供了主机名到IP地址的映射。
每个DNS应答包含了一条或多条记录。RR是一个包含如下资源的4元组 (Name, Value, Type, TTL)

TTL是记录的生存时间，它决定了记录应该从缓存中删除的时间，Name和Value的值取决于Type
+ Type=A，那么Name是主机名，Value是该主机名对应的IP地址. (foo.com, 1.1.1.1, A)
+ Type=NS，则Name是个域，如foo.com，Value是个知道如何获得该域中主句IP地址的权威服务器的主机名。(foo.com, dns.foo.com, NS)
+ Type=CNAME，那么Value是别名为Name的主机对应的规范主机名。(foo.com, wocao.nima.com, CNAME)
+ Type=MX, 那么Value是别名为Name的邮件服务器的规范主机名。(foo.com, mail.bar.foo.com, MX)。通过MX记录，一个公司的邮件服务器和其他服务器能够使用相同的别名。

nslookup 可以查询DNS地址。

http://he.net/ 这个网站可以查询到he线路IP对应的信息

本地会默认向cat /etc/resolv.conf 里边设置的地址，作为上连DNS服务器地址。

dig如何使用？Domain Information Groper
+ dig 向默认的上连DNS查询根服务器的NS记录
+ dig @8.8.8.8 www.baidu.com A 向googlednsserver查询百度的A记录
+ 一些常用的选项
  + -c 设置协议类型 包括IN，CH，HS
  + -f 从文件内容执行批量查询
  + -4 -6 使用哪种IP协议
  + -t 查询类型 A MX
  + -q 显式设置要查询的域名
  + -x 逆向查询
+ 特有的查询选项 +开头的命令
  + TCP代替UDP dig +tcp baidu.com
  + 默认追加域 +domain
  + 跟踪dig全过程 +trace 从根域查询一直跟踪到最终结果
  + 精简输出 +nocmd +short +nocomment + nostat

如何在DNS服务器中插入记录？
向注册登录机构register 验证。ICANN 向各处的注册机构授权。授权机构列表：www.internic.net

## 2.6 P2P应用

文件分发，bittorrent协议；分布在大型对等方社区中的数据库，分布式hash表。

### 2.6.1 文件分发

![P2P文件分发时间](http://o9hjg7h8u.bkt.clouddn.com/p2p%E6%96%87%E4%BB%B6%E5%88%86%E5%8F%91%E6%97%B6%E9%97%B4.png)

![BitTorrent分发文件](http://o9hjg7h8u.bkt.clouddn.com/bittorrent%E5%88%86%E5%8F%91%E6%96%87%E4%BB%B6.png)

参与一个特定文件分发的所有对等方的集合被称为torrent洪流，在一个torrent每个对等方下载等长度的文件块。
每个torrent具有一个基础设施节点叫做tracker追踪器。当一个对等方加入torrent时，它注册自己，并周期性地通知追踪器。

为了决定一个对等方响应哪些请求，bittorrent采用了对换算法，以最高速率给自己提供数据的对等方，具有优先权。

### 2.6.2 分布式散列表 

考虑在P2P网络中，实现一个最简单的数据库。

Distributed Hash Table DHT

+ 使用散列函数将key映射为2^n-1范围内的整数，后者标识网络中的一个对等方。
+ 最邻近对等方，最邻近后继，用于寻找某个key该存储在那个peer中

对等方的网络结构
+ 环形对等方
+ 对等方扰动

bittorrent 使用Kademlia DHT来产生一个追踪器tracker.

## TCP socket编程

![tcpServer的两个socket](http://o9hjg7h8u.bkt.clouddn.com/tcpserver%E8%BF%9B%E7%A8%8B%E7%9A%84%E4%B8%A4%E4%B8%AA%E5%A5%97%E6%8E%A5%E5%AD%97.png)

监听socket与连接socket




 
 

 
















