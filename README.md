# 简介

Nightingale接入prometheus query接口做数据源，替代Alertmanager作为CNCF监控体系中告警模块的补充。


# 使用方式

正常安装夜莺后，安装插件。

直接从prometheus的API接口拿取prometheus的已存储数据，方便在已有prometheus+gragana监控体系的情况下，进一步提供告警机制。



## 快速构建 

    $ mkdir -p $GOPATH/src/github.com/n9e
    $ cd $GOPATH/src/github.com/n9e
    $ git clone https://github.com/JamesJoe-C/n9e-query-prometheus.git
    $ cd n9e-query-prometheus
    $ export GO111MODULE=on
    $ export GOPROXY=https://goproxy.cn
    $ go build
    $ cat plugin.test.json | ./n9e-query-prometheus

 
## 插件使用
夜莺中stdin如下：
```
{
  "exporter_urls": [
    "http://prometheus.url"
  ],
  "query": "service_memory_load",
  "timeout": 5000
}
```

目前仅测试到Gauge数据格式，后续更新其他数据格式DEMO。

