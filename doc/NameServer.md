[TOC]

## 一. 较传统DNS的优势

使用传统的DNS服务，主要面临电信运营商 Local DNS 的问题：

| 传统DNS                                                      | HTTPDNS                                                      |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| 域名劫持，缓存污染                                           | 防支持。绕过运营商Local DNS，HTTPs通信                       |
| 调度不精准。解析服务转发、Local DNS部署不均、解析规则配置不当 | 直接IP访问，直接获取客户端IP ，基于客户端IP获得就近接入业务节点 |
| 解析生效滞后。Local DNS缓存不刷新、TTL不生效                 | 热点预解析，缓存解析结果                                     |
| 延迟大。每次查询需要递归遍历多个DNS服务器以获取最终结果      | 遵循TTL，全网快速同步                                        |
| 无                                                           | 多级别(/线路/区域/业务字段)精细到单个IP和业务字段的灵活调度  |



## 二. 概要

- 支持业务接入点分组，组内流量均衡负载。
- 业务接入点应该分主、备两个组。客户端永远都应该按主、备顺序尝试接入。
- 支持客户端流量灰度，支持按业务字段解析。
- 客户端需要有重试&降级逻辑。
- 可以用作LocalDNS，可以用作权威DNS，可以用作GLB。




## 三. 接口

支持HTTP2、HTTTPs、HTTP、DNS等接入协议。用户业务使用HTTPDNS时，应做好异常情况下的出错兼容逻辑，主要包括**异步请求**、**重试**、**降级**、**网络切换**、**状态切换**等。

DNS协议遵循标准。HTTPx协议详细如下：

### 1. 域名解析请求

HTTPDNS通过HTTP接口对外提供域名解析服务，用户接入直接使用IP地址访问HTTPDNS服务。一次请求可解析多个域名。

```
GET http://1.2.3.4/d?dn=x.com&dn=y.com&ip=1.1.1.1&k=v
```

URL参数说明：

| 名称 | 是否必须 | 描述|
| ---- | -------- | ----------------------------|
| dn   | 必须     | 要解析的域名，可以有多个，一次最多10个。 |
| ip   | 可选     | 用户的来源IP，如果没指定这个参数，默认使用请求连接的源IP |
| k  | 可选     | 业务字段；可用于区分应用、版本、平台、用户…等业务层逻辑；参数名标识调度策略，参数值参于策略计算。 |

请求成功时，HTTP响应的状态码为200，响应结果用JSON格式表示，示例如下：

```
{
	"s":0,
	"e":"错误信息",
	"u":"1.1.1.1",
	"v":"中国广东广州_电信",
	"data":[
		{
      "n":"x.com",
			"s":0,
			"rs": [
    			{
    			"type":"A",
    			"host": "10.1.80.172",
          "ttl": 60,
        	"http": 80,
        	"http2":443
    			},
    			... ...    
    		]
		}
	]
}
```

字段说明：
| 名称 | 描述                                                         |
| ---- | ------------------------------------------------------------ |
| s    | 应答结果状态码。0：成功处理；1：服务端出错；2：请求参数错误；3：找不到相关域名记录； |
| e    | 应答结果状态码为非0时，显示错误信息。                        |
| u    | 客户端IP                                                     |
| v    | 客户端线路和地域信息                                         |
| data | 解析结果数组                                                 |
| n    | 请求解析的域名。                                             |
| ttl  | 该域名解析结果的TTL缓存时间。 用户应该按这个TTL时间对域名解析结果进行缓存。在TTL过期之前，直接使用缓存的IP；TTL值为0不缓存。 |
| rs   | 该域名的解析结果列表，包括解析结果记录以及记录类型(type)：A、CNAME等。和这个地址上(A记录)提供服务的协议及端口号；协议包括：http、http2中的一个或多个。 |

### 2. NameServer服务器IP管理

因为NameServer服务器IP是会变更的，而且接入质量也有差异。服务端会跟据终端不同返回最优的NameServer服务器IP列表，供客户端轮询重试。

假定客户端所在的国家，全国部署有3个机房。NameServer自身服务器IP需要解决两个问题：

- **客户端如何同步更新NameServer服务器IP**
  1. 在客户端发版时，把这些IP及属性(运营商及所在省份)写在配置文件中一同发布。
  2. 客户端在运行过程中， 每隔1周(在wifi状态下)通过HTTPDNS接口获得最新服务器IP列表，并保存到配置文件中。只要所有IP不时时变更，就可获得最新的服务器IP列表。
  3. 如果所有服务器IP都不通，DNS解析请求降级为LocalDNS。并会通过监控系统告警，人工介入处理。
- **客户端进行域名解析请求时，如何优化选择NameServer服务器IP**
  1. 客户端第一次DNS解析请求，采用并发方式同时向所有IP发起HTTPDNS请求，以及LocalDNS发超UDP请求。使用优先返回的HTTPDNS结果，如果所有HTTPDNS都超时则使用LocalDNS结果。
  2. 跟据服务端返回的结果，缓存本机的线路及区域信息(v字段)，及对应的IP列表。
  3. 跟据HTTPDNS接口返回的本机的线路及区域信息选择对应的Ip列表发起HTTPDNS查询。

NameServer服务器IP更新接口与HTTPDNS使用同一个接口，当显示请求解析`httpdns.com`域名时，表示获取最佳服务器IP列表:

```json
GET http://1.2.3.4/d?dn=httpdns.com&dn=y.com&ip=1.1.1.1&k=v
```



## 四. 解析规则

当一个客户端请求解析一个域名，解析会分3个优先级别：

- 如果带有业务参数(非dn、ip)，用业务定制规则解析。暂时需要定制开发。
- 如果有配置线路&区域解析规则，命中规则即返回结果。
- 如果前两步都没有命中，并配置有DNS记录，查询所配置的DNS服务并返回结果。

### 1. 业务定制解析

业务定制解析优先级最高。只要请求参数里带有除dn、ip之外的第3个参数，就表明开启业务定制解析。暂时只支持一个参数维度。

例如按用户id调度。客户端发如下请求：

```
GET http://1.2.3.4/d?uid=123456
```

并且服务端有关键字uid相关的规则配置。可能是：

```
数据类型：int64
规则1：< 10000 group1
规则2：>= 10000 group2
```

这时候，会用123456去依次匹配解析规则。

> 暂不支持配合线路/区域多维度解析， 后续如果业务有需求可以支持。

### 2. 线路&区域就近解析

跟据客户端ip，查询IP段得到客户端所在的线路和区域，跟据线路和区域就近调度。

**线路** 国内用户区分到电信、联通、移动…等运营商。国外区分到国家或按业务实际情况定义。

**区域** 国内用户区分粒度到省。可配置单个客户端IP为独立的解析区域。

- 查询客户端IP所在的段， 得到线路&区域键值如：` /电信/中国广东`。
- 业务服务接入点，分成多个不同的组，如：`南方电信A组`，包含多条接入服务IP或CNAME域名。
- 不同的客户端线路&区域配置对应到不同的接入点组。
- 同组内的接入点记录，按时间轮询返回给客户端。`current_timestamp % count(group)`

### 3. 权威DNS查询

这种情况本服务相当于Local DNS角色，所有调度策略都配置在权威DNS。通过EDNS查询权威DNS并返回结果。

### 4. 默认解析

 每个域名的每个线路都有一个默认解析值，当没有更精细的解析规则命中，返回默认值。配置时以`*`表示。



## 五. 负载均衡

业务接入点之间加权轮询，第一版不支持。



## 六. 数据结构

**解析规则配置保存在mysql里**：

```
CREATE SCHEMA `ns` DEFAULT CHARACTER SET utf8 COLLATE utf8_bin ;
```

**资源记录(Resource Record, RR)**：
```
CREATE TABLE `rr` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain` varchar(128) NOT NULL,
  `ttl` int(11) NOT NULL DEFAULT '600',
  `type` tinyint(4) NOT NULL,
  `class` tinyint(4) NOT NULL DEFAULT '1',
  `data` varchar(45) NOT NULL,
  `group` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_domain_type` (`domain`,`type`),
  KEY `index_group` (`group`,`domain`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

```

**客户分区**:

```
CREATE TABLE `ns`.`zone` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `line` varchar(45) COLLATE utf8_bin NOT NULL,
  `area` varchar(45) COLLATE utf8_bin NOT NULL,
  `zone` varchar(45) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_line_area` (`line`,`area`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

```

**解析映射**：
```
CREATE TABLE `ns`.`rule` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain` varchar(256) COLLATE utf8_bin NOT NULL,
  `zone` varchar(45) COLLATE utf8_bin NOT NULL,
  `group` varchar(45) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
```



## 七. 配置管理

1. 添加域名记录接口
2. 删除域名记录接口
3. 设置IP zone接口
4. 设置解析策略接口



## 八. 部署

![ns](./ns.png)

> 1. 为保证服务可用，必须多线路多机房部署
> 2. 解析规则保存在MySQL，管理后台读写MySQL主节点。每部署一个新的NameServer，添加一个MySQL备节点。
> 3. 权威DNS需要支持EDNS，保证就近解析。



## 九. 质量评估

后续支持



## 十. 监控

