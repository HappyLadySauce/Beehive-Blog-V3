# Beehive Blog v3 Edge 与 Gateway 路由设计

## 1. 目标

本文件专门定义 v3 在多 `gateway` 实例场景下的接入选路、连接承载、路由状态模型与故障恢复策略。

目标是解决两个核心问题：

1. 用户第一次连接时应该进入哪一台 `gateway`
2. 连接建立后，业务服务如何准确找到用户所在 `gateway`

## 2. 角色分工

## 2.1 `edge`

`edge` 是边缘接入与路由决策层。

负责：

- 接收首次接入请求
- 感知可用 `gateway` 实例
- 根据路由策略选择最优 `gateway`
- 返回目标 `wsEndpoint`
- 在重连时优先回原 `gateway`

不负责：

- 不长期代理 WebSocket
- 不持有业务状态真相
- 不承载业务 HTTP 接口

## 2.2 `gateway`

`gateway` 是连接承载层与统一业务出口。

负责：

- 普通 HTTP 透传
- WebSocket 握手
- 连接注册与断开清理
- 用户连接映射
- 本地消息下发
- 推送消费与限流

不负责：

- 不做全局实例发现
- 不承担边缘选路职责
- 不承担业务聚合职责

## 3. 接入模式

v3 第一阶段固定采用：

**`edge` 分配，客户端直连 `gateway`。**

## 3.1 完整握手流程

```text
1. Client -> Edge
   请求获取可连接的 gateway

2. Edge -> Registry / Route Store
   读取 gateway 实例信息与用户历史绑定

3. Edge
   根据路由策略选择目标 gateway

4. Edge -> Client
   返回 gatewayId / wsEndpoint / expiresIn

5. Client -> Gateway
   直接发起 WebSocket 握手

6. Gateway
   建立连接、注册 conn、写入在线映射
```

## 3.2 为什么第一阶段不做 edge 长期反向代理

不采用 “Client -> Edge -> 长期代理 -> Gateway” 的原因：

- edge 长时间持有大量连接，压力更大
- 连接代理增加故障面与排障复杂度
- 网络路径更长
- 成本与复杂度都更高

结论：

- edge 负责分配
- gateway 负责承载

## 4. 路由状态存储模型

## 4.1 gateway registry

由 `etcd` 维护：

```text
gatewayId -> {
  host,
  port,
  region,
  zone,
  status,
  weight,
  capacity
}
```

用途：

- edge 发现可用 gateway
- 按区域、健康、容量做路由选择

## 4.2 user route

由 `redis` 维护：

```text
userId -> {
  gatewayId,
  deviceId?,
  expiresAt
}
```

用途：

- 重连优先路由
- 推送寻址

## 4.3 connection route

由 `redis` 维护：

```text
connId -> gatewayId
```

用途：

- 精确定位连接归属
- 断连时快速清理

## 4.4 gateway online index

由 `redis` 维护：

```text
gatewayId -> set(userId / connId)
```

用途：

- 本地在线用户反查
- 实例清理与排障

## 5. etcd 与 redis 分工

## 5.1 `etcd`

负责：

- 实例注册
- lease 与 keepalive
- 健康状态
- 权重、区域、可用区、容量

适合存储：

- 低频变动的实例元数据
- 可 watch 的注册状态

## 5.2 `redis`

负责：

- 高频在线态写入
- `user -> gateway`
- `conn -> gateway`
- TTL 过期回收

适合存储：

- 高频变化的在线映射与会话态

## 6. edge 路由算法

`edge` 的路由顺序固定为：

1. 已有绑定优先
2. 区域 / 延迟优先
3. 健康优先
4. 负载优先

## 6.1 已有绑定优先

若用户存在有效 `user -> gateway` 映射：

- 优先回原 `gateway`
- 重连保持连接粘性

适用：

- 页面刷新
- 短时断线重连
- 同一设备重复连接

## 6.2 区域优先

若用户无有效绑定：

- 优先同 region
- 再优先同 zone
- 再考虑延迟更低的实例

## 6.3 健康优先

候选实例必须满足：

- etcd 中状态健康
- lease 未过期
- 未被标记为摘除

## 6.4 负载优先

在健康候选中再考虑：

- 当前连接数
- 当前容量水位
- 最近错误率
- 实例权重

## 7. gateway 注册与续约流程

## 7.1 启动注册

`gateway` 启动时：

1. 生成或加载 `gatewayId`
2. 向 `etcd` 写入实例注册信息
3. 绑定 lease
4. 初始化本地连接管理器
5. 初始化 Redis 在线态存储

## 7.2 周期性续约

运行期间：

- 向 `etcd` 持续 keepalive
- 定期刷新本实例基础状态
- 按需刷新本地容量指标

## 7.3 写入在线路由状态

建立连接并完成用户绑定后：

- 写入 `user -> gateway`
- 写入 `conn -> gateway`
- 写入 `gateway -> online index`

以上映射都必须带 TTL。

## 8. 用户连接生命周期

## 8.1 connect

- 客户端从 edge 获得目标 gateway
- 客户端直连 gateway
- gateway 建立连接并注册 `connId`

## 8.2 auth bind

- 用户完成握手鉴权
- gateway 将 `connId` 与 `userId` 绑定
- 写入 `user -> gateway`

## 8.3 heartbeat

- 连接保活
- 刷新在线映射 TTL
- 更新本地活跃状态

## 8.4 disconnect

- 清理 `conn -> gateway`
- 清理 `user -> gateway` 或减少引用
- 从本地在线索引移除

## 8.5 reconnect

- 客户端重新请求 edge
- edge 若命中有效绑定，优先回原 gateway
- 原 gateway 不可用则重新分配

## 9. gateway 故障场景

## 9.1 gateway 不可用

如果 gateway 宕机或失联：

- etcd lease 到期后实例自动摘除
- edge 不再把新请求分配给该实例

## 9.2 旧映射失效

旧的 `user -> gateway`、`conn -> gateway` 映射：

- 依靠 Redis TTL 自动过期
- 避免把请求继续发往失效实例

## 9.3 客户端重连

当客户端重连时：

- 若原实例已摘除，edge 自动重新分配
- 客户端获得新的 `wsEndpoint`

## 10. 推送寻址流程

## 10.1 业务侧发起

业务服务或 worker 产出需要下发的消息。

## 10.2 查询用户路由

推送组件查询：

- `user -> gateway`

若命中，则得到目标 `gatewayId`。

## 10.3 投递到目标 gateway

推送组件将消息投递到目标 `gateway` 的 push channel。

## 10.4 本地下发

目标 `gateway` 收到消息后：

- 命中本地连接
- 通过 realtime 模块下发给对应客户端

## 11. 安全与风控

## 11.1 edge

- 只做接入选路
- 不持有业务状态真相
- 不长期代理连接

## 11.2 gateway

- 握手鉴权
- 连接级限流
- 非法订阅校验
- 空闲连接回收
- 广播范围控制

## 11.3 风险控制

必须具备：

- route TTL
- instance lease
- reconnect fallback
- 连接级限流
- 推送目标校验

## 12. 对后续开发的指导结论

多 `gateway` 实例下，v3 的路由与实时设计最终收口为：

- `edge` 做选路，不做代理承载
- `gateway` 做透传、连接承载与推送
- `etcd` 管实例注册与健康
- `redis` 管用户路由与在线映射
- 路由顺序固定为“已有绑定 > 区域优先 > 健康优先 > 负载优先”
- 所有在线映射必须带 TTL
- 所有实例注册必须带 lease
