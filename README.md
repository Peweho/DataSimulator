# 数据模拟器

## 使用

### 1、填写配置文件

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