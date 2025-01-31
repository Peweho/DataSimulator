# 数据模拟器

## 使用

使用分为真实环境和虚拟环境，提供配置文件的Env参数确定

- 真实环境：时间戳为真实时间，数据发送到Kafka消息队列，没有配置输出到控制台，支持为每个数据定义生成频率
- 虚拟环境：根据time的配置生成时间戳，最后输出Excel文件，不支持每个数据定义生成频率

### 1、真实环境配置

```yaml
# 数据部分
data:
  - title: "温度"
    id: "temperature"
    frequency: 10 #生成频率。单位s
    min:
    max: 30
    model: gaussian   # 生成模型类型，默认随机生成
    params: [25,2]   # 模型参数

  - title: "湿度"
    id: "humidity"
    frequency: 10
    min: 50
    max: 60
    model: exp
    params: [-0.5]
mq:
  addr:
  topic:
  partition: 0
  timeout: 500  #ms

env: truth
```

### 2、虚拟环境配置

```yaml
time:
  start: "2021-01-01 00:00:00"  #起始时间
  end: "2021-01-01 01:00:00"    #终止时间
  interval:  100 #时间间隔，单位秒

# 数据部分
data:
  - title: "tempture"
    min: 0
    max: 30
    model: gaussian   # 生成模型类型，默认随机生成
    params: [25,2]   # 模型参数
  - title: "humidity"
    min: 0
    max: 100
    model: random
    params:
    
env: virtual
```

time接收三项配置

1. start：表示模拟数据的起始时间
2. end：表示模拟数据的终止时间
3. interval：表示每隔多久生成一次数据

data接收一个数组，每一项表示一个数据项

1. title：数据项标题
2. min：数据的最小值
3. max：数据的最大值
4. model：生成数据的模型；模型提供正态分布，指数分布和随机生成
5. params：模型接收参数

### 2、启动程序

```powershell
simulator path
```

`path`为配置文件路径，默认为`./etc.yaml`

在当前目录下生成excel文件

## 生成数据模型

模型提供正态分布（gaussian），指数分布（exp）和均匀分布（random）

### 1、Gaussian

```yaml
model: gaussian
params: [均值, 标准差] #接收两个参数
```

当模型选为`gaussian`时不接受最大值和最小值

### 2、Exp

```yaml
model: exp
params: [lambda] 		#接收一个参数
```

当模型选为`exp`时

- 如果lambda为正数时，生成数据大于等于min
- 如果lambda为负数时，生成数据小于等于max

### 3、均匀分布

```yaml
model: random
params:             #不接受参数
```

数据范围为[min,max)

## Kafka

提供将生成的数据发送到Kafka上，需要填写对应信息以及选择真实环境，无需填写

```yaml
mq:
  addr:
  topic:
  partition: 0
  timeout: 500  #ms
env: truth
```

## HTTP服务

提供三个接口，分别是修改数据的生成频率、修改绑定的数据库、查询数据详情

### GET 查询数据详情

GET /list

> 返回示例

> 200 Response

```json
{
  "id": "string",
  "frequency": "string"
}
```

#### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline   |

#### 返回数据结构

状态码 **200**

| 名称        | 类型   | 必选 | 约束 | 中文名   | 说明 |
| ----------- | ------ | ---- | ---- | -------- | ---- |
| » id        | string | true | none | 数据项id | none |
| » frequency | string | true | none | 采集频率 | none |

### POST 修改数据采集频率

POST /set/frequency

> Body 请求参数

```yaml
id: temperature
frequency: 12

```

#### 请求参数

| 名称        | 位置 | 类型    | 必选 | 说明                   |
| ----------- | ---- | ------- | ---- | ---------------------- |
| body        | body | object  | 否   | none                   |
| » id        | body | string  | 否   | 数据项id               |
| » frequency | body | integer | 否   | 改后的采集频率，单位秒 |

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

#### 返回结果

| 状态码 | 状态码含义                                                   | 说明 | 数据模型 |
| ------ | ------------------------------------------------------------ | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)      | none | Inline   |
| 400    | [Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1) | none | Inline   |

#### 返回数据结构

状态码 **200**

| 名称  | 类型   | 必选 | 约束 | 中文名   | 说明 |
| ----- | ------ | ---- | ---- | -------- | ---- |
| » msg | string | true | none | 返回信息 | none |

状态码 **400**

| 名称  | 类型   | 必选 | 约束 | 中文名   | 说明 |
| ----- | ------ | ---- | ---- | -------- | ---- |
| » msg | string | true | none | 返回信息 | none |

### POST 绑定数据库

POST /bind/db

> Body 请求参数

```yaml
id: ""

```

#### 请求参数

| 名称 | 位置 | 类型   | 必选 | 说明     |
| ---- | ---- | ------ | ---- | -------- |
| body | body | object | 否   | none     |
| » id | body | string | 否   | 数据库id |

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

#### 返回结果

| 状态码 | 状态码含义                                                   | 说明 | 数据模型 |
| ------ | ------------------------------------------------------------ | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)      | none | Inline   |
| 400    | [Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1) | none | Inline   |

#### 返回数据结构

状态码 **200**

| 名称  | 类型   | 必选 | 约束 | 中文名   | 说明 |
| ----- | ------ | ---- | ---- | -------- | ---- |
| » msg | string | true | none | 返回信息 | none |

状态码 **400**

| 名称  | 类型   | 必选 | 约束 | 中文名   | 说明 |
| ----- | ------ | ---- | ---- | -------- | ---- |
| » msg | string | true | none | 返回信息 | none |
