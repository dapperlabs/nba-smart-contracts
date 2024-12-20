/*
    The JumpBall contract facilitates competitive games by securely managing players' NFTs.
    It creates game-specific instances to isolate gameplay, store NFTs, and enforce game rules.
    NFTs are released to the winner or returned to their original owners based on the game's outcome.

    Authors:
        Corey Humeston: corey.humeston@dapperlabs.com
*/

import NonFungibleToken from 0xf8d6e0586b0a20c7

access(all) contract JumpBall {
    // Events
    access(all) event GameCreated(gameID: String, creator: Address, startTime: UFix64)
    access(all) event OpponentAdded(gameID: String, opponent: Address)
    access(all) event NFTDeposited(gameID: String, nftID: UInt64, owner: Address)
    access(all) event NFTReturned(gameID: String, nftID: UInt64, owner: Address)
    access(all) event NFTAwarded(gameID: String, nftID: UInt64, previousOwner: Address, winner: Address)
    access(all) event WinnerDetermined(gameID: String, winner: Address)
    access(all) event TimeoutClaimed(gameID: String, claimant: Address)

    // Player resource for game participation
    access(all) resource Player {
        access(all) let address: Address
        access(all) let collectionCap: Capability<&{NonFungibleToken.Collection}>
        access(all) var metadata: {String: AnyStruct}

        init(address: Address, collectionCap: Capability<&{NonFungibleToken.Collection}>) {
            self.address = address
            self.collectionCap = collectionCap
            self.metadata = {}
        }

        // Helper function for player metadata
        access(all) fun setMetadata(key: String, value: AnyStruct) {
            self.metadata[key] = value
        }

        access(all) fun getMetadata(key: String): AnyStruct? {
            return self.metadata[key]
        }

        // Create a new game
        access(all) fun createGame(gameID: String, startTime: UFix64, gameDuration: UFix64, selectedStatistic: String): String {
            pre {
                self.collectionCap.check(): "Player does not have a valid collection capability to create a game."
            }

            JumpBall.games[gameID] <-! create Game(
                id: gameID,
                creator: self.address,
                startTime: startTime,
                gameDuration: gameDuration,
                selectedStatistic: selectedStatistic,
                creatorCap: self.collectionCap
            )

            JumpBall.addGameForUser(user: self.address, gameID: gameID)

            emit GameCreated(gameID: gameID, creator: self.address, startTime: startTime)
            return gameID
        }

        // Add an opponent to a game
        access(all) fun addOpponent(gameID: String) {
            let gameRef = &JumpBall.games[gameID] as &Game?
                    ?? panic("Game does not exist.")

            if !self.collectionCap.check() {
                panic("Player does not have a valid collection capability to join a game.")
            }

            gameRef.setOpponent(opponentAddress: self.address, collectionCap: self.collectionCap)
            JumpBall.addGameForUser(user: self.address, gameID: gameID)

            emit OpponentAdded(gameID: gameID, opponent: self.address)
        }
    }

    // Resource to manage game-specific data
    access(all) resource Game {
        access(all) let id: String
        access(all) let creator: Address
        access(all) var opponent: Address? // Opponent can be added later
        access(all) let startTime: UFix64
        access(all) let gameDuration: UFix64
        access(self) var selectedStatistic: String
        access(self) var nfts: @{UInt64: {NonFungibleToken.NFT}}
        access(self) var ownership: {UInt64: Address}
        access(self) var metadata: {String: AnyStruct}
        access(self) let creatorCap: Capability<&{NonFungibleToken.Collection}>
        access(self) var opponentCap: Capability<&{NonFungibleToken.Collection}>?

        init(id: String, creator: Address, startTime: UFix64, gameDuration: UFix64, selectedStatistic: String, creatorCap: Capability<&{NonFungibleToken.Collection}>) {
            self.id = id
            self.creator = creator
            self.opponent = nil
            self.startTime = startTime
            self.gameDuration = gameDuration
            self.selectedStatistic = selectedStatistic
            self.nfts <- {}
            self.ownership = {}
            self.metadata = {}
            self.creatorCap = creatorCap
            self.opponentCap = nil
        }

        // Getter for creatorCap
        access(all) fun getCreatorCap(): Capability<&{NonFungibleToken.Collection}> {
            return self.creatorCap
        }

        // Getter for opponentCap
        access(all) fun getOpponentCap(): Capability<&{NonFungibleToken.Collection}>? {
            return self.opponentCap
        }

        access(all) fun setMetadata(key: String, value: AnyStruct) {
            self.metadata[key] = value
        }

        access(all) fun getMetadata(key: String): AnyStruct? {
            return self.metadata[key]
        }

        // Getter function to check if NFTs exist
        access(all) fun hasNoNFTs(): Bool {
            return self.nfts.keys.length == 0
        }

        access(all) fun getNFTKeys(): [UInt64] {
            return self.nfts.keys
        }

        access(all) fun getOwnership(key: UInt64): Address? {
            return self.ownership[key]
        }

        access(all) fun getDepositCapForAddress(owner: Address): Capability<&{NonFungibleToken.Collection}> {
            let acct = getAccount(owner)
            let depositCap = acct.capabilities.get<&{NonFungibleToken.Collection}>(/public/NFTReceiver)

            if !depositCap.check() {
                panic("Deposit capability for owner is invalid.")
            }

            return depositCap
        }

        access(all) fun setOpponent(opponentAddress: Address, collectionCap: Capability<&{NonFungibleToken.Collection}>) {
            pre {
                self.opponent == nil: "Opponent has already been added to this game."
            }

            self.opponent = opponentAddress
            self.opponentCap = collectionCap
        }

        // Deposit an NFT into the game
        access(all) fun depositNFT(nft: @{NonFungibleToken.NFT}, owner: Address) {
            // Time-based check
            if JumpBall.getCurrentTime() >= self.startTime + self.gameDuration {
                panic("Cannot deposit NFT after the game has ended.")
            }

            let nftID = nft.id

            // Safely remove and destroy any existing NFT
            if let oldNFT <- self.nfts.remove(key: nftID) {
                destroy oldNFT
            }

            // Add the new NFT to the dictionary
            self.nfts[nftID] <-! nft

            // Track ownership
            self.ownership[nftID] = owner

            emit NFTDeposited(gameID: self.id, nftID: nftID, owner: owner)
        }

        // Return all NFTs to their original owners
        access(contract) fun returnAllNFTs() {
            let keys = self.nfts.keys
            for key in keys {
                let owner = self.ownership[key] ?? panic("Owner not found for NFT")
                let depositCap = self.getDepositCapForAddress(owner: owner)
                self.returnNFT(nftID: key, depositCap: depositCap)
            }
        }

        // Return a specific NFT to its original owner
        access(contract) fun returnNFT(nftID: UInt64, depositCap: Capability<&{NonFungibleToken.Collection}>) {
            pre {
                self.ownership.containsKey(nftID): "NFT does not exist in this game."
            }
            let ownerAddress = self.ownership.remove(key: nftID)!
            let receiver = depositCap.borrow() ?? panic("Failed to borrow receiver capability.")
            receiver.deposit(token: <-self.nfts.remove(key: nftID)!)
            emit NFTReturned(gameID: self.id, nftID: nftID, owner: ownerAddress)
        }

        // Award all NFTs to the winner
        access(contract) fun transferAllToWinner(winner: Address, winnerCap: Capability<&{NonFungibleToken.Collection}>) {
            let winnerCollection = winnerCap.borrow() ?? panic("Failed to borrow winner capability.")
            let keys = self.nfts.keys
            for key in keys {
                let nft <- self.nfts.remove(key: key) ?? panic("NFT not found.")
                let previousOwner = self.ownership.remove(key: key)!
                winnerCollection.deposit(token: <-nft)
                emit NFTAwarded(gameID: self.id, nftID: key, previousOwner: previousOwner, winner: winner)
            }
        }
    }

    // Admin resource for game management
    access(all) resource Admin {
        access(contract) fun determineWinner(gameID: String, stats: {UInt64: UInt64}) {
            let game = &JumpBall.games[gameID] as &Game?
                ?? panic("Game does not exist.")

            let creatorTotal = self.calculateTotal(game: game, address: game.creator, stats: stats)
            let opponentTotal = self.calculateTotal(game: game, address: game.opponent, stats: stats)

            if creatorTotal > opponentTotal {
                self.awardWinner(game: game, winner: game.creator, winnerCap: game.getCreatorCap())
            } else if opponentTotal > creatorTotal {
                let opponentCap = game.getOpponentCap() ?? panic("Opponent capability not found.")
                self.awardWinner(game: game, winner: game.opponent!, winnerCap: opponentCap)
            } else {
                game.returnAllNFTs()
            }
        }

        access(contract) fun calculateTotal(game: &Game, address: Address?, stats: {UInt64: UInt64}): UInt64 {
            if address == nil {
                return 0
            }

            var total: UInt64 = 0
            let nftKeys = game.getNFTKeys()

            for key in nftKeys {
                let owner = game.getOwnership(key: key) ?? panic("Owner not found.")
                if owner == address {
                    total = total + (stats[key] ?? 0)
                }
            }

            return total
        }

        access(contract) fun awardWinner(game: &Game, winner: Address, winnerCap: Capability<&{NonFungibleToken.Collection}>) {
            emit WinnerDetermined(gameID: game.id, winner: winner)
            game.transferAllToWinner(winner: winner, winnerCap: winnerCap)
        }
    }

    // Factory function to create a new Player resource
    access(all) fun createPlayer(playerAddress: Address, collectionCap: Capability<&{NonFungibleToken.Collection}>): @Player {
        if !collectionCap.check() {
            panic("A valid NFT collection capability is required to create a player.")
        }

        return <- create Player(
            address: playerAddress,
            collectionCap: collectionCap
        )
    }

    access(all) var metadata: {String: AnyStruct}
    access(all) var games: @{String: Game}
    access(all) var userGames: {Address: [String]}
    access(all) let admin: @Admin

    init() {
        self.metadata = {}
        self.admin <- create Admin()
        self.games <- {}
        self.userGames = {}

        // Initialize Player resource for contract owner
        let collectionCap = self.account.capabilities.storage.issue<&{NonFungibleToken.Collection}>(/storage/OwnerNFTCollection)
        self.account.capabilities.publish(collectionCap, at: /public/OwnerNFTCollection)

        // Save the Player resource in storage
        self.account.storage.save<@Player>(<- create Player(
            address: self.account.address,
            collectionCap: collectionCap
        ), to: /storage/JumpBallPlayer)
    }

    // Helper functions for contract metadata
    access(all) fun setMetadata(key: String, value: AnyStruct) {
        self.metadata[key] = value
    }

    access(all) fun getMetadata(key: String): AnyStruct? {
        return self.metadata[key]
    }

    // Add a game to a user's list of games
    access(self) fun addGameForUser(user: Address, gameID: String) {
        if JumpBall.userGames[user] == nil {
            JumpBall.userGames[user] = []
        }
        JumpBall.userGames[user]?.append(gameID)
    }

    // Retrieve all games for a given user
    access(all) fun getGamesByUser(user: Address): [String] {
        return JumpBall.userGames[user] ?? []
    }

    // Get a reference to a specific game
    access(all) fun getGame(gameID: String): &Game? {
        return &JumpBall.games[gameID] as &Game?
    }

    access(all) fun gameExists(gameID: String): Bool {
        return JumpBall.games[gameID] != nil
    }

    // Handle timeout: allows users to reclaim their moments
    access(all) fun claimTimeout(gameID: String, claimant: Address) {
        let gameRef = &JumpBall.games[gameID] as &Game?
            ?? panic("Game does not exist.")

        if JumpBall.getCurrentTime() < gameRef.startTime + gameRef.gameDuration {
            panic("Game is still in progress.")
        }

        emit TimeoutClaimed(gameID: gameID, claimant: claimant)
        gameRef.returnAllNFTs()
    }

    // Destroy a game and clean up resources
    access(all) fun destroyGame(gameID: String) {
        // Safely borrow a reference to the game
        let gameRef = &JumpBall.games[gameID] as &Game?
            ?? panic("Game does not exist.")

        if !gameRef.hasNoNFTs() {
            panic("All NFTs must be withdrawn before destroying the game.")
        }

        // Remove the game from the dictionary
        let game <- JumpBall.games.remove(key: gameID)
            ?? panic("Game does not exist.")

        // Remove the game from the user's games list
        JumpBall.removeGameForUser(user: game.creator, gameID: gameID)
        if let opponent = game.opponent {
            JumpBall.removeGameForUser(user: opponent, gameID: gameID)
        }

        // Destroy the game resource
        destroy game
    }

    // Remove a game from a user's list of games
    access(self) fun removeGameForUser(user: Address, gameID: String) {
        if JumpBall.userGames[user] != nil {
            JumpBall.userGames[user] = JumpBall.userGames[user]!.filter(view fun(id: String): Bool {
                return id != gameID
            })
        }
    }

    // Helper function to get the current time (mocked in Cadence)
    access(all) fun getCurrentTime(): UFix64 {
        return UFix64(getCurrentBlock().timestamp)
    }

    // Securely call determine winner from the admin resource
    access(all) fun determineWinner(gameID: String, stats: {UInt64: UInt64}) {
        self.admin.determineWinner(gameID: gameID, stats: stats)
    }
}
