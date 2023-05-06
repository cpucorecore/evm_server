pragma solidity 0.6;

contract hello {
    int i;

    function say_hello() public {
        i=i+1;
    }

    function get_i() public view returns(int)  {
        return i;
    }
}
