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


