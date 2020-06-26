# Vcache
Vcache 基本上模仿了 groupcache 的实现，总体实现上，还是与 groupcache 非常接近的。支持特性有：  单机缓存和基于 HTTP 的分布式缓存 最近最少访问(Least Recently Used, LRU) 缓存策略 使用 Go 锁机制防止缓存击穿 使用一致性哈希选择节点，实现负载均衡 使用 protobuf 优化节点间二进制通信
