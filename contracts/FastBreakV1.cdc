/*
      _____                 __    ___.                         __
    _/ ____\____    _______/  |_  \_ |_________   ____ _____  |  | __
    \   __\\__  \  /  ___/\   __\  | __ \_  __ \_/ __ \\__  \ |  |/ /
     |  |   / __ \_\___ \  |  |    | \_\ \  | \/\  ___/ / __ \|    <
     |__|  (____  /____  > |__|    |___  /__|    \___  >____  /__|_ \
                \/     \/              \/            \/     \/     \/

    fast break game contract & oracle
    micro coder: jer ahrens <jer.ahrens@dapperlabs.com>

*/

import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS

/// Game & Oracle Contract for Fast Break V1
///
pub contract FastBreakV1: NonFungibleToken {

    /// Contract events
    ///
    pub event ContractInitialized()

    pub event Withdraw(id: UInt64, from: Address?)

    pub event Deposit(id: UInt64, to: Address?)

    pub event FastBreakPlayerCreated(
        id: UInt64,
        playerName: String
    )

    pub event FastBreakRunCreated(
        id: String,
        name: String,
        runStart: UInt64,
        runEnd: UInt64,
        fatigueModeOn: Bool
    )

    pub event FastBreakRunStatusChange(id: String, newRawStatus: UInt8)

    pub event FastBreakGameCreated(
        id: String,
        name: String,
        fastBreakRunID: String,
        submissionDeadline: UInt64,
        numPlayers: UInt64
    )

    pub event FastBreakGameStatusChange(id: String, newRawStatus: UInt8)

    pub event FastBreakNFTBurned(id: UInt64, serialNumber: UInt64)

    pub event FastBreakGameTokenMinted(
        id: UInt64,
        fastBreakGameID: String,
        serialNumber: UInt64,
        mintingDate: UInt64,
        topShots: [UInt64],
        mintedTo: UInt64
    )

    pub event FastBreakGameSubmissionUpdated(
        playerId: UInt64,
        fastBreakGameID: String,
        topShots: [UInt64],
    )

    pub event FastBreakGameWinner(
        playerId: UInt64,
        submittedAt: UInt64,
        fastBreakGameID: String,
        topShots: [UInt64]
    )

    pub event FastBreakGameStatAdded(
        fastBreakGameID: String,
        name: String,
        type: UInt8,
        valueNeeded: UInt64
    )

    /// Named Paths
    ///
    pub let CollectionStoragePath:      StoragePath
    pub let CollectionPublicPath:       PublicPath
    pub let OracleStoragePath:          StoragePath
    pub let OraclePrivatePath:          PrivatePath
    pub let PlayerStoragePath:          StoragePath
    pub let PlayerPrivatePath:          PrivatePath

    /// Contract variables
    ///
    pub var totalSupply:        UInt64
    pub var nextPlayerId:        UInt64

    /// Game Enums
    ///

    /// A game of Fast Break has the following status transitions
    ///
    pub enum GameStatus: UInt8 {
        pub case SCHEDULED /// Game is schedules but closed for submission
        pub case OPEN /// Game is open for submission
        pub case STARTED /// Game has started
        pub case CLOSED /// Game is over and rewards are being distributed
    }

    /// A Fast Break Run has the following status transitions
    ///
    pub enum RunStatus: UInt8 {
        pub case SCHEDULED
        pub case RUNNING /// The first Fast Break game of the run has started
        pub case CLOSED /// The last Fast Break game of the run has ended
    }

    /// A Fast Break Statistic can be met by an individual or group of top shots
    ///
    pub enum StatisticType: UInt8 {
        pub case INDIVIDUAL /// Each top shot must meet or exceed this statistical value
        pub case CUMMULATIVE /// All top shots in the submission must meet or exceed this statistical value
    }

    /// Metadata Dictionaries
    ///
    access(self) let fastBreakRunByID:      {String: FastBreakRun}
    access(self) let fastBreakGameByID:     {String: FastBreakGame}
    access(self) let fastBreakPlayerByID:   {UInt64: PlayerData}
    access(self) let playerAccountMapping:  {UInt64: Address}
    access(self) let accountPlayerMapping:  {Address: UInt64}

    /// A top-level Fast Break Run, the container for Fast Break Games
    /// A Fast Break Run contains many Fast Break games & is a mini-season.
    /// Fatigue mode applies submission limitations for the off-chain version of the game
    /// Fatigue mode limits top shot usage by tier. 4 uses legendary. 2 uses rare. 1 use other.
    ///
    pub struct FastBreakRun {
        pub let id: String /// The off-chain uuid of the Fast Break Run
        pub let name: String /// The name of the Run (R0, R1, etc)
        pub var status: FastBreakV1.RunStatus /// The status of the run
        pub let runStart: UInt64 /// The block timestamp starting the run
        pub let runEnd: UInt64 /// The block timestamp ending the run
        pub let runWinCount: {UInt64: UInt64} /// win count by playerId
        pub let fatigueModeOn: Bool /// Fatigue mode is a game rule limiting usage of top shots by tier

        init (id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool) {
            if let fastBreakRun = FastBreakV1.fastBreakRunByID[id] {
                self.id = fastBreakRun.id
                self.name = fastBreakRun.name
                self.status = fastBreakRun.status
                self.runStart = fastBreakRun.runStart
                self.runEnd = fastBreakRun.runEnd
                self.runWinCount = fastBreakRun.runWinCount
                self.fatigueModeOn = fastBreakRun.fatigueModeOn
            } else {
                self.id = id
                self.name = name
                self.status = FastBreakV1.RunStatus.SCHEDULED
                self.runStart = runStart
                self.runEnd = runEnd
                self.runWinCount = {}
                self.fatigueModeOn = fatigueModeOn
            }
        }

        /// Update status of the Fast Break Run
        ///
        access(contract) fun updateStatus(status: FastBreakV1.RunStatus) { self.status = status }

        /// Write a new win to the Fast Break Run runWinCount
        ///
        access(contract) fun incrementRunWinCount(playerId: UInt64) {
            let runWinCount = self.runWinCount
            runWinCount[playerId] = (runWinCount[playerId] ?? 0) + 1
        }
    }

    /// Get a Fast Break Run by Id
    ///
    pub fun getFastBreakRun(id: String): FastBreakV1.FastBreakRun? {
        return FastBreakV1.fastBreakRunByID[id]
    }

    /// A single Game of Fast Break
    /// A Fast Break is played on any day NBA games are scheduled
    /// It is the intention of this contract to allow private & public Fast Break games
    /// A private Fast Break is visible on-chain but is restricted to private accounts
    /// A public Fast Break can be played by custodial and non-custodial accounts
    ///
    pub struct FastBreakGame {
        pub let id: String /// The off-chain uuid of the Fast Break
        pub let name: String /// The name of the Fast Break (eg FB0, FB1, FB2)
        pub let submissionDeadline: UInt64 /// The block timestamp restricting submission to the Fast Break
        pub let numPlayers: UInt64 /// The number of top shots a player should submit to the Fast Break
        pub var status: FastBreakV1.GameStatus /// The game status
        pub var winner: UInt64 /// The playerId of the winner of Fast Break
        pub var submissions: {UInt64: FastBreakV1.FastBreakSubmission} /// Map of player submission to the Fast Break
        pub let fastBreakRunID: String /// The off-chain uuid of the Fast Break Run containing this Fast Break
        pub var stats: [FastBreakStat] /// The NBA statistical requirements for this Fast Break

        init (
            id: String,
            name: String,
            fastBreakRunID: String,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        ) {
            if let fb = FastBreakV1.fastBreakGameByID[id] {
                self.id = fb.id
                self.name = fb.name
                self.submissionDeadline = fb.submissionDeadline
                self.numPlayers = fb.numPlayers
                self.status = fb.status
                self.winner = fb.winner
                self.submissions = fb.submissions
                self.fastBreakRunID = fb.fastBreakRunID
                self.stats = fb.stats
            } else {
                self.id = id
                self.name = name
                self.submissionDeadline = submissionDeadline
                self.numPlayers = numPlayers
                self.status = FastBreakV1.GameStatus.SCHEDULED
                self.submissions = {}
                self.fastBreakRunID = fastBreakRunID
                self.stats = []
                self.winner = 0
            }
        }

        /// Get a account's active Fast Break Submission
        ///
        pub fun getFastBreakSubmissionByPlayerId(playerId: UInt64): FastBreakV1.FastBreakSubmission? {
            let fastBreakSubmissions = self.submissions

            return fastBreakSubmissions[playerId]
        }

        /// Add a statistic to the Fast Break during game creation
        ///
        access(contract) fun addStat(stat: FastBreakV1.FastBreakStat) {
            self.stats.append(stat)
        }

        /// Update status and winner of a Fast Break
        ///
        access(contract) fun update(status: FastBreakV1.GameStatus, winner: UInt64) {
            self.status = status
            self.winner = winner
        }

        /// Submit a Fast Break
        ///
        access(contract) fun submitFastBreak(submission: FastBreakV1.FastBreakSubmission) {
            pre {
                FastBreakV1.isValidSubmission(submissionDeadline: self.submissionDeadline) : "Submission missed deadline"
            }

            self.submissions[submission.playerId] = submission
        }

        /// Update a Fast Break with new topshot moments
        ///
        access(contract) fun updateFastBreakTopshots(playerId: UInt64, topshotMoments: [UInt64]) {
            pre {
                FastBreakV1.isValidSubmission(submissionDeadline: self.submissionDeadline) : "Submission update missed deadline"
            }

            let submission = &self.submissions[playerId] as &FastBreakV1.FastBreakSubmission?
                ?? panic("Could not find submission for playerId: ".concat(playerId.toString()))

            submission.updateTopshots(topshotMomentIds: topshotMoments)
        }

        /// Update the Fast Break score of an account
        ///
        access(contract) fun updateScore(playerId: UInt64, points: UInt64, win: Bool): Bool {
            let submissions = self.submissions

            let submission: FastBreakV1.FastBreakSubmission = submissions[playerId]
                ?? panic("Unable to find fast break submission for playerId: ".concat(playerId.toString()))

            let isPrevSubmissionWin = submission.win

            submission.setPoints(points: points, win: win)

            self.submissions[playerId] = submission

            if win && !isPrevSubmissionWin {
                return true
            }

            return false
        }
    }

    /// Validate Fast Break Submission
    ///
    pub fun isValidSubmission(submissionDeadline: UInt64): Bool {
        return submissionDeadline > UInt64(getCurrentBlock().timestamp) 
    }

    /// Get a Fast Break Game by Id
    ///
    pub fun getFastBreakGame(id: String): FastBreakV1.FastBreakGame? {
        return FastBreakV1.fastBreakGameByID[id]
    }

    /// Get the game stats of a Fast Break
    ///
    pub fun getFastBreakGameStats(id: String): [FastBreakV1.FastBreakStat] {
        if let fastBreak = FastBreakV1.getFastBreakGame(id: id) {
            return fastBreak.stats
        }
        return []
    }

    /// Get a Fast Break account by playerId
    ///
    pub fun getFastBreakPlayer(id: UInt64): Address? {
        return FastBreakV1.playerAccountMapping[id]
    }

    /// A statistical structure used in Fast Break Games
    /// This structure names the NBA statistic top shots must match or exceed
    /// An example is points as the statistic and 30 as the value
    /// A top shot or group of top shots must meet or exceed 30 points
    ///
    pub struct FastBreakStat {
        pub let name: String
        pub let type: FastBreakV1.StatisticType
        pub let valueNeeded: UInt64

        init (
            name: String,
            type: FastBreakV1.StatisticType,
            valueNeeded: UInt64
        ) {
            self.name = name
            self.type = type
            self.valueNeeded = valueNeeded
        }
    }

    /// An account submission to a Fast Break
    ///
    pub struct FastBreakSubmission {
        pub let playerId: UInt64
        pub var submittedAt: UInt64
        pub let fastBreakGameID: String
        pub var topShots: [UInt64]
        pub var points: UInt64
        pub var win: Bool

        init (
            playerId: UInt64,
            fastBreakGameID: String,
            topShots: [UInt64],
        ) {
            self.playerId = playerId
            self.fastBreakGameID = fastBreakGameID
            self.topShots = topShots
            self.submittedAt = UInt64(getCurrentBlock().timestamp)
            self.points = 0
            self.win = false
        }

        /// Set the points of a submission
        ///
        access(contract) fun setPoints(points: UInt64, win: Bool) {
            self.points = points
            self.win = win
        }

        access(contract) fun updateTopshots(topshotMomentIds: [UInt64]) {
            self.topShots = topshotMomentIds
        }
    }

    /// Resource for playing Fast Break
    /// The Fast Break Player plays the game & mints game tokens
    ///
    pub resource Player: FastBreakPlayer, NonFungibleToken.INFT {

        pub let id: UInt64
        pub let playerName: String      /// username
        pub var tokensMinted: UInt64    /// num games played

        access(contract) var gameTokensPlayed: [UInt64]

        init(playerName: String) {
            self.id = FastBreakV1.nextPlayerId
            self.playerName = playerName
            self.gameTokensPlayed = []
            self.tokensMinted = 0

            FastBreakV1.fastBreakPlayerByID[self.id] = PlayerData(playerName: playerName)
        }

        /// Play the game of Fast Break with an array of Top Shots
        /// Each account must own a top shot collection to play fast break
        ///
        pub fun play(
            fastBreakGameID: String,
            topShots: [UInt64]
        ): @FastBreakV1.NFT {
            pre {
                FastBreakV1.fastBreakGameByID.containsKey(fastBreakGameID): "No such fast break game with gameId: ".concat(fastBreakGameID)
            }

            /// Update player address mapping
            if let ownerAddress = self.owner?.address {
                FastBreakV1.playerAccountMapping[self.id] = ownerAddress
                FastBreakV1.accountPlayerMapping[ownerAddress] = self.id
            }

            /// Validate Top Shots
            let acct = getAccount(self.owner?.address!)
            let collectionRef = acct.getCapability(/public/MomentCollection)
                .borrow<&{TopShot.MomentCollectionPublic}>() ?? panic("Player does not have top shot collection")

            /// Must own Top Shots to play Fast Break
            /// more efficient to borrow ref than to loop
            ///
            for flowId in topShots {
                let topShotRef = collectionRef.borrowMoment(id: flowId)
                    ?? panic("Top shot not owned in collection with flowId: ".concat(flowId.toString()))
            }

            let fastBreakGame = (&FastBreakV1.fastBreakGameByID[fastBreakGameID] as &FastBreakV1.FastBreakGame?)
                 ?? panic("Fast break does not exist with gameId: ".concat(fastBreakGameID))

            /// Cannot mint two tokens for the same Fast Break
            let existingSubmission = fastBreakGame.getFastBreakSubmissionByPlayerId(playerId: self.id)
            if existingSubmission != nil {
                panic("Account already submitted to fast break with playerId: ".concat(self.id.toString()))
            }

            let fastBreakSubmission = FastBreakV1.FastBreakSubmission(
                playerId: self.id,
                fastBreakGameID: fastBreakGameID,
                topShots: topShots
            )

            fastBreakGame.submitFastBreak(submission: fastBreakSubmission)

            let fastBreakNFT <- create NFT(
                fastBreakGameID: fastBreakGameID,
                serialNumber: self.tokensMinted + 1,
                topShots: topShots,
                mintedTo: self.id
            )

            self.tokensMinted = self.tokensMinted + 1
            self.gameTokensPlayed.append(fastBreakNFT.id)

            emit FastBreakGameTokenMinted(
                id: fastBreakNFT.id,
                fastBreakGameID: fastBreakNFT.fastBreakGameID,
                serialNumber: fastBreakNFT.serialNumber,
                mintingDate: fastBreakNFT.mintingDate,
                topShots: fastBreakNFT.topShots,
                mintedTo: fastBreakNFT.mintedTo
            )

            FastBreakV1.totalSupply = FastBreakV1.totalSupply + 1
            return <- fastBreakNFT
        }

        /// Update FastBreak Game Submission with an array of Top Shots
        /// Each account must have a submission before being able to update
        ///
        pub fun updateSubmission(
            fastBreakGameID: String,
            topShots: [UInt64]
        ) {
            pre {
                FastBreakV1.fastBreakGameByID.containsKey(fastBreakGameID): "No such fast break game with gameId: ".concat(fastBreakGameID)
            }

            /// Update player address mapping
            if let ownerAddress = self.owner?.address {
                FastBreakV1.playerAccountMapping[self.id] = ownerAddress
                FastBreakV1.accountPlayerMapping[ownerAddress] = self.id
            }

            /// Validate Top Shots
            let acct = getAccount(self.owner?.address!)
            let collectionRef = acct.getCapability(/public/MomentCollection)
                .borrow<&{TopShot.MomentCollectionPublic}>() ?? panic("Player does not have top shot collection")

            /// Must own Top Shots to play Fast Break
            /// more efficient to borrow ref than to loop
            ///
            for flowId in topShots {
                let topShotRef = collectionRef.borrowMoment(id: flowId)
                    ?? panic("Top shot not owned in collection with flowId: ".concat(flowId.toString()))
            }

            let fastBreakGame = (&FastBreakV1.fastBreakGameByID[fastBreakGameID] as &FastBreakV1.FastBreakGame?)
                ?? panic("Fast break does not exist with gameId: ".concat(fastBreakGameID))

            /// Check that the user has a submission for Fast Break game we can update
            let pastSubmission = fastBreakGame.getFastBreakSubmissionByPlayerId(playerId: self.id)
                ?? panic("Account already with playerID: ".concat(self.id.toString())
                    .concat(" has not played FastBreak with ID: ".concat(fastBreakGameID)))

            fastBreakGame.updateFastBreakTopshots(playerId: self.id, topshotMoments: topShots)

            // Get the updated submission with new topshot moment Ids
            let updatedSubmission = fastBreakGame.getFastBreakSubmissionByPlayerId(playerId: self.id)
                ?? panic("Account already with playerID: ".concat(self.id.toString())
                    .concat(" has not played FastBreak with ID: ".concat(fastBreakGameID)))

            emit FastBreakGameSubmissionUpdated(
                playerId: self.id,
                fastBreakGameID: fastBreakGameID,
                topShots: updatedSubmission.topShots,
            )
        }
    }

    pub struct PlayerData {

        pub let id: UInt64
        pub let playerName: String

        init(playerName: String) {
            self.id = FastBreakV1.nextPlayerId
            self.playerName = playerName
        }
    }

    /// Get a player id by account address
    ///
    pub fun getPlayerIdByAccount(accountAddress: Address): UInt64 {
        return FastBreakV1.accountPlayerMapping[accountAddress]!
    }

    /// Validate Fast Break Submission topShots
    ///
    pub fun validatePlaySubmission(fastBreakGame: FastBreakGame, topShots: [UInt64]): Bool {

        if (topShots.length < 1) {
            return false
        }

        if topShots.length > Int(fastBreakGame.numPlayers) {
            return false
        }

        return true
    }


    /// The Fast Break game token
    ///
    pub resource NFT: NonFungibleToken.INFT {
        pub let id: UInt64
        pub let fastBreakGameID: String /// The uuid of the Fast Break Game
        pub let serialNumber: UInt64 /// Each account mints game tokens from 1 => n
        pub let mintingDate: UInt64 /// The block timestamp of the tokens minting
        pub let mintedTo: UInt64 /// The playerId of the minter.
        pub let topShots: [UInt64] /// The top shot ids of the game tokens submission

        destroy() {
            emit FastBreakNFTBurned(id: self.id, serialNumber: self.serialNumber)
        }

        init(
            fastBreakGameID: String,
            serialNumber: UInt64,
            topShots: [UInt64],
            mintedTo: UInt64,
        ) {
            pre {
                FastBreakV1.fastBreakGameByID[fastBreakGameID] != nil: "No such fast break with gameId: ".concat(fastBreakGameID)
            }

            self.id = self.uuid
            self.fastBreakGameID = fastBreakGameID
            self.serialNumber = serialNumber
            self.mintingDate = UInt64(getCurrentBlock().timestamp)
            self.topShots = topShots
            self.mintedTo = mintedTo
        }

        pub fun isWinner(): Bool {
            if let fastBreak = FastBreakV1.fastBreakGameByID[self.fastBreakGameID] {
                if let submission = fastBreak.submissions[self.mintedTo] {
                    return submission.win
                }
            }
            return false
        }

        pub fun points(): UInt64 {
            if let fastBreak = FastBreakV1.fastBreakGameByID[self.fastBreakGameID] {
                if let submission = fastBreak.submissions[self.mintedTo] {
                    return submission.points
                }
            }
            return 0
        }
    }

    /// The Fast Break game token collection
    ///
    pub resource interface FastBreakNFTCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun batchDeposit(tokens: @NonFungibleToken.Collection)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowNFTSafe(id: UInt64): &NonFungibleToken.NFT?
        pub fun borrowFastBreakNFT(id: UInt64): &FastBreakV1.NFT? {
            post {
                (result == nil) || (result?.id == id):
                    "Cannot borrow Fast Break NFT reference: The ID of the returned reference is incorrect"
            }
        }
    }

    /// Capabilities of Fast Break Players
    ///
    pub resource interface FastBreakPlayer {
        pub fun play(
            fastBreakGameID: String,
            topShots: [UInt64]
        ): @FastBreakV1.NFT
    }

    /// Fast Break game collection
    ///
    pub resource Collection:
        NonFungibleToken.Provider,
        NonFungibleToken.Receiver,
        NonFungibleToken.CollectionPublic,
        FastBreakNFTCollectionPublic
    {

        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}

        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) 
                ?? panic("Could not find a fast break with the given ID in the Fast Break collection. Fast break Id: ".concat(withdrawID.toString()))

            emit Withdraw(id: token.id, from: self.owner?.address)

            return <-token
        }

        pub fun deposit(token: @NonFungibleToken.NFT) {
            let token <- token as! @FastBreakV1.NFT
            let id: UInt64 = token.id

            let oldToken <- self.ownedNFTs[id] <- token

            emit Deposit(id: id, to: self.owner?.address)

            destroy oldToken
        }

        pub fun batchDeposit(tokens: @NonFungibleToken.Collection) {
            let keys = tokens.getIDs()

            for key in keys {
                self.deposit(token: <-tokens.withdraw(withdrawID: key))
            }

            destroy tokens
        }

        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            return (&self.ownedNFTs[id] as &NonFungibleToken.NFT?)!
        }

        pub fun borrowNFTSafe(id: UInt64): &NonFungibleToken.NFT? {
            return (&self.ownedNFTs[id] as &NonFungibleToken.NFT?)
        }

        pub fun borrowFastBreakNFT(id: UInt64): &FastBreakV1.NFT? {
            if self.ownedNFTs[id] != nil {
                let ref = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)!
                return ref as! &FastBreakV1.NFT
            } else {
                return nil
            }
        }

        destroy() {
            destroy self.ownedNFTs
        }

        init() {
            self.ownedNFTs <- {}
        }
    }

    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create Collection()
    }

    pub fun createPlayer(playerName: String): @FastBreakV1.Player {
        FastBreakV1.nextPlayerId = FastBreakV1.nextPlayerId + UInt64(1)

        emit FastBreakPlayerCreated(
            id: FastBreakV1.nextPlayerId,
            playerName: playerName,
        )

        return <- create FastBreakV1.Player(playerName: playerName)
    }

    /// Capabilities of the Game Oracle
    ///
    pub resource interface GameOracle {
        pub fun createFastBreakRun(id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool)
        pub fun updateFastBreakRunStatus(id: String, status: UInt8)
        pub fun createFastBreakGame(
            id: String,
            name: String,
            fastBreakRunID: String,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        )
        pub fun updateFastBreakGame(id: String, status: UInt8, winner: UInt64)
        pub fun updateFastBreakScore(fastBreakGameID: String, playerId: UInt64, points: UInt64, win: Bool)
        pub fun addStatToFastBreakGame(fastBreakGameID: String, name: String, rawType: UInt8, valueNeeded: UInt64)
    }

    /// Fast Break Daemon game oracle implementation
    ///
    pub resource FastBreakDaemon: GameOracle {

        /// Create a Fast Break Run
        ///
        pub fun createFastBreakRun(id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool) {
            let fastBreakRun = FastBreakV1.FastBreakRun(
                id: id,
                name: name,
                runStart: runStart,
                runEnd: runEnd,
                fatigueModeOn: fatigueModeOn
            )
            FastBreakV1.fastBreakRunByID[fastBreakRun.id] = fastBreakRun
            emit FastBreakRunCreated(
                id: fastBreakRun.id,
                name: fastBreakRun.name,
                runStart: fastBreakRun.runStart,
                runEnd: fastBreakRun.runEnd,
                fatigueModeOn: fastBreakRun.fatigueModeOn
            )
        }

        /// Update the status of a Fast Break Run
        ///
        pub fun updateFastBreakRunStatus(id: String, status: UInt8) {
            let fastBreakRun = (&FastBreakV1.fastBreakRunByID[id] as &FastBreakV1.FastBreakRun?)
                ?? panic("Fast break run does not exist with Id: ".concat(id))

            let runStatus: FastBreakV1.RunStatus = FastBreakV1.RunStatus(rawValue: status)
                ?? panic("Run status does not exist with rawValue: ".concat(status.toString()))

            fastBreakRun.updateStatus(status: runStatus)

            emit FastBreakRunStatusChange(id: fastBreakRun.id, newRawStatus: fastBreakRun.status.rawValue)
        }

        /// Create a game of Fast Break
        ///
        pub fun createFastBreakGame(
            id: String,
            name: String,
            fastBreakRunID: String,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        ) {
            let fastBreakGame: FastBreakV1.FastBreakGame = FastBreakV1.FastBreakGame(
                id: id,
                name: name,
                fastBreakRunID: fastBreakRunID,
                submissionDeadline: submissionDeadline,
                numPlayers: numPlayers
            )
            FastBreakV1.fastBreakGameByID[fastBreakGame.id] = fastBreakGame
            emit FastBreakGameCreated(
                id: fastBreakGame.id,
                name: fastBreakGame.name,
                fastBreakRunID: fastBreakGame.fastBreakRunID,
                submissionDeadline: fastBreakGame.submissionDeadline,
                numPlayers: fastBreakGame.numPlayers
            )
        }

        /// Add a Fast Break Statistic to a game of Fast Break during game creation
        ///
        pub fun addStatToFastBreakGame(fastBreakGameID: String, name: String, rawType: UInt8, valueNeeded: UInt64) {

            let fastBreakGame: &FastBreakV1.FastBreakGame = (&FastBreakV1.fastBreakGameByID[fastBreakGameID] as &FastBreakV1.FastBreakGame?)
                ?? panic("Fast break does not exist with Id: ".concat(fastBreakGameID))

            let statType: FastBreakV1.StatisticType = FastBreakV1.StatisticType(rawValue: rawType)
                ?? panic("Fast break stat type does not exist with rawType: ".concat(rawType.toString()))

            let fastBreakStat : FastBreakV1.FastBreakStat = FastBreakV1.FastBreakStat(
                name: name,
                type: statType,
                valueNeeded: valueNeeded
            )

            fastBreakGame.addStat(stat: fastBreakStat)
            emit FastBreakGameStatAdded(
                fastBreakGameID: fastBreakGame.id,
                name: fastBreakStat.name,
                type: fastBreakStat.type.rawValue,
                valueNeeded: fastBreakStat.valueNeeded
            )

        }

        /// Update the status of a Fast Break
        ///
        pub fun updateFastBreakGame(id: String, status: UInt8, winner: UInt64) {

            let fastBreakGame: &FastBreakV1.FastBreakGame = (&FastBreakV1.fastBreakGameByID[id] as &FastBreakV1.FastBreakGame?)
                ?? panic("Fast break does not exist with Id: ".concat(id))

            let fastBreakStatus: FastBreakV1.GameStatus = FastBreakV1.GameStatus(rawValue: status)
                ?? panic("Fast break status does not exist with rawValue: ".concat(status.toString()))

            fastBreakGame.update(status: fastBreakStatus, winner: winner)

            emit FastBreakGameStatusChange(id: fastBreakGame.id, newRawStatus: fastBreakGame.status.rawValue)

        }

        /// Updates the submission scores of a Fast Break
        ///
        pub fun updateFastBreakScore(fastBreakGameID: String, playerId: UInt64, points: UInt64, win: Bool) {
            let fastBreakGame: &FastBreakV1.FastBreakGame = (&FastBreakV1.fastBreakGameByID[fastBreakGameID] as &FastBreakV1.FastBreakGame?)
                ?? panic("Fast break does not exist with Id: ".concat(fastBreakGameID))

            let isNewWin = fastBreakGame.updateScore(playerId: playerId, points: points, win: win)

            if isNewWin {
                let fastBreakRun = (&FastBreakV1.fastBreakRunByID[fastBreakGame.fastBreakRunID] as &FastBreakV1.FastBreakRun?)
                    ?? panic("Could not obtain reference to fast break run with Id: ".concat(fastBreakGame.fastBreakRunID))

                fastBreakRun.incrementRunWinCount(playerId: playerId)

                let submission = fastBreakGame.submissions[playerId]!

                emit FastBreakGameWinner(
                    playerId: playerId,
                    submittedAt: submission.submittedAt,
                    fastBreakGameID: submission.fastBreakGameID,
                    topShots: submission.topShots
                )

            }
        }
    }

    init() {
        self.CollectionStoragePath = /storage/FastBreakGameV1
        self.CollectionPublicPath = /public/FastBreakGameV1
        self.OracleStoragePath = /storage/FastBreakOracleV1
        self.OraclePrivatePath = /private/FastBreakOracleV1
        self.PlayerStoragePath = /storage/FastBreakPlayerV1
        self.PlayerPrivatePath = /private/FastBreakPlayerV1

        self.totalSupply = 0
        self.nextPlayerId = 0
        self.fastBreakRunByID = {}
        self.fastBreakGameByID = {}
        self.fastBreakPlayerByID = {}
        self.playerAccountMapping = {}
        self.accountPlayerMapping = {}

        let oracle <- create FastBreakDaemon()
        self.account.save(<-oracle, to: self.OracleStoragePath)
        self.account.link<&FastBreakV1.FastBreakDaemon{FastBreakV1.GameOracle}>(
            self.OraclePrivatePath,
            target: self.OracleStoragePath
        )

        emit ContractInitialized()
    }
}
