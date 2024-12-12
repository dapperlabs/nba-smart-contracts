/*
    The JumpBall contract facilitates competitive games by securely managing players' NFTs.
    It creates game-specific instances to isolate gameplay, store NFTs, and enforce game rules.
    NFTs are released to the winner or returned to their original owners based on the game's outcome.

    Authors:
        Corey Humeston: corey.humeston@dapperlabs.com
*/

import NonFungibleToken from 0x631e88ae7f1d7c20

access(all) contract JumpBall {
    // Events
    access(all) event GameCreated(gameID: UInt64, creator: Address, startTime: UFix64)
    access(all) event OpponentAdded(gameID: UInt64, opponent: Address)
    access(all) event NFTDeposited(gameID: UInt64, nftID: UInt64, owner: Address)
    access(all) event NFTReturned(gameID: UInt64, nftID: UInt64, owner: Address)
    access(all) event NFTAwarded(gameID: UInt64, nftID: UInt64, previousOwner: Address, winner: Address)
    access(all) event StatisticSelected(gameID: UInt64, statistic: String, player: Address)
    access(all) event WinnerDetermined(gameID: UInt64, winner: Address)
    access(all) event TimeoutClaimed(gameID: UInt64, claimant: Address)

    // Resource to manage game-specific data
    access(all) resource Game {
        access(all) let id: UInt64
        access(all) let creator: Address
        access(all) var opponent: Address? // Opponent can be added later
        access(all) let startTime: UFix64
        access(all) let gameDuration: UFix64
        access(self) var selectedStatistic: String?
        access(self) var nfts: @{UInt64: NonFungibleToken.NFT}
        access(self) var ownership: {UInt64: Address}
        access(self) var metadata: {String: AnyStruct}
        access(self) let creatorCap: Capability<&{NonFungibleToken.Collection}>
        access(self) let opponentCap: Capability<&{NonFungibleToken.Collection}>?

        init(id: UInt64, creator: Address, startTime: UFix64, gameDuration: UFix64, creatorCap: Capability<&{NonFungibleToken.Collection}>) {
            self.id = id
            self.creator = creator
            self.opponent = nil
            self.startTime = startTime
            self.gameDuration = gameDuration
            self.selectedStatistic = nil
            self.nfts <- {}
            self.ownership = {}
            self.metadata = {}
            self.creatorCap = creatorCap
            self.opponentCap = nil
        }

        access(all) fun setMetadata(key: String, value: AnyStruct) {
            self.metadata[key] = value
        }

        access(all) fun getMetadata(key: String): AnyStruct? {
            return self.metadata[key]
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

        access(all) fun getDepositCapForAddress(owner: Address): Capability<&{NonFungibleToken.Collection}> {
            let depositCap = getAccount(owner).getCapability<&{NonFungibleToken.Collection}>(/public/NFTReceiver)
            if !depositCap.check() {
                panic("Deposit capability for owner is invalid.")
            }
            return depositCap
        }

        // Deposit an NFT into the game
        access(all) fun depositNFT(nft: @NonFungibleToken.NFT, owner: Address) {
            pre {
                JumpBall.getCurrentTime() < self.startTime + self.gameDuration: "Cannot deposit NFT after the game has ended."
            }

            let nftID = nft.id

            // Remove the existing NFT if it exists and destroy it
            if let oldNFT <- self.nfts[nftID] {
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
                let depositCap = JumpBall.getDepositCapForAddress(owner: owner)
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

    // Player resource for game participation
    access(all) resource Player {
        access(all) let address: Address
        access(all) let collectionCap: Capability<&{NonFungibleToken.Collection}>

        init(address: Address, collectionCap: Capability<&{NonFungibleToken.Collection}>) {
            self.address = address
            self.collectionCap = collectionCap
        }

        access(all) fun canCreateGame(): Bool {
            return self.collectionCap.check()
        }

        access(all) fun createPlayer(account: AuthAccount) {
            let collectionCap = account.capabilities.storage.issue<&{NonFungibleToken.Collection}>(/storage/UserNFTCollection)
            account.capabilities.publish(collectionCap, at: /public/UserNFTCollection)

            // Save the Player resource in storage
            account.storage.save<@Player>(<- create Player(
                address: account.address,
                collectionCap: collectionCap
            ), to: /storage/JumpBallPlayer)
        }

        // Store the Player resource in the account's storage
        access(all) fun savePlayer(account: AuthAccount, player: @Player) {
            account.save(<-player, to: /storage/JumpBallPlayer)
            account.link<&Player>(/public/JumpBallPlayer, target: /storage/JumpBallPlayer)
        }

        access(all) fun getPlayer(account: PublicAccount): &Player? {
            return account.getCapability<&Player>(/public/JumpBallPlayer).borrow()
        }
    }

    // Admin resource for game management
    access(all) resource Admin {
        access(all) fun determineWinner(gameID: UInt64, stats: {UInt64: UInt64}) {
            let game = JumpBall.games[gameID] ?? panic("Game does not exist.")

            let creatorTotal = self.calculateTotal(game: game, address: game.creator, stats: stats)
            let opponentTotal = self.calculateTotal(game: game, address: game.opponent, stats: stats)

            if creatorTotal > opponentTotal {
                self.awardWinner(game: game, winner: game.creator, winnerCap: game.creatorCap)
            } else if opponentTotal > creatorTotal {
                let opponentCap = game.opponentCap ?? panic("Opponent capability not found.")
                self.awardWinner(game: game, winner: game.opponent!, winnerCap: opponentCap)
            } else {
                game.returnAllNFTs()
            }
        }

        access(self) fun calculateTotal(game: &Game, address: Address?, stats: {UInt64: UInt64}): UInt64 {
            if address == nil {
                return 0
            }

            return game.nfts.keys.reduce(0, fun(acc: UInt64, key: UInt64): UInt64 {
                let owner = game.ownership[key] ?? panic("Owner not found.")
                if owner == address {
                    return acc + (stats[key] ?? 0)
                }
                return acc
            })
        }

        access(self) fun awardWinner(game: &Game, winner: Address, winnerCap: Capability<&{NonFungibleToken.Collection}>) {
            emit WinnerDetermined(gameID: game.id, winner: winner)
            game.transferAllToWinner(winner: winner, winnerCap: winnerCap)
        }
    }

    // Game registry
    access(self) var nextGameID: UInt64
    access(all) var games: @{UInt64: Game}
    access(all) var userGames: {Address: [UInt64]}

    init() {
        JumpBall.nextGameID = 1
        JumpBall.games <- {}
        JumpBall.userGames = {}

        // Initialize Player resource for contract owner
        let collectionCap = self.account.capabilities.storage.issue<&{NonFungibleToken.Collection}>(/storage/OwnerNFTCollection)
        self.account.capabilities.publish(collectionCap, at: /public/OwnerNFTCollection)

        // Save the Player resource in storage
        self.account.storage.save<@Player>(<- create Player(
            address: self.account.address,
            collectionCap: collectionCap
        ), to: /storage/JumpBallPlayer)
    }

    // Create a new game
    access(all) fun createGame(player: &Player, startTime: UFix64, gameDuration: UFix64): UInt64 {
        pre {
            player.canCreateGame(): "Player does not have a valid collection capability to create a game."
        }

        let gameID = JumpBall.nextGameID
        JumpBall.games[gameID] <- create Game(
            id: gameID,
            creator: player.address,
            startTime: startTime,
            gameDuration: gameDuration,
            creatorCap: player.collectionCap
        )
        JumpBall.nextGameID = JumpBall.nextGameID + 1
        JumpBall.addGameForUser(player.address, gameID)

        emit GameCreated(gameID: gameID, creator: player.address, startTime: startTime)
        return gameID
    }

    // Add an opponent to a game
    access(all) fun addOpponent(gameID: UInt64, player: &Player) {
        let game = JumpBall.games[gameID] ?? panic("Game does not exist.")

        pre {
            game.opponent == nil: "Opponent has already been added."
            player.canCreateGame(): "Player does not have a valid collection capability to be added as an opponent."
        }

        game.opponent = player.address
        game.opponentCap = player.collectionCap
        JumpBall.addGameForUser(player.address, gameID)
        emit OpponentAdded(gameID: gameID, opponent: player.address)
    }

    // Add a game to a user's list of games
    access(self) fun addGameForUser(user: Address, gameID: UInt64) {
        if JumpBall.userGames[user] == nil {
            JumpBall.userGames[user] = []
        }
        JumpBall.userGames[user]?.append(gameID)
    }

    // Retrieve all games for a given user
    access(all) fun getGamesByUser(user: Address): [UInt64] {
        return JumpBall.userGames[user] ?? []
    }

    // Get a reference to a specific game
    access(all) fun getGame(gameID: UInt64): &Game? {
        return &JumpBall.games[gameID] as &Game?
    }

    // Handle timeout: allows users to reclaim their moments
    access(all) fun claimTimeout(gameID: UInt64, claimant: Address) {
        pre {
            JumpBall.getCurrentTime() > JumpBall.games[gameID]?.startTime + JumpBall.games[gameID]?.gameDuration:
            "Game is still in progress."
        }

        let game = JumpBall.games[gameID] ?? panic("Game does not exist.")

        emit TimeoutClaimed(gameID: gameID, claimant: claimant)
        game.returnAllNFTs()
    }

    // Destroy a game and clean up resources
    access(all) fun destroyGame(gameID: UInt64) {
        let game <- JumpBall.games.remove(key: gameID) ?? panic("Game does not exist.")

        // Ensure all NFTs have been withdrawn
        pre {
            game.nfts.isEmpty: "All NFTs must be withdrawn before destroying the game."
        }

        // Remove game from userGames mapping
        JumpBall.removeGameForUser(game.creator, gameID)
        if let opponent = game.opponent {
            JumpBall.removeGameForUser(opponent, gameID)
        }

        // Destroy the game resource
        destroy game
    }

    // Remove a game from a user's list of games
    access(self) fun removeGameForUser(user: Address, gameID: UInt64) {
        if JumpBall.userGames[user] != nil {
            JumpBall.userGames[user] = JumpBall.userGames[user]!.filter(fun(id: UInt64): Bool {
                return id != gameID
            })
        }
    }

    // Helper function to get the current time (mocked in Cadence)
    access(all) fun getCurrentTime(): UFix64 {
        return UFix64(getCurrentBlock().timestamp)
    }
}
