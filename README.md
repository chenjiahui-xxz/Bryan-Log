# Bryan-study
Bryan的日志收集系统，很好玩的一个文件输入后收集，就像你在写日志将你写的日志收集起来，可以通过etcd配置，不需要重启服务就能够实时监听到需要收集的那些文件，通过etcd的watch监听。将收集到的信息通过管道发送到kafka里面。然后在起一个消费者服务，将kafka的数据进行消费，存储到elasticsearch中

## 包名介绍
* **ByranLog**：ByranLog为日志收集，会先初始化kafka实例，然后在连接etcd服务。拉去etcd中的配置来创建出需要的读取文件实例，同时还能实时的通过 __etcd.watch__ 监控etcd配置文件的变化，来决定是否新增新的配置读取文件实例，或者是删除配置文件中已经没有的实例。实现不需要ByranLog重启服务而能动态增删改灵活读取所需要的文件日志。将读取的到的文件内容通过go管道发送到kafka中进行持久化。如etcd中的配置 *[{"path":"D:\\studyGo\\s4.log","topic":"web_log"},{"path":"D:\\studyGo\\s5.log","topic":"s5_log"}]*，那就是分别创建两个收集文件的实例，监控 _D:\\studyGo\\s4.log_ ,和 _D:\\studyGo\\s5.log_ 文件，并且发送到kafka的 _web_log_ 和 _s5_log_ 的 __topic__ 中  
* **etcdemo**：etcdemo为配置etcd中参数的服务。主要是为了服务于ByranLog所需的配置灵活测试。也就是单纯的与etcd打交道，增加配置和读取配置   
* **log_transfer**：log_transfer为kafka消费服务。通过订阅kafka中的 __topic__ 将写入kafka的数据消费出来，并发送到elasticsearch进行持久化存储，然后在通过es的经典搭档kiban。配置好索引模式，将收集到的数据展示出来。  
