# RestCheck
一个用来实现 Http API 接口的工具

## 优势
- 命令行实现, 文件形式保存所有配置
  - 方便 ci 集成
  - 可以用 git 管理整个配置目录, 团队协作
- 结果复用
    - 如果接口B 依赖于接口 A 的结果里面的 user.token 字段
    - 那么可以在 B 的请求参数里面用 {restcheck:listLogisticsTracking -> user.token}


## 使用
命令行使用

`-init {project name}` 
- 启动一个项目
- 新建一个文件夹, 里面包含一个默认的 restcheck.conf
    - restchec.conf
      - 规定项目名称, 项目情况
      - env 配置, 多 env (默认env)
      - 不同 env 可以有不同的 域名 / header 配置

`-init {api-name}`

`-check {api-name}`

检查某个 api
- 如果没有 api, 那么就会提示创建一个 api, 并且在 http 里面去实现

`-save`
请求的结果作为文件保存下来, 作为后续的校验结果

`-checkall`
一次新检查所有 api

`-env {env}`
设置本次的环境

## 使用

```shell

# 初始化 hades 项目
restcheck -init ezecopy

# 查看当前目录的详情, env 情况, tree 打印所有 api 情况, 类似于 tree 命令
restcheck -msg

# 初始化一个接口配置
restcheck -initapi listNote

# 初始化一个组里面的一个接口
restcheck -initapi note/listNote

# 检查 listNote 接口并且将结果保存下来
restcheck -check note/listNote -save

# 检查 listNote 接口
restcheck -check note/listNote

# 检查 note 组下面的所有接口
restcheck -check note/

# 检查所有接口
restcheck -checkall

```

## 工作目录
```shell
└── restcheck.conf     ---> 项目配置, env / 域名 / 通用 header 之类的
│
├── listNote                                  ---> 一个 API 目录, api 名称: listTracking
│  ├── request.http                           ---> API 的配置
│  │
│  ├── prd                                    ---> prd 环境的测试情况
│  │  ├── 20231129_141731_results.txt
│  │  └── 20231129_141851_test.txt
│  └── test                                   ---> test 环境的测试情况
│     ├── 20231129_141731_results.txt         ---> 用来做校对的 result
│     └── 20231129_141851_test.txt            ---> 最新测试得到的 result
| 
|── listShare                                 ---> 一个 API 目录, api 名称: listUsers
│  ├── request.http                           ---> ...
|  ........

```


## 风险
- 某些 post 跑到了正式环境上可能会造成 prd 数据污染？
  - 来自 restchecck 的 header 里面默认带一个 rest-check-id
  - prd 如果检测到这个 header 就直接拒绝掉?


