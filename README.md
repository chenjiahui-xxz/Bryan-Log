# Bryan-study
Bryan的日志收集系统，很好玩的一个文件输入后收集，就像你在写日志将你写的日志收集起来，可以通过etcd配置，不需要重启服务就能够实时监听到需要收集的那些文件，通过etcd的watch监听。将收集到的信息通过管道发送到kafka里面。然后在起一个消费者服务，将kafka的数据进行消费，存储到elasticsearch中

## 包名介绍
* **ByranLog**：ByranLog为日志收集，会先初始化kafka实例，然后在连接etcd服务。拉去etcd中的配置来创建出需要的读取文件实例，同时还能实时的通过 __etcd.watch__ 监控etcd配置文件的变化，来决定是否新增新的配置读取文件实例，或者是删除配置文件中已经没有的实例。实现不需要ByranLog重启服务而能动态增删改灵活读取所需要的文件日志。将读取的到的文件内容通过go管道发送到kafka中进行持久化。如etcd中的配置 *[{"path":"D:\\studyGo\\s4.log","topic":"web_log"},{"path":"D:\\studyGo\\s5.log","topic":"s5_log"}]*，那就是分别创建两个收集文件的实例，监控 _D:\\studyGo\\s4.log_ ,和 _D:\\studyGo\\s5.log_ 文件，并且发送到kafka的 _web_log_ 和 _s5_log_ 的 __topic__ 中  
* **etcdemo**：etcdemo为配置etcd中参数的服务。主要是为了服务于ByranLog所需的配置灵活测试。也就是单纯的与etcd打交道，增加配置和读取配置   
* **log_transfer**：log_transfer为kafka消费服务。通过订阅kafka中的 __topic__ 将写入kafka的数据消费出来，并发送到elasticsearch进行持久化存储，然后在通过es的经典搭档kiban。配置好索引模式，将收集到的数据展示出来。  

## 效果展示
* 开启本地的zookeeper，kafka，etcd，elasticsearch服务
![图片名称](https://github.com/chenjiahui-xxz/IMG/blob/main/zookeeper.png)  
![图片名称](https://github.com/chenjiahui-xxz/IMG/blob/main/kafka.png) 
![图片名称](https://github.com/chenjiahui-xxz/IMG/blob/main/etcd.png) 
![图片名称](https://github.com/chenjiahui-xxz/IMG/blob/main/elasticsearch.png) 

* 启动etcd服务，配置好etcd里面的参数
![图片名称](https://github.com/chenjiahui-xxz/IMG/blob/main/etcddemo.png)

* 启动ByranLog服务，根据etcd中的参数配置开启监控文件实例，并且将数据收集发送到kafka
![图片名称](https://github.com/chenjiahui-xxz/IMG/blob/main/runBryanLog.png)

* 启动log_transfer服务，消费kafka中的数据，并且发送到elasticsearch
![图片名称](https://github.com/chenjiahui-xxz/IMG/blob/main/runConsume.png)

* 启动kibnan服务，建立索引模式，将数据显示到页面来
![图片名称](https://github.com/chenjiahui-xxz/IMG/blob/main/kibnan.png)
![图片名称](https://github.com/chenjiahui-xxz/IMG/blob/main/kibnanShow.png)

## 踩过的坑记录
```
0.安装kafak还有zoopkeeper下来后，最好修改下运行时kafka产生的日志位置，修改下面文件。
kafka：.\config\server.properties
zookeeper：.\config\zookeeper.properties

1.开启zookeeper，kafak，消费kafka命令。启动kafka需要事先启动zookeeper。可以分别安装。也可以使用kafka自带的zookeeper（我就是使用自带的）
.\bin\windows\zookeeper-server-start.bat .\config\zookeeper.properties

.\bin\windows\kafka-server-start.bat .\config\server.properties

.\bin\windows\kafka-console-consumer.bat --bootstrap-server 127.0.0.1:9092 --topic web_log --from-beginning
eg：--bootstrap-server：连接的是哪一个kafka ；--topic：主题；--from-beginning：从最开始读取

2.kafka的包 "github.com/Shopify/sarama" window下面要使用v1.19之前的版本，
v1.20之后的版本加入了zstd压缩算法，需要用到cgo。在window平台下面编译时会提示错误gcc找不到（最新版本的saram包支持window版本，可以尝试）

3.监听文件变化的包 "github.com/hpcloud/tail",下载的时候有些依赖包已经变了，需要修改包依赖，可百度
https://www.cnblogs.com/kenLoong/p/15452313.html

4.etcd默认是v2版本没有put命令，要设置下环境变量SET ETCDCTL_API=3,然后就可以设置
etcdctl.exe --endpoints=http://127.0.0.1:2379 put bryan "dd"
etcdctl.exe --endpoints=http://127.0.0.1:2379 get bryan

5.安装etcd客户端、事先要安装好etcd的服务端，然后去包里面打开etcd.exe即可
go get go.etcd.io/etcd/client/v3

6.安装infludb，因为v1跟v2版本的go语言命令是不一样的，我这边安装v1.7.7，直接下载
https://dl.influxdata.com/influxdb/releases/influxdb-1.7.7_windows_amd64.zip
相关操作包，命令介绍可以看这个https://www.fdevops.com/docs/golang-tutorial/influxdb-operating
可以改下配置文件，因为不是linux，最好改下：如地址
# Where the metadata/raft database is stored
  dir = "D:\studyGo\influxdb\meta"
相关地址我都修改了，如果是mac可以直接不管

7.influxDB 1.x版本
go get github.com/influxdata/influxdb1-client/v2
influxDB 2.x版本
go get github.com/influxdata/influxdb-client-go

8.influxdb数据库要手动创建
操作文档：https://docs.influxdata.com/influxdb/v1.7/introduction/getting-started/

9.grafana下载，下载的是6.2.5版本；也可以下载最新
https://grafana.com/grafana/download/8.4.4?pg=get&plcmt=selfmanaged-box1-cta1&platform=windows

10.根据官网描述，下载下来后需要修改conf里面，将defaults.ini复制一份修改为custom.ini
11.发现6.2.5打开的时候一直会报错 Service init failed: License token file not found: D:\\place\\grafana-enterprise-6.2.5.windows-amd64\\grafana-6.2.5\\data\\license.jwt；最后我直接换成最新版

12.grafana 要先新增一个数据源，再去新增一个仪表盘，不然仪表盘里面是读不到这个数据源，也就是读不到那你配置的数据库
13.Database:  "monitor",这个是数据库，需要提前创建，
name是表名；比如memory，会自己通过代码创建

13.下载Elasticsearch，https://www.elastic.co/cn/downloads/elasticsearch
14.启动es的服务端，cd D:\place\elasticsearch-8.1.2-windows-x86_64\elasticsearch-8.1.2
然后 .\bin\elasticsearch.bat
15.es用localhost：9200报错；received plaintext http traffic on an https channel, closing connection Netty4HttpChannel；
其实是因为开启了ssl；将配置文件xpack.security.http.ssl:改为false
16.es进去后需要账号跟密码，直接改为免密码登录：xpack.security.enabled：false

17.上面是用最新版本的elasticsearch-8.1.2；如果你使用post请求，按照以往的
curl -H "ContentType:application/json" -X POST 127.0.0.1:9200/user/person -d '
{
	"name": "dsb",
	"age": 9000,
	"married": true
}'绝对出错，因为8.0版本已经移除了这个type类型的语法，所以我后面装回低版本
https://mirrors.huaweicloud.com/elasticsearch/7.3.0/

18.用低版本就是舒服：cd D:\place\elasticsearch-7.3.0-windows-x86_64\elasticsearch-7.3.0\bin
然后.\elasticsearch.bat

19.es的版本要跟kibana对应上 现在都是7.3.0

20.kibana改下配置文件，elasticsearch.hosts: ["http://localhost:9200"]；i18n.locale: "zh-CN"（支持中文）
启动：D:\place\kibana-7.3.0\kibana-7.3.0-windows-x86_64\bin，然后.\kibana.bat；就可以在仪表旁那边看到数据了

21.kibana新建了一个index，在索引模式下面得kibana里面创建索引后，就能看到数据啦
```

### 最终
**使用到了zookeeper，kafka，etcd，elasticsearch，kibnan。可以拉下来玩玩。还是挺好玩的收集系统。My name is Bryan**
