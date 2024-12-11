/*
    The JumpBall contract facilitates competitive games by securely managing players' NFTs.
    It creates game-specific instances to isolate gameplay, store NFTs, and enforce game rules.
    NFTs are released to the winner or returned to their original owners based on the game's outcome.

    Authors:
        Corey Humeston: corey.humeston@dapperlabs.com
*/

import NonFungibleToken from 0xNFTADDRESS

access(all) contract JumpBall {
    // Events
    access(all) event GameCreated(gameID: UInt64, creator: Address)
    access(all) event NFTDeposited(gameID: UInt64, nftID: UInt64, owner: Address)
    access(all) event NFTReturned(gameID: UInt64, nftID: UInt64, owner: Address)
    access(all) event NFTAwarded(gameID: UInt64, nftID: UInt64, previousOwner: Address, winner: Address)
    access(all) event StatisticSelected(gameID: UInt64, statistic: String, player: Address)

    // Resource to manage game-specific data
    access(all) resource Game {
        pub let id: UInt64
        pub let creator: Address
        pub let opponent: Address
        access(self) var selectedStatistic: String?
        access(self) var nfts: @{UInt64: NonFungibleToken.NFT}
        access(self) var ownership: {UInt64: Address}

        init(id: UInt64, creator: Address) {
            self.id = id
            self.creator = creator
            self.opponent = Address.zero()
            self.selectedStatistic = nil
            self.nfts <- {}
            self.ownership <- {}
        }

        access(all) fun selectStatistic(statistic: String, player: Address) {
            pre {
                self.selectedStatistic == nil: "Statistic has already been selected for this game."
            }
            self.selectedStatistic = statistic
            emit StatisticSelected(gameID: self.id, statistic: statistic, player: player)
        }

        // Retrieve the selected statistic
        access(all) fun getStatistic(): String? {
            return self.selectedStatistic
        }

        // Deposit an NFT into the game
        access(all) fun depositNFT(nft: @NonFungibleToken.NFT, owner: Address) {
            let nftID = nft.id
            self.nfts[nftID] <-! nft
            self.ownership[nftID] = owner
            emit NFTDeposited(gameID: self.id, nftID: nftID, owner: owner)
        }

        // Award all NFTs to the winner
        access(contract) fun transferAllToWinner(winnerCap: Capability<&{NonFungibleToken.Collection}>) {
            let winnerCollection = winnerCap.borrow() ?? panic("Failed to borrow winner capability.")

            for (nftID, nft) in self.nfts {
                let previousOwner = self.ownership.remove(key: nftID)!
                winnerCollection.deposit(token: <-nft)
                emit NFTAwarded(gameID: self.id, nftID: nftID, previousOwner: previousOwner, winner: winnerCap.address)
            }

            // Clear all NFTs after awarding them to the winner
            self.nfts = {}
        }

        // Return an NFT to its original owner
        access(contract) fun returnNFT(nftID: UInt64, depositCap: Capability<&{NonFungibleToken.Collection}>) {
            pre {
                self.ownership.containsKey(nftID): "NFT does not exist in this game."
            }

            let ownerAddress = self.ownership.remove(key: nftID)!
            let receiver = depositCap.borrow() ?? panic("Failed to borrow receiver capability.")
            receiver.deposit(token: <-self.nfts.remove(key: nftID)!)
            emit NFTReturned(gameID: self.id, nftID: nftID, owner: ownerAddress)
        }

        destroy() {
            destroy self.nfts
        }
    }

    // Game registry
    access(self) var nextGameID: UInt64
    access(all) var games: @{UInt64: Game}

    // Mapping of users to their associated gameIDs
    access(all) var userGames: {Address: [UInt64]}

    init() {
        self.nextGameID = 1
        self.games <- {}
        self.userGames <- {}
    }

    // Create a new game
    access(all) fun createGame(creator: Address, opponent: Address): UInt64 {
        let gameID = self.nextGameID
        self.games[gameID] <- create Game(id: gameID, creator: creator, opponent: opponent)
        self.nextGameID = self.nextGameID + 1

        // Update userGames mapping
        self.addGameForUser(creator, gameID)
        self.addGameForUser(opponent, gameID)

        emit GameCreated(gameID: gameID, creator: creator)
        return gameID
    }

    // Add a game to a user's list of games
    access(self) fun addGameForUser(user: Address, gameID: UInt64) {
        if self.userGames[user] == nil {
            self.userGames[user] = []
        }
        self.userGames[user]?.append(gameID)
    }

    // Retrieve all games for a given user
    access(all) fun getGamesByUser(user: Address): [UInt64] {
        return self.userGames[user] ?? []
    }

    // Get a reference to a specific game
    access(all) fun getGame(gameID: UInt64): &Game? {
        return &self.games[gameID] as &Game?
    }

    // Destroy a game and clean up resources
    access(all) fun destroyGame(gameID: UInt64) {
        let game <- self.games.remove(key: gameID) ?? panic("Game does not exist.")

        // Remove game from userGames mapping
        self.removeGameForUser(game.creator, gameID)
        self.removeGameForUser(game.opponent, gameID)

        destroy game
    }

    // Remove a game from a user's list of games
    access(self) fun removeGameForUser(user: Address, gameID: UInt64) {
        if let gameIDs = self.userGames[user] {
            self.userGames[user] = gameIDs.filter { $0 != gameID }
        }
    }

    destroy() {
        destroy self.games
    }
}
