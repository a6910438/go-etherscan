# go-etherscan
监听ETH区块代币记录

逻辑过程和现有的区块链项目基本类似

首先:  
1.boltdb记录区块高度    
2.根据区块高度遍历交易找系统内已经上线的合约币    
3.找到后更新用户资产 没找到跳过    

![图片说明1](https://github.com/a6910438/go-etherscan/blob/master/1.png)
