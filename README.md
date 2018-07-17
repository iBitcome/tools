# tools

##luckdraw
抽奖规则：

1. 抽奖活动开始前，选取以太坊未来某个时间之后第一个区块hash(去掉0x前缀)作为种子
2. 活动结束后，对参与到活动中的每一个用户计算score, score计算方式为hash(seed + inviteCode)后取模得到
3. 对活动用户依据score排名分配奖励
