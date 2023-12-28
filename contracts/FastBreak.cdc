/*
    Fast Break Game Contract
    Author: Jeremy Ahrens jer.ahrens@dapperlabs.com
*/

import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS

/// Game & Oracle Contract for Fast Break
///
pub contract FastBreak: NonFungibleToken {

    /// Contract events
    ///
    pub event ContractInitialized()
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)
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
        isPublic: Bool,
        submissionDeadline: UInt64,
        numPlayers: UInt64
    )
    pub event FastBreakGameStatusChange(id: String, newRawStatus: UInt8)
    pub event FastBreakNFTBurned(id: UInt64, serialNumber: UInt64)
    pub event FastBreakNFTMinted(
        id: UInt64,
        fastBreakGameID: String,
        serialNumber: UInt64,
        mintingDate: UInt64,
        topShots: [UInt64],
        mintedTo: Address
    )
    pub event FastBreakSubmissionSent(
        wallet: Address,
        submittedAt: UInt64,
        fastBreakGameID: String,
        topShots: [UInt64]
    )
    pub event FastBreakGameWinner(
        wallet: Address?,
        submittedAt: UInt64,
        fastBreakGameID: String,
        topShots: [UInt64]
    )
    pub event FastBreakGameStatAdded(
        fastBreakGameID: String,
        name: String,
        type: String,
        valueNeeded: UInt64
    )

    /// Named Paths
    ///
    pub let CollectionStoragePath:  StoragePath
    pub let CollectionPublicPath:   PublicPath
    pub let OracleStoragePath:       StoragePath
    pub let OraclePrivatePath:      PrivatePath

    /// Contract variables
    ///
    pub var totalSupply:        UInt64

    /// Game Enums
    ///
    pub enum GameStatus: UInt8 {
        pub case SCHEDULED
        pub case OPEN
        pub case STARTED
        pub case CLOSED
    }

    pub enum RunStatus: UInt8 {
        pub case SCHEDULED
        pub case RUNNING
        pub case CLOSED
    }

    /// Metadata Dictionaries
    ///
    access(self) let fastBreakRunByID:        {String: FastBreakRun}
    access(self) let fastBreakGameByID:           {String: FastBreakGame}

    /// A top-level Fast Break Run, the container for Fast Break Games
    ///
    pub struct FastBreakRun {
        pub let id: String
        pub let name: String
        pub var status: FastBreak.RunStatus
        pub let runStart: UInt64
        pub let runEnd: UInt64
        pub let leaderboard: {Address: UInt64}
        pub let fatigueModeOn: Bool

        init (id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool) {
            if let fastBreakRun = FastBreak.fastBreakRunByID[id] {
                self.id = fastBreakRun.id
                self.name = fastBreakRun.name
                self.status = fastBreakRun.status
                self.runStart = fastBreakRun.runStart
                self.runEnd = fastBreakRun.runEnd
                self.leaderboard = fastBreakRun.leaderboard
                self.fatigueModeOn = fastBreakRun.fatigueModeOn
            } else {
                self.id = id
                self.name = name
                self.status = FastBreak.RunStatus.SCHEDULED
                self.runStart = runStart
                self.runEnd = runEnd
                self.leaderboard = {}
                self.fatigueModeOn = fatigueModeOn
            }
        }

        /// Update status of the Fast Break Run
        ///
        access(contract) fun updateStatus(status: FastBreak.RunStatus) { self.status = status }

        /// Write a new win to the Fast Break Run leaderboard
        ///
        access(contract) fun incrementLeaderboardWins(wallet: Address) {
            let leaderboard = self.leaderboard
            var wins = leaderboard[wallet] ?? 0
            wins = wins + 1
            leaderboard[wallet] = wins
        }
    }

    /// Get a Fast Break Run by Id
    ///
    pub fun getFastBreakRun(id: String): FastBreak.FastBreakRun? {
        return FastBreak.fastBreakRunByID[id]
    }

    /// A single Game of Fast Break
    ///
    pub struct FastBreakGame {
        pub let id: String
        pub let name: String
        pub let isPublic: Bool
        pub let submissionDeadline: UInt64
        pub let numPlayers: UInt64
        pub var status: FastBreak.GameStatus
        pub var winner: Address?
        pub var submissions: {Address: FastBreak.FastBreakSubmission}
        pub let fastBreakRunID: String
        pub var stats: [FastBreakStat]

        init (
            id: String,
            name: String,
            fastBreakRunID: String,
            isPublic: Bool,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        ) {
            if let fb = FastBreak.fastBreakGameByID[id] {
                self.id = fb.id
                self.name = fb.name
                self.isPublic = fb.isPublic
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
                self.isPublic = isPublic
                self.submissionDeadline = submissionDeadline
                self.numPlayers = numPlayers
                self.status = FastBreak.GameStatus.SCHEDULED
                self.submissions = {}
                self.fastBreakRunID = fastBreakRunID
                self.stats = []
                self.winner = nil
            }
        }

        /// Get a wallet's active Fast Break Submission
        ///
        pub fun getFastBreakSubmissionByWallet(wallet: Address): FastBreak.FastBreakSubmission? {
            let fastBreakSubmissions = self.submissions!

            return fastBreakSubmissions[wallet]
        }

        /// Add a statistic to the Fast Break during game creation
        ///
        access(contract) fun addStat(stat: FastBreak.FastBreakStat) {
            self.stats.append(stat)
        }

        /// Update status and winner of a Fast Break
        ///
        access(contract) fun update(status: FastBreak.GameStatus, winner: Address?) {
            self.status = status
            self.winner = winner
        }

        /// Submit a Fast Break
        ///
        access(contract) fun submitFastBreak(submission: FastBreak.FastBreakSubmission) {
            pre {
                FastBreak.isValidSubmission(submissionDeadline: self.submissionDeadline) : "submission missed deadline"
            }

            self.submissions[submission.wallet] = submission
        }

        /// Update the Fast Break score of a wallet
        ///
        access(contract) fun updateScore(wallet: Address, points: UInt64, win: Bool): Bool {
            var isNewWin = false
            let submissions = self.submissions
            let submission: FastBreak.FastBreakSubmission = submissions[wallet]!
            if win && !submission.win {
                isNewWin = true
            }
            submission.setPoints(points: points, win: win)
            self.submissions[wallet] = submission
            return isNewWin
        }
    }

    /// Validate Fast Break Submission
    ///
    pub fun isValidSubmission(submissionDeadline: UInt64): Bool {
        if submissionDeadline > UInt64(getCurrentBlock().timestamp) {
            return true
        }

        return false
    }

    /// Get a Fast Break Game by Id
    ///
    pub fun getFastBreakGame(id: String): FastBreak.FastBreakGame? {
        return FastBreak.fastBreakGameByID[id]
    }

    /// Get the game stats of a Fast Break
    ///
    pub fun getFastBreakGameStats(id: String): [FastBreak.FastBreakStat] {
        let fastBreak = FastBreak.getFastBreakGame(id: id)!
        return fastBreak.stats
    }

    /// A statistical structure used in Fast Break Games
    ///
    pub struct FastBreakStat {
        pub let name: String
        pub let type: String
        pub let valueNeeded: UInt64

        init (
            name: String,
            type: String,
            valueNeeded: UInt64
        ) {
            self.name = name
            self.type = type
            self.valueNeeded = valueNeeded
        }
    }

    /// A wallet submission to a Fast Break
    ///
    pub struct FastBreakSubmission {
        pub let wallet: Address
        pub var submittedAt: UInt64
        pub let fastBreakGameID: String
        pub var topShots: [UInt64]
        pub var points: UInt64
        pub var win: Bool

        init (
            wallet: Address,
            fastBreakGameID: String,
            topShots: [UInt64],
        ) {
            self.wallet = wallet
            self.fastBreakGameID = fastBreakGameID
            self.topShots = topShots
            self.submittedAt = UInt64(getCurrentBlock().timestamp)
            self.points = 0
            self.win = false

            emit FastBreakSubmissionSent(
                wallet: self.wallet,
                submittedAt: self.submittedAt,
                fastBreakGameID: self.fastBreakGameID,
                topShots: self.topShots
            )
        }

        /// Set the points of a submission
        ///
        access(contract) fun setPoints(points: UInt64, win: Bool) {
            self.points = points
            self.win = win
        }

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
        pub let fastBreakGameID: String
        pub let serialNumber: UInt64
        pub let mintingDate: UInt64
        pub let mintedTo: Address
        pub let topShots: [UInt64]

        destroy() {
            emit FastBreakNFTBurned(id: self.id, serialNumber: self.serialNumber)
        }

        init(
            fastBreakGameID: String,
            serialNumber: UInt64,
            topShots: [UInt64],
            mintedTo: Address,
        ) {
            pre {
                FastBreak.fastBreakGameByID[fastBreakGameID] != nil: "no such fast break"
            }

            self.id = self.uuid
            self.fastBreakGameID = fastBreakGameID
            self.serialNumber = serialNumber
            self.mintingDate = UInt64(getCurrentBlock().timestamp)
            self.topShots = topShots
            self.mintedTo = mintedTo

            emit FastBreakNFTMinted(
                id: self.id,
                fastBreakGameID: self.fastBreakGameID,
                serialNumber: self.serialNumber,
                mintingDate: self.mintingDate,
                topShots: self.topShots,
                mintedTo: self.mintedTo
            )
        }

        pub fun isWinner(): Bool {
            let fastBreak : FastBreak.FastBreakGame = FastBreak.fastBreakGameByID[self.fastBreakGameID]!
            let submission : FastBreak.FastBreakSubmission = fastBreak.submissions[self.mintedTo]!
            return submission.win
        }

        pub fun points(): UInt64 {
            let fastBreak : FastBreak.FastBreakGame = FastBreak.fastBreakGameByID[self.fastBreakGameID]!
            let submission : FastBreak.FastBreakSubmission = fastBreak.submissions[self.mintedTo]!
            return submission.points
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
        pub fun borrowFastBreakNFT(id: UInt64): &FastBreak.NFT? {
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
        ): @FastBreak.NFT
    }

    /// Fast Break game collection
    ///
    pub resource Collection:
        NonFungibleToken.Provider,
        NonFungibleToken.Receiver,
        NonFungibleToken.CollectionPublic,
        FastBreakNFTCollectionPublic,
        FastBreakPlayer
    {

        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}
        pub var numMinted: UInt64

        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("Could not find a fast break with the given ID in the Fast Break collection")

            emit Withdraw(id: token.id, from: self.owner?.address)

            return <-token
        }

        pub fun deposit(token: @NonFungibleToken.NFT) {
            let token <- token as! @FastBreak.NFT
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

        pub fun borrowFastBreakNFT(id: UInt64): &FastBreak.NFT? {
            if self.ownedNFTs[id] != nil {
                let ref = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)!
                return ref as! &FastBreak.NFT
            } else {
                return nil
            }
        }

        /// Play the game of Fast Break with an array of Top Shots
        ///
        pub fun play(
            fastBreakGameID: String,
            topShots: [UInt64]
        ): @FastBreak.NFT {
            pre {
                FastBreak.fastBreakGameByID.containsKey(fastBreakGameID): "no such fast break game"
                FastBreak.fastBreakGameByID[fastBreakGameID]!.isPublic: "fast break game is private"
            }

            let acct = getAccount(self.owner?.address!)
            let collectionRef = acct.getCapability(/public/MomentCollection)
                                    .borrow<&{TopShot.MomentCollectionPublic}>()!

            /// Must own Top Shots to play Fast Break
            for flowId in topShots {
                if !collectionRef.getIDs().contains(flowId) {
                    panic("top shot not owned in collection")
                }
            }

            let fastBreakGame = (&FastBreak.fastBreakGameByID[fastBreakGameID] as &FastBreak.FastBreakGame?)!

            /// Cannot mint two tokens for the same Fast Break
            let existingSubmission = fastBreakGame.getFastBreakSubmissionByWallet(wallet: self.owner?.address!)
            if existingSubmission != nil {
                panic("wallet already submitted to fast break")
            }

            let fastBreakSubmission = FastBreak.FastBreakSubmission(
                wallet: self.owner?.address!,
                fastBreakGameID: fastBreakGame.id,
                topShots: topShots
            )

            fastBreakGame.submitFastBreak(submission: fastBreakSubmission)

            let fastBreakNFT <- create NFT(
                fastBreakGameID: fastBreakGame.id,
                serialNumber: self.numMinted + 1,
                topShots: topShots,
                mintedTo: self.owner?.address!
            )

            FastBreak.totalSupply = FastBreak.totalSupply + 1
            return <- fastBreakNFT
        }

        destroy() {
            destroy self.ownedNFTs
        }

        init() {
            self.numMinted = 0
            self.ownedNFTs <- {}
        }
    }

    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create Collection()
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
            isPublic: Bool,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        )
        pub fun updateFastBreakGame(id: String, status: UInt8, winner: Address?)
        pub fun submitFastBreak(wallet: Address, submission: FastBreak.FastBreakSubmission)
        pub fun updateFastBreakScore(fastBreakGameID: String, wallet: Address, points: UInt64, win: Bool)
        pub fun addStatToFastBreakGame(fastBreakGameID: String, name: String, type: String, valueNeeded: UInt64)
    }

    /// Fast Break Daemon game oracle implementation
    ///
    pub resource FastBreakDaemon: GameOracle {

        /// Create a Fast Break Run
        ///
        pub fun createFastBreakRun(id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool) {
            let fastBreakRun = FastBreak.FastBreakRun(
                id: id,
                name: name,
                runStart: runStart,
                runEnd: runEnd,
                fatigueModeOn: fatigueModeOn
            )
            FastBreak.fastBreakRunByID[fastBreakRun.id] = fastBreakRun
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
            let fastBreakRun = (&FastBreak.fastBreakRunByID[id] as &FastBreak.FastBreakRun?)!
            let runStatus: FastBreak.RunStatus = FastBreak.RunStatus(rawValue: status)!

            fastBreakRun.updateStatus(status: runStatus)
            emit FastBreakRunStatusChange(id: fastBreakRun.id, newRawStatus: fastBreakRun.status.rawValue)
        }

        /// Create a game of Fast Break
        ///
        pub fun createFastBreakGame(
            id: String,
            name: String,
            fastBreakRunID: String,
            isPublic: Bool,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        ) {
            let fastBreakGame: FastBreak.FastBreakGame = FastBreak.FastBreakGame(
                id: id,
                name: name,
                fastBreakRunID: fastBreakRunID,
                isPublic: isPublic,
                submissionDeadline: submissionDeadline,
                numPlayers: numPlayers
            )
            FastBreak.fastBreakGameByID[fastBreakGame.id] = fastBreakGame
            emit FastBreakGameCreated(
                id: fastBreakGame.id,
                name: fastBreakGame.name,
                fastBreakRunID: fastBreakGame.fastBreakRunID,
                isPublic: fastBreakGame.isPublic,
                submissionDeadline: fastBreakGame.submissionDeadline,
                numPlayers: fastBreakGame.numPlayers
            )
        }

        /// Add a Fast Break Statistic to a game of Fast Break during game creation
        ///
        pub fun addStatToFastBreakGame(fastBreakGameID: String, name: String, type: String, valueNeeded: UInt64) {
            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[fastBreakGameID] as &FastBreak.FastBreakGame?)!
            let fastBreakStat : FastBreak.FastBreakStat = FastBreak.FastBreakStat(
                name: name,
                type: type,
                valueNeeded: valueNeeded
            )

            fastBreakGame.addStat(stat: fastBreakStat)
            emit FastBreakGameStatAdded(
                fastBreakGameID: fastBreakGame.id,
                name: fastBreakStat.name,
                type: fastBreakStat.type,
                valueNeeded: fastBreakStat.valueNeeded
            )
        }

        /// Update the status of a Fast Break
        ///
        pub fun updateFastBreakGame(id: String, status: UInt8, winner: Address?) {
            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[id] as &FastBreak.FastBreakGame?)!
            let fastBreakStatus: FastBreak.GameStatus = FastBreak.GameStatus(rawValue: status)!

            fastBreakGame.update(status: fastBreakStatus, winner: winner)
            emit FastBreakGameStatusChange(id: fastBreakGame.id, newRawStatus: fastBreakGame.status.rawValue)
        }

        /// Submit a Fast Break on behalf of a wallet
        ///
        pub fun submitFastBreak(wallet: Address, submission: FastBreak.FastBreakSubmission) {
            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[submission.fastBreakGameID] as &FastBreak.FastBreakGame?)!
            fastBreakGame.submitFastBreak(submission: submission)
        }

        /// Updates the submission scores of a Fast Break
        ///
        pub fun updateFastBreakScore(fastBreakGameID: String, wallet: Address, points: UInt64, win: Bool) {
            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[fastBreakGameID] as &FastBreak.FastBreakGame?)!
            let isNewWin = fastBreakGame.updateScore(wallet: wallet, points: points, win: win)
            if isNewWin {

                let fastBreakRun: &FastBreak.FastBreakRun =
                    (&FastBreak.fastBreakRunByID[fastBreakGame.fastBreakRunID] as &FastBreak.FastBreakRun?)!

                fastBreakRun.incrementLeaderboardWins(wallet: wallet)

                let submission = fastBreakGame.submissions[wallet]!

                emit FastBreakGameWinner(
                    wallet: submission.wallet,
                    submittedAt: submission.submittedAt,
                    fastBreakGameID: submission.fastBreakGameID,
                    topShots: submission.topShots
                )
            }
        }
    }

    init() {
        self.CollectionStoragePath = /storage/FastBreakGame
        self.CollectionPublicPath = /public/FastBreakGame
        self.OracleStoragePath = /storage/FastBreakDaemon
        self.OraclePrivatePath = /private/FastBreakDaemon

        self.totalSupply = 0
        self.fastBreakRunByID = {}
        self.fastBreakGameByID = {}

        let oracle <- create FastBreakDaemon()
        self.account.save(<-oracle, to: self.OracleStoragePath)
        self.account.link<&FastBreak.FastBreakDaemon{FastBreak.GameOracle}>(
            self.OraclePrivatePath,
            target: self.OracleStoragePath
        )

        emit ContractInitialized()
    }
}
