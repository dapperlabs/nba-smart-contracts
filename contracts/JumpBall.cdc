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

        init(id: UInt64, creator: Address, startTime: UFix64, gameDuration: UFix64) {
            self.id = id
            self.creator = creator
            self.opponent = nil
            self.startTime = startTime
            self.gameDuration = gameDuration
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
            self.nfts[nftID] <-! nft
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

    // Admin resource for game management
    access(all) resource Admin {
        access(all) fun determineWinner(gameID: UInt64, winnerCap: Capability<&{NonFungibleToken.Collection}>, stats: {UInt64: UInt64}) {
            let game = JumpBall.games[gameID] ?? panic("Game does not exist.")

            let creatorTotal = self.calculateTotal(game: game, stats: stats, player: game.creator)
            let opponentTotal = self.calculateTotal(game: game, stats: stats, player: game.opponent)

            if creatorTotal > opponentTotal {
                self.awardWinner(game: game, winner: game.creator, winnerCap: winnerCap)
            } else if opponentTotal > creatorTotal {
                self.awardWinner(game: game, winner: game.opponent, winnerCap: winnerCap)
            } else {
                self.returnAllNFTs(game: game)
            }
        }

        access(self) fun calculateTotal(game: &Game, address: Address, stats: {UInt64: UInt64}): UInt64 {
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
        JumpBall.userGames <- {}
    }

    // Create a new game
    access(all) fun createGame(creator: Address, startTime: UFix64, gameDuration: UFix64): UInt64 {
        let gameID = JumpBall.nextGameID
        JumpBall.games[gameID] <- create Game(id: gameID, creator: creator, startTime: startTime, gameDuration: gameDuration)
        JumpBall.nextGameID = JumpBall.nextGameID + 1
        JumpBall.addGameForUser(creator, gameID)
        emit GameCreated(gameID: gameID, creator: creator, startTime: startTime)
        return gameID
    }

    // Add an opponent to a game
    access(all) fun addOpponent(gameID: UInt64, opponent: Address) {
        pre {
            game.opponent == nil: "Opponent has already been added."
        }
        let game = JumpBall.games[gameID] ?? panic("Game does not exist.")
        game.opponent = opponent
        JumpBall.addGameForUser(opponent, gameID)
        emit OpponentAdded(gameID: gameID, opponent: opponent)
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
        let currentTime = getCurrentTime()
        let gameStartTime = game.startTime
        let gameDuration = game.gameDuration

        emit TimeoutClaimed(gameID: gameID, claimant: claimant)
        game.returnAllNFTs()
    }

    // Destroy a game and clean up resources
    access(all) fun destroyGame(gameID: UInt64) {
        let game <- JumpBall.games.remove(key: gameID) ?? panic("Game does not exist.")

        // Remove game from userGames mapping
        JumpBall.removeGameForUser(game.creator, gameID)
        JumpBall.removeGameForUser(game.opponent, gameID)

        destroy game
    }

    // Remove a game from a user's list of games
    access(self) fun removeGameForUser(user: Address, gameID: UInt64) {
        if let gameIDs = JumpBall.userGames[user] {
            JumpBall.userGames[user] = gameIDs.filter(fun(id: UInt64): Bool { return id != gameID })
        }
    }

    // Helper function to get the current time (mocked in Cadence)
    access(all) fun getCurrentTime(): UFix64 {
        return UFix64(getCurrentBlock().timestamp)
    }
}
