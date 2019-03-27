pragma solidity ^0.4.25;


contract StorageAccount{

    struct Block {
        bytes32 hash;
        bytes32 peerInfo;
        address beneficiary;
    }

    address fileAddress;
    mapping(uint => Block[]) blocks;
    mapping(address => uint256) balance;
    uint128 uploadedBlockNums;
    uint128 blockNums;
    uint downloadTotal;
    uint initBalance;

    constructor(address _fileAddress, uint128 _block_nums) public payable{
        fileAddress = _fileAddress;
        blockNums = _block_nums;
        initBalance = msg.value;
    }

    // Save one block of the file and get some rewards.
    function CommitBlockInfo(address _fileAddress, uint index, bytes32 hash, bytes32 peerInfo, string proof) public {

        require(index < blockNums);
        require(_fileAddress == fileAddress);
        require(blocks[index].length < 3);
        // check dup
        for(uint i = 0; i < blocks[index].length; i++){
            if (blocks[index][i].peerInfo == peerInfo && blocks[index][i].beneficiary == msg.sender){
                revert();
            }
        }

        uploadedBlockNums++;
        blocks[index].push(Block(hash, peerInfo, msg.sender));
        // rewards 20% for 3 miners
        msg.sender.transfer(initBalance / (5*uploadedBlockNums*3));
    }

    function getBlockInfo(uint index) public view returns(bytes32 blockHash, bytes32 peerInfo){
        // TODO: return all peersInfo

        blockHash = blocks[index][0].hash;
        peerInfo = blocks[index][0].peerInfo;
        return;
    }

    // 获取文件所有块的存储信息
	function getAllBlocksInfo() public view returns (bytes32[] blocksHash, bytes32[] peersInfo){
	    return;
	}

	function DownloadSuccess() public  {

	    downloadTotal++;
	    if (downloadTotal > 10) {
	        return;
	    }
	    // 80%为下载奖励，分10次发放。

	    uint total = 0;
	    for(uint i = 0; i < blockNums; i++){
	        total += blocks[i].length;
	    }
	    uint reward = initBalance*4/(5*10*total);

	    for(i = 0; i < blockNums; i++){
	        for( uint j = 0; j < blocks[i].length; j ++){
	            blocks[i][j].beneficiary.transfer(reward);
	        }
	    }
	}
}