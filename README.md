# 基于LRU实现的本地缓存管理
  
## 功能点
1. 基于HTTP提供缓存读写服务
2. 缓存管理采用LRU算法
3. 自带缓存穿透保护机制

## 默认配置项
1. 默认配置LRU容量为1GB
2. URL默认地址：http://127.0.0.1:9999/t_cache/

## 使用方法
启动服务后，通过HTTP请求进行读写操作：
1. 读：
   - GET http://127.0.0.1:9999/_tcache/缓存名/key
   - JSON格式返回，key为实际的key，Value则为缓存的值
2. 写：
   - POST http://127.0.0.1:9999/_tcache/缓存名/key
   - Value在body中

可根据状态码判断读写是否请求成功，状态码200表示成功，其他表示失败

## 例
1. 启动服务：
```shell
go run main.go
```
2. 设置一个分数缓存，名称为score，key为学生名
```shell
curl --location '127.0.0.1:9999/_tcache/score/zhangsan' \
--header 'Content-Type: text/plain' \
--data '561'
```
3. 读取缓存内容
```shell
curl --location '127.0.0.1:9999/_tcache/score/zhangsan'
```
返回结果示例：
```json
{
    "key": "zhangsan",
    "value": "561"
}
```
