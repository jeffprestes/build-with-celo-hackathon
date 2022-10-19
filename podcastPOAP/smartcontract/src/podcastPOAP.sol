// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (utils/Strings.sol)

pragma solidity 0.8.17;

import {ERC1155} from "@solmate/tokens/ERC1155.sol";
import {Strings} from "@openzeppelin/utils/Strings.sol";
import {Ownable} from "@openzeppelin/access/Ownable.sol";

/// @title ERC 1155 Podcast Proof of Audience NFT to give away to Podcast listeners
/// @title PodcastPOAP
/// @author Jeff Prestes
contract PodcastPOAP is ERC1155, Ownable {
    using Strings for uint256;

    mapping(uint256 => string) public uris;
    uint256 public latestEpisodeID = 0;

    /// @notice Add an Episode in the token records
    /// @param _episodeURI the URI of the new episode.
    function addEpisode(string memory _episodeURI) external onlyOwner {
        latestEpisodeID++;
        uris[latestEpisodeID] = _episodeURI;
    }

    /// @notice Mint NFT function.
    /// @param _member the recipient of the new token 
    function mintNFT(address _member) external onlyOwner {
        _mint(_member, latestEpisodeID, 1, "");
    }

    /// @notice Return the URI of the specific episode.
    /// @param _episodeID the episode id. 
    function uri(uint256 _episodeID) public override view virtual returns (string memory) {
        return uris[_episodeID];
    }

    /// @notice Return true if the updates runs OK.
    /// @param _episodeID the episode id. 
    function updateUriEpisode(uint256 _episodeID, string memory _episodeURI) external onlyOwner returns (bool) {
        require(_episodeID <= latestEpisodeID, "invalid episode number");
        uris[_episodeID] = _episodeURI;
        return true;
    }

    function safeTransferFrom(
        address from,
        address to,
        uint256 id,
        uint256 amount,
        bytes calldata data
    ) public override virtual {
    }

    function safeBatchTransferFrom(
        address from,
        address to,
        uint256[] calldata ids,
        uint256[] calldata amounts,
        bytes calldata data
    ) public override virtual { }
}