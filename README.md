# Bark-Serverless
bark-server tencent cloud SCF(Serverless Cloud Function) version

bark-server的腾讯云SCF版本

# ✈️ 介绍

- 本项目为`bark-Server`的`腾讯云SCF`个人重构版本

- 需要`1个域名`用于Bark APP添加私有服务器

- 无需数据库和存储空间

- 支持bark 1.1.5

- 依赖bark-server@v2

# 📚 安装说明

1. 从Github Realase下载编译好的可执行文件

2. 登录腾讯云

3. 创建一个Serverless云函数，运行环境选择Go1，上传可执行文件的zip

4. 进入创建好的云函数`触发管理`页面，创建一个触发器，选择`API网关触发`，请求方法选择ANY

5. 进入腾讯云`API网关`管理页面，选择刚刚创建的服务，编辑刚刚创建的API，路径改成"/"，保存

6. 修改完服务路径之后再获取API默认的访问地址，如"[https://service-00wc1lm6-12********.gz.apigw.tencentcs.com/](https://service-00wc1lm6-1258029428.gz.apigw.tencentcs.com)release/"，这样是不能直接作为Bark APP的私有服务器的地址，所以接下来有2种方法解决

## 1. 利用域名的隐性URL解析（推荐）

域名DNS解析中添加记录类型为`隐性URL`的解析，其中记录值为上面第6步获取的API访问地址

> 本人使用的是`阿里云`提供的域名解析服务，需要`ICP备案`，其他域名服务提供商没试过，这个属于骚操作，不保证每个人都能实现

## 2. 绑定API网关自定义域名

1. 在API网关中，Bark-Serverless服务管理页面选择`自定义域名`并新建

2. 尽量申请免费的证书开启HTTPS访问，注意网络安全

3. （关键）在路径映射选项中选择自定义路径映射，并设置`发布`环境的路径为"/"

4. 提交之后就能直接使用自定义域名访问API服务

# ☘️ 使用说明

1. Bark APP添加私有服务器之后，程序会输出以下内容，请前往云函数的`日志查询`中查找

```JSON
{"level":"info","ts":1623430963.414276,"caller":"controller/register.go:66","msg":"设备绑定信息","router":"register","key":"9GMMk5JhTEL*****","token":"7008fb1e25ff2f91aa80db4ff56141456e**********","old_key":"9GMMk5JhTEL*****","old_token":"7008fb1e25ff2f91aa80db4ff56141456e**********"}
```

1. 在上述信息中，需要key和token，请注意不要泄漏这两项数据

2. 进入云函数的`函数配置`，点击编辑，环境变量中新增一项，键是"device_"前缀加key的值（示例：device_9GMMk5JhTEL***** ），值是token（示例：7008fb1e25ff2f91aa80db4ff56141456e**********）

3. 保存

# 🥺 发送消息提示failed to get token from env

请重复`使用说明`的过程添加key和token到云函数的环境变量

# ✨ 为什么用SCF

1. 要推送一些比较隐私（如验证码）的内容

2. 没有AWS账号

3. 腾讯云的SCF免费提供40万GBs资源使用量和100万次事件型函数调用次数

4. 腾讯云API网关第一年每月（自然月）前100万次调用免费

5. 没钱

# 👍 请支持Bark项目

- [Bark]([https://github.com/Finb/Bark](https://github.com/Finb/Bark))

- [bark-server]([https://github.com/Finb/bark-server](https://github.com/Finb/bark-server))

# 📢 开源声明

MIT License