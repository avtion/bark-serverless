# Bark-Serverless

bark-server 腾讯云SCF版本 - bark-server tencent cloud SCF version
> [云函数SCF](https://cloud.tencent.com/document/product/583) 是腾讯云为企业和开发者们提供的无服务器执行环境，帮助您在无需购买和管理服务器的情况下运行代码。

# ✈️ 介绍

- 本项目为 **bark-Server** 的 **腾讯云SCF云函数** 个人重构版本
- 需要1个 **域名** 用于 Bark 客户端添加私有服务器
- 支持 **bark 1.2.0**
- 依赖 **bark-server v2.0.2**
- 无需数据库和存储空间
- `v1.1.0` 基于 GORM 支持 PostgreSQL、MySQL、SQL Server 和 Clickhouse

# 📚 安装说明

1. 从 [Github Realase](https://github.com/avtion/bark-serverless/releases/) 下载编译好的可执行文件
2. 登录腾讯云
3. 创建一个Serverless云函数，运行环境选择Go1，上传可执行文件的zip
4. 进入创建好的云函数 **触发管理** 页面，创建一个触发器，选择 **API网关触发** ，请求方法选择 `ANY`
5. 进入腾讯云 **API网关** 管理页面，选择刚刚创建的服务，编辑刚刚创建的API，路径改成"/"，保存
6. 修改完服务路径之后再获取API默认的访问地址，如 https://service-00wc1lm6-12********.gz.apigw.tencentcs.com/release/ ，这样是不能直接作为Bark
   APP的私有服务器的地址，所以接下来有2种方法解决

## 6.1 利用域名的隐性URL解析（推荐）

域名DNS解析中添加记录类型为 **隐性URL** 的解析，其中记录值为上面第6步获取的API访问地址
> 本人使用的是 **阿里云** 提供的域名解析服务，需要 **ICP备案** ，不保证每个人都能实现

## 6.2 绑定API网关自定义域名

1. 在API网关中，Bark-Serverless服务管理页面选择 **自定义域名** 并新建
2. 尽量申请免费的证书开启HTTPS访问，注意网络安全
3. （关键）在路径映射选项中选择自定义路径映射，并设置 **发布** 环境的路径为"/"
4. 提交之后就能直接使用自定义域名访问API服务

# ☘️ 使用说明

1. Bark APP添加私有服务器之后，程序会输出以下内容，请前往云函数的`日志查询`中查找

```JSON
{
  "level": "info",
  "ts": 1623430963.414276,
  "caller": "controller/register.go:66",
  "msg": "设备绑定信息",
  "router": "register",
  "key": "9GMMk5JhTEL*****",
  "token": "7008fb1e25ff2f91aa80db4ff56141456e**********",
  "old_key": "9GMMk5JhTEL*****",
  "old_token": "7008fb1e25ff2f91aa80db4ff56141456e**********"
}
```

2. 在上述信息中，需要key和token，请注意不要泄漏这两项数据
3. 进入云函数的 **函数配置** ，点击编辑，环境变量中新增一项，**键是 device_前缀加key的值**（eg. device_9GMMk5JhTEL*****
   ），值是 **token**（eg. 7008fb1e25ff2f91aa80db4ff56141456e**********）
4. 点击 **保存** 按钮

# 🥺 发送消息提示 failed to get token from env

请重复 **使用说明** 的过程添加key和token到云函数的环境变量

# ✨ 为什么用SCF

1. 要推送一些比较隐私（如验证码）的内容
3. 腾讯云的SCF免费提供40万GBs资源使用量和100万次事件型函数调用次数
4. 腾讯云API网关第一年每月（自然月）前100万次调用免费
5. 快乐白嫖

# 👍 请支持Bark项目

- [Bark](https://github.com/Finb/Bark) - IOS 客户端

- [bark-server](https://github.com/Finb/bark-server) - Golang 服务端

# 📢 开源声明

MIT License