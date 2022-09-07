# cmp
├── main.go        # main 入口
├── cmd            # 可独立执行程序
├── configs        # 存放配置文件
├── internal
│   ├── api/v1     # api 接口逻辑，控制路由进来后执行那个函数  /k8s/cluster
│   ├── conf       # 项目配置
│   ├── model      # model定义  是实体模型定义，一般是一个实体一个文件，比如 k8s_cluster.go  定义集群相关的结构体
│   ├── router     # 路由入口
│   ├── service    # 是实体和数据库处理逻辑存放，比如新增一个实体；操作数据库对数据进行增删改查
│   ├── middleware # 中间件
├── pkg            # 可供外部使用的 package
└── web            # 前端
list 名称 ip  状态 角色 cpu 内存
Info
基本信息
名称 存活时间 标签 注释

系统信息
系统架构 操作系统 操作系统内核  容器云runtime kubelt版本

已分配资源
cpureq lilmt  内存req limit pods

健康状态  
网络可用性 内存压力 磁盘压力 进程压力 节点就绪状态

pod

镜像 