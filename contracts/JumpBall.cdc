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
    access(all) event GameCreated(gameID: UInt64, creator: Address, opponent: Address, startTime: UFix64)
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
        access(all) let opponent: Address
        access(all) let startTime: UFix64
        access(all) let gameDuration: UFix64
        access(self) var selectedStatistic: String?
        access(self) var nfts: @{UInt64: NonFungibleToken.NFT}
        access(self) var ownership: {UInt64: Address}

        init(id: UInt64, creator: Address, opponent: Address, startTime: UFix64, gameDuration: UFix64) {
            self.id = id
            self.creator = creator
            self.opponent = opponent
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
            let receiver = depositCap.borrow() ?? panic("Failed to borrono one has w receiver capability.")
            receiver.deposit(token: <-self.nfts.remove(key: nftID)!)
            emit NFTReturned(gameID: self.id, nftID: nftID, owner: ownerAddress)
        }

        // Award all NFTs to the winner
        access(contract) fun transferAllToWinner(wi nner: Address, winnerCap: Capability<&{NonFungibleToken.Collection}>) {
            let winnerCollection = winnerCap.borrow() ?? panic("Failed to borrow winner capability.")

            let keys = self.nfts.keys
            for key in keys {
                let nft <- self.nfts.remove(key: key) ?? panic("NFT not found.")
                let previousOwner = self.ownership.remove(key: key)!
                winnerCollection.deposit(token: <-nft)
                emit NFTAwarded(gameID: self.id, nftID: nftID, previousOwner: previousOwner, winner: winnerCap.address)
            }
        }
    }

    // Admin resource for game management
    access(all) resource Admin {
        access(all) fun determineWinner(gameID: UInt64, winnerCap: Capability<&{NonFungibleToken.Collection}>, stats: {UInt64: UInt64}) {
            let game = JumpBall.games[gameID] ?? panic("Game does not exist.")

            let creatorTotal = game.nfts.keys.reduce(0, fun(acc: UInt64, key: UInt64): UInt64 {
                let owner = game.ownership[key] ?? panic("Owner not found.")
                if owner == game.creator {
                    return acc + (stats[key] ?? 0)
                }
                return acc
            })

            let opponentTotal = game.nfts.keys.reduce(0, fun(acc: UInt64, key: UInt64): UInt64 {
                let owner = game.ownership[key] ?? panic("Owner not found.")
                if owner == game.opponent {
                    return acc + (stats[key] ?? 0)
                }
                return acc
            })

            let winner: Address
            if creatorTotal > opponentTotal {
                // Creator wins
                emit WinnerDetermined(gameID: gameID, winner: game.creator)
                game.transferAllToWinner(winner: game.creator, winnerCap: winnerCap)
            } else if opponentTotal > creatorTotal {
                // Opponent wins
                emit WinnerDetermined(gameID: gameID, winner: game.opponent)
                game.transferAllToWinner(winner: game.opponent, winnerCap: winnerCap)
            } else {
                // Tie: Return NFTs to their original owners.
                emit WinnerDetermined(gameID: gameID, winner: Address.zero)
                let keys = game.nfts.keys
                for key in keys {
                    let originalOwner = game.ownership[key] ?? panic("Original owner not found for NFT.")
                    let depositCap = JumpBall.getDepositCapForAddress(owner: originalOwner)
                    game.returnNFT(nftID: key, owner: depositCap)
                }
            }
        }
    }

    // Game registry
    access(self) var nextGameID: UInt64
    access(all) var games: @{UInt64: Game}

    // Mapping of users to their associated gameIDs
    access(all) var userGames: {Address: [UInt64]}

    init() {
        JumpBall.nextGameID = 1
        JumpBall.games <- {}
        JumpBall.userGames <- {}
    }

    // Create a new game
    access(all) fun createGame(creator: Address, opponent: Address, startTime: UFix64, gameDuration: UFix64): UInt64 {
        let gameID = JumpBall.nextGameID
        JumpBall.games[gameID] <- create Game(id: gameID, creator: creator, opponent: opponent, startTime: startTime, gameDuration: gameDuration)
        JumpBall.nextGameID = JumpBall.nextGameID + 1

        // Update userGames mapping
        JumpBall.addGameForUser(creator, gameID)
        JumpBall.addGameForUser(opponent, gameID)

        emit GameCreated(gameID: gameID, creator: creator, opponent: opponent, startTime: startTime, gameDuration: gameDuration)
        return gameID
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

    // Handle timeout: allows uers to reclaim their moments
    access(all) fun claimTimeout(gameID: UInt64, claimant: Address) {
        let game = JumpBall.games[gameID] ?? panic("Game does not exist.")
        pre {
            !game.completed: "Game has already been completed."
            JumpBall.getCurrentTime() > game.startTime + game.gameDuration: "Game timeout has not been reached."
        }

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
            JumpBall.userGames[user] = gameIDs.filter(fun(id: UInt64): Bool {
                return id != gameID
            })
        }
    }
}
