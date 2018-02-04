pragma solidity ^0.4.18;

contract Election {

    address public chairperson;

    string[] public candidates;

    mapping(uint => uint) public votes;

    mapping(address => bool) public voter;
    mapping(address => bool) public voted;

    function Election() public {
        chairperson = msg.sender;
        candidates.push("Путин, Владимир Владимирович");
        candidates.push("Грудинин, Павел Николаевич");
        candidates.push("Против всех");
    }

    function candidatesCount() public view returns (uint) {
        return candidates.length;
    }

    function addVoters(address[] voters) public {
        require(msg.sender == chairperson);
        for (uint i = 0; i < voters.length; i++) {
            voter[voters[i]] = true;
        }
    }

    function vote(uint id) public {
        require(voter[msg.sender]);
        require(!voted[msg.sender]);
        voted[msg.sender] = true;
        votes[id]++;
    }

    function winner() public view returns (uint) {
        uint max = 0;
        uint winnerID = 0;
        for (uint i = 0; i < candidates.length; i++) {
            if (votes[i] > max) {
                max = votes[i];
                winnerID = i;
            }
        }
        return winnerID;
    }
}