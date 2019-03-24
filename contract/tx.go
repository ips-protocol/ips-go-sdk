package contract

type Client interface {
	// 客户端新建上传任务
	// fileHash: 文件hash;
	// fsize: 文件大小;
	// blockNums: 文件分块数;
	// pay: 支持token，默认可传0。
	// 返回交易hash
	NewUploadJob(fileHash string, fsize, blockNums uint64, pay uint64) (tx string, err error)

	// 获取指定文件的存储帐户地址。
	GetStorageAccount(fileHash string) (string, error)
}

type StorageAccount interface {
	// 获取文件基本信息
	GetFileInfo() (fileHash string, fsize, blockNums uint64, err error)

	// 提交文件。获得奖励
	// fileHash:文件hash
	// index:块索引
	// blockHash:块哈唏
	// peerInfo:存储节点信息
	// proof:存储证据
	CommitBlockInfo(fileHash string, index int, blockHash, peerInfo string, proof []byte) (tx string, err error)

	// 根据索引号获取文件块存储信息
	GetBlockInfo(index int) (blockHash, peerInfo string, err error)

	// 根据文件所有块的存储信息
	GetAllBlocksInfo() (blocksHash, peersInfo []string, err error)

	// 下载成功，发放奖励
	DownloadSuccess() (tx string, err error)
}
