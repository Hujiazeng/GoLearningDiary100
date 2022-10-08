### 分布式缓存


##### 1.FIFO/LFU/LRU 算法简介

1.1 FIFO 全写First in First out, 先进先出, 最早添加的记录先删除

1.2 LFU 全写Least Frequent Used 最近频繁使用, 计算使用次数, 淘汰最少使用的

1.3 LRU 全写Least Recent Used 最近曾经使用, 使用则将元素排到最前, 淘汰最末尾的

##### 2. LRU 算法实现

双向链表 + 字典

使用go封装好的链表及方法: 

ll := list.New()创建链表
ll.MoveToFront(ele) 将节点移动到队列最前
ll.Remove(ele) 将节点删除
ll.PushFront(ele) 前方插入节点

整体思路: 查询/新增时从字典取出节点, 返回节点值并且将该节点移动到最前, 内存超过时, 删除最末尾节点