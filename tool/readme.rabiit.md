## 安装rabbitmq
在`47.96.172.8`服务器上一安装好一个, rabbit的账号密码是`client`/`wOtYTFRN`, 对应的url为 `amqp://client:wOtYTFRN@116.62.167.76:5672/`

### 安装
百度 `apt-get rabbitmq`
### 配置
rabbitmq默认不支持远程登录, 需要新建账号与权限
百度 `rabbitmq 远程登录`

## 如何使用rabbitmq
rabbitmq是消息队列, 由生产者和消费者组成, 可有多个消费者和生产者并且可以远程连接, 所以rabbit也可以用用在双向通讯

### 发送
rabbit 中有 `queue`概念, 生产者发送在同一个`queue`中的消息才能别消费者收到, 所以我们要给某个箱子发消息只需要这样定义`queue`:

`box_1`, 其中`1`是指箱子id, 发送消息代码如下
```
err = tool.Rabbit.Publish("box_1", body)
```
其中`body` 是[]byte, 自行定义

当我们生产者对`box_1`发送了消息, 硬件方需要消费才行, 让硬件方`订阅`这个`box_1` queue就行, 然后把`body`的协议内容告诉他, 他们就能收到并解析

### 接收
可接收来至任何queue的消息

代码
```
queue := "info"
err := tool.Rabbit.Receive(queue, handleInfo)
```
```
type BoxInfoFromHardware struct {
	B float32 `json:"a"`
	A float32 `json:"b"`
}
// 处理盒子上报自身信息
func handleInfo(msg amqp.Delivery) (error) {
	s := BoxInfoFromHardware{}
	err := json.Unmarshal(msg.Body, &s)
	if err != nil {
		return err
	}
	...

	return nil
}
```
现在叫硬件方往`info`queue里发送我们定义好的结构就可以了.

