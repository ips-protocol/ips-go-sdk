pragma solidity ^0.4.25;

import "./storageAccount.sol";

// client deposit ETH before upload file
contract StorageDeposit {

    mapping(address => address) storageAccounts;
    event NewUploadJob(address indexed fileAddress, address storageAccount, uint256 deposit);

    function minValue(uint fsize) internal pure returns(uint){
        return 0;
    }

    function newUploadJob(address fileAddress, uint fsize, uint128 block_nums) public payable{
        require(msg.value > minValue(fsize));
        require(storageAccounts[fileAddress] == address(0));

        StorageAccount s = (new StorageAccount).value(msg.value)(fileAddress, block_nums);
        storageAccounts[fileAddress] = address(s);
        emit NewUploadJob(fileAddress, address(s), msg.value);
    }

    function getStorageAccount(address fileAddress) public view returns(address){
        return storageAccounts[fileAddress];
    }
}
