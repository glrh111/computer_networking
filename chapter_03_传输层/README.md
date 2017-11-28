+ 两个实体怎么在一种会丢失数据的传输没接上可靠通信
+ 控制传输层实体的传输速率以避免网络拥塞

## 3.1 概述和传输层服务

逻辑通信 logic communication: 通过逻辑通信，运行不同进程的主机好像直接相连一样；实际上，这些主机可能相隔万里，中间有很多路由器和传输链路。

![传输层提供的逻辑通信](http://o9hjg7h8u.bkt.clouddn.com/%E4%BC%A0%E8%BE%93%E5%B1%82%E9%80%BB%E8%BE%91%E9%80%9A%E4%BF%A1.png)

传输层协议是在end system中实现，而不是在路由器中实现。

### 3.1.1 传输层和网络层的关系

传输层为运行在不同主机的进程之间提供了逻辑通信；网络层提供了主机之间的网络通信。

传输层只工作在end system中，负责将message移动到网络层中。

### 3.1.2 传输层概述

IP的服务模型是尽力而为交付服务best-effort delivery service，但是不做任何确保，
+ 不确保segment的交付
+ 不保证segment的按序交付
+ 不保证segment数据的完整性

所以，IP被称为不可靠服务unreliable service.

UDP和TCP最基本的责任是，将两个end system间IP的交付服务扩展为运行在end system上的两个进程之间的交付服务。
将主机间的交付扩展到进程间的交付，被称为传输层的多路复用transport-layer multiplexing, 与多路分解demultiplexing

进程到进程的数据交付和差错检查是两种最低限度的传输层服务，也是UDP能提供的仅有的两项服务。

TCP提供了几种附加服务
+ 可靠数据传输 reliable data transfer 通过使用流量控制，序号，确认和定时器，TCP确保正确的，按序地将数据从发送进程交付给接收进程
+ 拥塞控制 congestion control 力求为每一条拥塞网络链路的连接平等地共享网络连接带宽。

## 3.2 多路复用和多路分解

+ 多路分解demultiplexing 将传输层segment中的数据交付到正确的socket的工作
+ 多路复用multiplexing 在源主机中从不同socket手机数据块，并为每个数据块封装上首部信息，从而生成segment，然后将segment传递到网络层

多路复用的要求：
+ socket有唯一标识符
+ 每个segment有特殊字段来指示该segment所要交付到的socket：源端口号字段source port number field和目的端口号字段destination port number field

端口号是一个16bit的数，大小在0~65535, 0~1023为周知端口号well known port number.

### 1 无连接的多路复用与多路分解

一个UDPsocket可以由一个二元组全面标识，只要具有相同的目的地址和目的端口号，那么segment将通过相同的目的socket定向到相同的目的进程。

### 2 面向连接的多路复用与多路分解

TCPsocket是由一个四元组标识的，包括源地址，源端口号，目的地址，目的端口号，当segment到达一台主机时，全部4个值被用来将segment定向到对应的socket。
除非TCP segment携带了初始创建连接的请求。

![TCP多路复用与多路分解](http://o9hjg7h8u.bkt.clouddn.com/TCP%E5%A4%9A%E8%B7%AF%E5%A4%8D%E7%94%A8%E4%B8%8E%E5%A4%9A%E8%B7%AF%E5%88%86%E8%A7%A3.png)

### 3 WEB服务器与TCP

## 3.3 无连接传输UDP

选择UDP可能基于一下原因：
+ 关于何时、发送什么数据的应用层控制更加精细：因为TCP有拥塞控制，而且必须等接收方确认才会停止重发。
+ 无需建立连接
+ 无连接状态
+ packet首部开销小:每个TCP segment有20bytes的首部开销；UDP只有8个。

应用例子：
+ RIP路由选择表的更新：每隔5分钟更新一次
+ SNMP
+ DNS
+ NFS 远程文件服务器
+ 流媒体或者因特网电话的某些场合

UDP没有拥塞控制，导致较高的丢包率。

使用UDP可以通过在应用层建立某些机制，来保证数据的可靠传输，而无需受制于拥塞控制。

### 3.3.1 UDP segment 结构

RFC 768

首部一共4个字段，8字节
+ 源端口号，目的端口号：多路分解、复用
+ UDP segment的字节数：首部加数据
+ 校验和 checksum
+ 应用数据

### 3.3.2 UDP checksum

RFC 1071

提供了差错检测功能。检验和用于确定当UDP segment从源到达目的地移动时，其中的byte是否改变（例如，由于链路层噪音干扰或者存储在路由器中时引入问题）。
发送方的UDP对segment中所有16byte的和进行反码运算，求和时遇到的任何溢出都被回卷。得到的结果被放在UDP相关字段中。

目的方接收时，所有的字段加起来应该全为1，如果有0，那么传输过程中出差错了。

UDP为什么要提供差错检测？
+ 不是所有的链路层协议都提供了差错检测，不能靠这个
+ 数据脱离链路层到传输层这个阶段，同样有可能引起差错

虽然UDP提供差错检测，但是它对差错恢复无能为力，它可以选择
+ 丢弃受损的segment
+ 交给应用程序并提出警告


## 3.4 可靠数据传输原理 reliable data transfer

![可靠数据传输](http://o9hjg7h8u.bkt.clouddn.com/%E5%8F%AF%E9%9D%A0%E6%95%B0%E6%8D%AE%E4%BC%A0%E8%BE%93rdt.png)

### 3.4.1 构造可靠数据传输协议 从158(138)页开始

#### 1. 经完全可靠信道的可靠数据传输 rdt1.0



#### 2. 经具有比特差错信道的可靠数据传输 rdt2.0
 
以日常打电话为例，接听方在听到一句话时，可以
+ 肯定确认 positive acknowledgement
+ 否定确认 negative acknowledgement

在计算机网络环境中，基于这样的重传机制的可靠数据传输协议称为自动重传请求automatic repeat reQuest ARQ
+ 差错检测：时接收方可以检测到何时出现了byte差错
+ 接收方反馈：ACK, NAK, 让发送方知道接收方情况。
+ 重传

但是ACK，NAK packet受损怎么办？
+ 考虑打电话场景，发送方将问：你说啥？但是这个你说啥 受损了怎么办？来回重复无穷尽也
+ 增加足够的检验和byte：使发送方不仅可以检测差错，还可恢复差错。对于会产生差错但不会丢失packet的信道，可以
+ 直接重传。但是会在信道中引入冗余分组dupliate packet, 搞得接收方不知道上次发的ACK，NAK是否被发送方正确收到。
因此它无法事先知道，收到的分组是一次新的，还是重传。

解决这个问题的一个方法是，在数据packet中添加一个新字段，于是，接收方只需要检查序号即可确定收到的packet是否是一次重传。

再次改进版本是，接收方可以通过对同一上一个sequence number发送两次ack，来通知发送方没有收到下一个packet。

#### 3. 经具有比特差错的丢包信道的可靠数据传输 rdt3.0

怎样检测丢包，发生丢包后该做什么？

如果一个发送方，发出去的packet未被收到，或者ack丢失，那么发送方需要等待足够长的时间，以决定重传packet。

如果因为网络阻塞，导致发送方经历很久才收到ack，此时重发的packet称为冗余数据packet duplicate data packet.

基于时间的重传机制，需要倒计时定时器countdown timer, 发送方需要做到：
+ 每次发送一个packet，启动一个定时器
+ 响应定时器中断
+ 终止定时器

至此，我们得到了一个可靠的数据传输协议

### 3.4.2 流水线可靠数据传输协议

rdt3.0 是一个功能正确的协议，但是性能不咋地。核心在与它是一个停等协议。

发送方(或信道)利用率：传输时延（发送方实际忙于将发送比特送进信道的那部分时间）/ 发送时间。

U = (L/R) / (RTT+L/R) 停等协议的u非常低。

![停等协议与流水线协议](http://o9hjg7h8u.bkt.clouddn.com/%E5%81%9C%E7%AD%89%E5%8D%8F%E8%AE%AE%E4%B8%8E%E6%B5%81%E6%B0%B4%E7%BA%BF%E5%8D%8F%E8%AE%AE.png)

流水线协议对可靠数据传输协议可带来如下影响：
+ 必须增加序号范围，每个传送中的packet必须有一个唯一的序号
+ 协议的发送和接收两端也必须缓冲多个packet。发送方最少应该缓冲那些已经发送但是没有ack的packet；
接收方或许应该缓存已经正确接受的packet
+ 所需序号范围和对缓冲的要求取决于数据传输协议如何处理丢失、损坏及延时过大的packet。解决流水线差错恢复有两种基本方法
  + 回退N步 Go-Back N
  + 选择重传 Selective Repeat SR

### 3.4.3 回退N步GBN Page167

如果某个ack没有收到，那么发送从该ack之后的所有packet，容易产生较多的冗余packet。

### 3.4.4 选择重传 SR

可靠数据传输机制及其用途的总结
+ 检验和 用于检验在一个传输packet中的比特错误
+ 定时器 用于超时/重传一个packet，可能因为该packet或者它的ack在信道中丢失了。
+ 序号   用于为从发送方流向接收方的数据packet按顺序编号。
+ 确认   接收方用于通知发送方某个packet或一组packet已经被正确接收到。确认ack里边可以包含多个packet编号
+ 否定确认 接收方用于告诉发送方,某个packet未被正确接收，否定接收报文通常携带未被正确接收的packet的序号
+ 窗口、流水线 发送方也许被限制仅发送那些序号落在一个指定范围内的packet。这样就允许发送方一次发送多个packet。
窗口长度可以根据接收方接收和缓存报文的能力、网络中的拥塞程度，或两者情况来进行设置。



## 3.5 面向连接的传输 TCP

TCP依赖于差错检测、重传、累积确认、定时器、用于序号和确认号的首部字段。RFC 793, RFC 1122, RFC 1323, RFC 2018, RFC 2581

### 3.5.1 TCP连接

TCP被称为是面向连接的 connection-oriented , 这是因为在一个应用进程可以开始向另一个应用进程发送数据之前，
这两个进程必须先相互握手，即他们必须先发送某些预备报文段，以建立确保数据传输的参数。作为TCP连接建立的一部分，连接的双方都将初始化与TCP连接相关的许多TCP状态变量。
 
这种TCP连接不是一条像在电路交换网络中的端到端TDM、FDM电路，也不是一条虚电路，因为其两个状态完全保留在端系统中。
由于TCP协议只在端系统中运行，而不是在中间的网络元素（路由器和链路层交换机）中运行，所以中间的网络元素不会维持TCP连接状态。
事实上，中间的路由器对TCP完全视而不见，他们看到的是数据报，而不是连接。

+ TCP连接提供的是全双工服务 full-duplex service。
+ TCP连接是点对点point-to-point的，即在单个发送方与单个接收方间的连接。

三次握手 three-way handshake
+ 客户发送一个TCP segment
+ 服务器用另一个特殊的TCP segment来响应
+ 客户端再用第三个特殊segment响应

前两个segment不承载有效载荷，不包含应用层数据，第三个segment可以承载有效载荷。

![TCP发送缓存和接收缓存](http://o9hjg7h8u.bkt.clouddn.com/tcp%E5%8F%91%E9%80%81%E7%BC%93%E5%AD%98%E5%92%8C%E6%8E%A5%E6%94%B6%E7%BC%93%E5%AD%98.png)

一旦建立起一条TCP连接，两个应用程序之间就可以相互发送数据了
+ 客户端进程通过socket传递数据流，之后由TCP控制
+ TCP将这些数据引导到该连接的发送缓存send buffer里，三次握手初期设置的缓存之一，接下来TCP就会不时地从其中取出一条数据。
TCP可从缓存中取出并放入segment中的数据数量受限于MSS Maximum Segment Size最大报文段长度。
MSS通常根据最初确定的由本地之际发送的最大链路层帧长度，即最大传输单元MTU Maximum Transmission Unit来设置。
设置该MSS要保证一个TCP segment加上TCP/IP首部长度（通常40bytes）将适合单个链路层帧。以太网和PPP链路层协议都具有1500字节的MTU，
因此MSS的典型值是1460字节。
已经提出了多种发现路径MTU的方法，并基于路径MTU（从源到目的地的所有链路上发送的最大链路层帧）设置MSS。
注意，MSS是segment里应用层数据的最大长度。
+ TCP为每一个segment配上TCP首部，这些segment被下传给网络层，网络层将其分别封装在IP datagram中。
+ 当TCP在另一端接收到一个segment后，该segment的数据就被放入该TCP连接的接收缓存中。应用程序从此缓存中读取数据流。

### 3.5.2 TCP segment 结构

![TCP segment结构](http://o9hjg7h8u.bkt.clouddn.com/TCP_segment%E7%BB%93%E6%9E%84.png)

TCP segment由首部字段和数据字段组成。MSS限制了segment数据字段的最大长度。当TCP发送一个大文件，TCP通常将该文件划分成长度是MSS的若干块，
最后一块通常小于MSS。交互式应用通常传送长度小于MSS的数据，Telnet通常只有一个字节的数据。

TCP首部一般是20字节。

TCP标志位含义
+ SYN synchronous 建立联机
+ ACK acknowledgement 确认
+ PSH push 传送
+ FIN finish 结束
+ RST reset 重置
+ URG urgent 紧急

#### 1 序号和确认号

TCP把数据看做无结构的、有序的字节流。序号是建立在传送的字节流之上，而不是建立在传送的segment的序列之上。一个segment的编号sequence number for a segment
因此是该segment首字节的字节流编号（字节的序号）。

![文件数据划分为TCP segment](http://o9hjg7h8u.bkt.clouddn.com/%E6%96%87%E4%BB%B6%E6%95%B0%E6%8D%AE%E5%88%92%E5%88%86%E4%B8%BATCP%20segment.png)

主机A从主机B接收数据之后，所发送的确认号，是主机A期望从主机B收到的下一个字节的序号。

TCP连接的双方，均可以随机地选择初始序号，这样做可以减少将那些仍在网络中存在的两台主机之间先前已终止的连接segment,。。。

### 3.5.3 往返时间的估计与超时

TCP采用超时/重传来处理segment的丢失问题。

#### 1 估计往返时间  page 181

segment的样本RTT，即Sample RTT, 就是从某segment被发出到对该segment的确认被收到之间的时间量。
S RTT也是变化的，通过以下公式取平均值：EstimatedRTT = (1-alpha) * EstimatedRTT + alpha * SampleRTT.

alpha的一个RFC 6298的建议值是0.125.

RTT偏差 DevRTT = (1 - beta) * DevRTT + beta * | SampleRTT - EstimatedRTT |

beta的推荐值是0.25.

#### 2 设置和管理重传超时间隔

TimeoutInterval = EstimatedRTT + 4 * DevRTT

推荐的初始TimeoutInterval值为1秒，出现超时后，会加倍，一旦更新EstimatedRTT, 又重新计算。

### 3.5.4 可靠数据传输

推荐的定时器管理过程仅使用单一的重传定时器。即时有多个已发送但未被确认的segment。

处理三个事件：收到应用层的数据，收到ack，超时

一些额外的话题：
+ 超时间隔加倍
+ 快速重传fast retransmit：接收方通过发送冗余ack，来通知发送方某个segment需要重传
+ GBN还是SR？TCP使用他们的混合体

### 3.5.5 流量控制

TCP为它的应用程序提供了流量控制服务 flow-control service, 以清除发送方使接收方缓存溢出的可能性。
因此，流量控制是一个速度匹配服务，即发送方的发送速率与接收方应用的读取速率相匹配。

TCP发送方也可能因为IP网络的拥塞而被遏制，这种形式的发送方的控制称为拥塞控制congestion control，它和流量控制是针对不同原因而采取的措施。

发送方维护一个称为接收窗口receive window的变量，来进行流量控制。代表接收方还有多少接收缓存可用。

![接收缓存和接收窗口](http://o9hjg7h8u.bkt.clouddn.com/%E6%8E%A5%E6%94%B6%E7%AA%97%E5%8F%A3%E5%92%8C%E6%8E%A5%E6%94%B6%E7%BC%93%E5%AD%98.png)

### 3.5.6 连接管理

![TCP三次握手](http://o9hjg7h8u.bkt.clouddn.com/TCP%E4%B8%89%E6%AC%A1%E6%8F%A1%E6%89%8Bsegment%E4%BA%A4%E6%8D%A2.png)

参与一条TCP连接的两个进程中的任何一个都能终止该连接。

![关闭一条TCP连接](http://o9hjg7h8u.bkt.clouddn.com/%E5%85%B3%E9%97%AD%E4%B8%80%E6%9D%A1TCP%E8%BF%9E%E6%8E%A5.png)

SYN洪范攻击 SYN flood attack：攻击者发送SYN，服务器分配资源并响应SYNACK，如果客户端不发送ACK来完成第三次握手，通常在1分钟多后，资源将被回收。
一种有效的防御是SYN cookie.
+ 服务端收到SYN时，不分配资源，此时根据源端口号和IP生成一个cookie，作为初始序列号放在SYNACK作为响应。服务器不存储该cookie的值
+ 如果客户合法，返回ACK，此ACK=cookie+1, 那么服务器此时分配资源
+ 如果客户端没有发送ack，那么服务器也不受其他损失。

如果服务器不接受某个端口的TCP连接，那么会在响应RST segment；如果收到的是UDP数据，那么返回一段特殊的ICMP。

例如，nmap向目标主机的某个端口发送一个特殊的TCP SYN segment，有三种可能的输出：
+ 收到源主机一个 TCP SYNACK segment，意味着该端口打开
+ 收到一个TCP RST segment，表示segment到达，但是端口不可用；segment没有被防火墙阻挡
+ 啥都没收到 无法到达目的主机


## 3.6 拥塞控制原理



































