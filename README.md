# larkGPT-notice
It can interface with various third-party services, or use GRPC to complete Feishu group notifications, and can call ChatGPT to complete Q&amp;A

## 配置
### 1.配置app.test.yaml
    - openai的配置是为了完成飞书自动对话，使用chatGPT完成对话。因为墙的原因，在项目根目录下提供了nginx.conf模版，可包部署到海外服务器代理openai的请求。
    - 另外飞书的相关配置请参考开放平台https://open.feishu.cn/app?lang=zh-CN，应用机器人的配置请看文档，这里不多做复述。另外注意，飞书开放api的权限需要申请，并添加到自己的事件中，才能触发自动对话。
    - 项目还提供了grpc的相关功能，可以通过grpc调用飞书的机器人，完成消息通知。grpc的配置请参考项目根目录下的grpc.proto文件，这里不多做复述。

## 运行
### make build
    - 编译项目会生成bin/lark可执行文件
### bin/lark http
    - 启动http服务
### bin/lark grpc
    - 启动grpc服务
    - 注意，修改了proto后需要执行 `make protoc` 重新生成grpc相关代码

## 演示
### 1.演示图片
消息通知
![image](https://github.com/zealerFT/feishuGPT/blob/main/resource/one.png)
机器人对话，点赞
![image](https://github.com/zealerFT/feishuGPT/blob/main/resource/two.png)
