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

node 名称 标签 准备就绪 CPU 下限 (cores) CPU 上限 (cores) 内存下限 (bytes) 内存上限 (bytes) Pods 创建时间
minikube
beta.kubernetes.io/arch: amd64
beta.kubernetes.io/os: linux
kubernetes.io/arch: amd64
True
1.15 (28.75%)
400.00m (10.00%)
470.00Mi (5.98%)
270.00Mi (3.43%)
17 (15.45%)
8 days ago