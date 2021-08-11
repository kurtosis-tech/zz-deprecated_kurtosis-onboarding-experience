// SPDX-License-Identifier: MIT
// SOURCE: https://solidity-by-example.org/state-variables/

// compiler version must be greater than or equal to 0.7.6 and less than 0.8.0
pragma solidity >=0.7.0 <0.9.0;

contract SimpleStorage {
    // State variable to store a number
    uint public num;

    // You need to send a transaction to write to a state variable.
    function set(uint _num) public {
        num = _num;
    }

    // You can read from a state variable without sending a transaction.
    function get() public view returns (uint) {
        return num;
    }
}