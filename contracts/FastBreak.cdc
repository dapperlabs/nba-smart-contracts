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
        mintedTo: Address
    )

    pub event FastBreakGameWinner(
        accountAddress: Address?,
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

    /// Contract variables
    ///
    pub var totalSupply:        UInt64

    /// Game Enums
    ///

    /// A game of Fast Break has the following status transitions
    ///
    pub enum GameStatus: UInt8 {
        pub case SCHEDULED
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

    /// A top-level Fast Break Run, the container for Fast Break Games
    /// A Fast Break Run contains many Fast Break games & is a mini-season.
    /// Fatigue mode applies submission limitations for the off-chain version of the game
    /// Fatigue mode limits top shot usage by tier. 4 uses legendary. 2 uses rare. 1 use other.
    ///
    pub struct FastBreakRun {
        pub let id: String /// The off-chain uuid of the Fast Break Run
        pub let name: String /// The name of the Run (R0, R1, etc)
        pub var status: FastBreak.RunStatus /// The status of the run
        pub let runStart: UInt64 /// The block timestamp starting the run
        pub let runEnd: UInt64 /// The block timestamp ending the run
        pub let runWinCount: {Address: UInt64} /// win count by address
        pub let fatigueModeOn: Bool /// Fatigue mode is a game rule limiting usage of top shots by tier

        init (id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool) {
            if let fastBreakRun = FastBreak.fastBreakRunByID[id] {
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
                self.status = FastBreak.RunStatus.SCHEDULED
                self.runStart = runStart
                self.runEnd = runEnd
                self.runWinCount = {}
                self.fatigueModeOn = fatigueModeOn
            }
        }

        /// Update status of the Fast Break Run
        ///
        access(contract) fun updateStatus(status: FastBreak.RunStatus) { self.status = status }

        /// Write a new win to the Fast Break Run runWinCount
        ///
        access(contract) fun incrementRunWinCount(accountAddress: Address) {
            let runWinCount = self.runWinCount
            runWinCount[accountAddress] = (runWinCount[accountAddress] ?? 0) + 1
        }
    }

    /// Get a Fast Break Run by Id
    ///
    pub fun getFastBreakRun(id: String): FastBreak.FastBreakRun? {
        return FastBreak.fastBreakRunByID[id]
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
        pub var status: FastBreak.GameStatus /// The game status
        pub var winner: Address? /// The address of the winner of Fast Break
        pub var submissions: {Address: FastBreak.FastBreakSubmission} /// Map of each submission to the Fast break
        pub let fastBreakRunID: String /// The off-chain uuid of the Fast Break Run containing this Fast Break
        pub var stats: [FastBreakStat] /// The NBA statistical requirements for this Fast Break

        init (
            id: String,
            name: String,
            fastBreakRunID: String,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        ) {
            if let fb = FastBreak.fastBreakGameByID[id] {
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
                self.status = FastBreak.GameStatus.SCHEDULED
                self.submissions = {}
                self.fastBreakRunID = fastBreakRunID
                self.stats = []
                self.winner = nil
            }
        }

        /// Get a account's active Fast Break Submission
        ///
        pub fun getFastBreakSubmissionByAccount(accountAddress: Address): FastBreak.FastBreakSubmission? {
            let fastBreakSubmissions = self.submissions

            return fastBreakSubmissions[accountAddress]
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

            self.submissions[submission.accountAddress] = submission
        }

        /// Update the Fast Break score of an account
        ///
        access(contract) fun updateScore(accountAddress: Address, points: UInt64, win: Bool): Bool {
            let submissions = self.submissions

            let submission: FastBreak.FastBreakSubmission = submissions[accountAddress]
                ?? panic("unable to find fast break submission for account address")

            submission.setPoints(points: points, win: win)

            self.submissions[accountAddress] = submission

            if win && !submission.win {
                return true
            }


            return false
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
        if let fastBreak = FastBreak.getFastBreakGame(id: id) {
            return fastBreak.stats
        }
        return []
    }

    /// A statistical structure used in Fast Break Games
    /// This structure names the NBA statistic top shots must match or exceed
    /// An example is points as the statistic and 30 as the value
    /// A top shot or group of top shots must meet or exceed 30 points
    ///
    pub struct FastBreakStat {
        pub let name: String
        pub let type: FastBreak.StatisticType
        pub let valueNeeded: UInt64

        init (
            name: String,
            type: FastBreak.StatisticType,
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
        pub let accountAddress: Address
        pub var submittedAt: UInt64
        pub let fastBreakGameID: String
        pub var topShots: [UInt64]
        pub var points: UInt64
        pub var win: Bool

        init (
            accountAddress: Address,
            fastBreakGameID: String,
            topShots: [UInt64],
        ) {
            self.accountAddress = accountAddress
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
        pub let mintedTo: Address /// The address of the minter. Used for composability.
        pub let topShots: [UInt64] /// The top shot ids of the game tokens submission

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

            emit FastBreakGameTokenMinted(
                id: self.id,
                fastBreakGameID: self.fastBreakGameID,
                serialNumber: self.serialNumber,
                mintingDate: self.mintingDate,
                topShots: self.topShots,
                mintedTo: self.mintedTo
            )
        }

        pub fun isWinner(): Bool {
            if let fastBreak = FastBreak.fastBreakGameByID[self.fastBreakGameID] {
                if let submission = fastBreak.submissions[self.mintedTo] {
                    return submission.win
                }
            }
            return false
        }

        pub fun points(): UInt64 {
            if let fastBreak = FastBreak.fastBreakGameByID[self.fastBreakGameID] {
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
        /// Each account must own a top shot collection to play fast break
        ///
        pub fun play(
            fastBreakGameID: String,
            topShots: [UInt64]
        ): @FastBreak.NFT {
            pre {
                FastBreak.fastBreakGameByID.containsKey(fastBreakGameID): "no such fast break game"
            }

            /// Validate Top Shots
            let acct = getAccount(self.owner?.address!)
            let collectionRef = acct.getCapability(/public/MomentCollection)
                .borrow<&{TopShot.MomentCollectionPublic}>() ?? panic("player does not have top shot collection")

            /// Must own Top Shots to play Fast Break
            for flowId in topShots {
                if !collectionRef.getIDs().contains(flowId) {
                    panic("top shot not owned in collection")
                }
            }

            let fastBreakGame = (&FastBreak.fastBreakGameByID[fastBreakGameID] as &FastBreak.FastBreakGame?)
                 ?? panic("fast break does not exist")

            /// Cannot mint two tokens for the same Fast Break
            let existingSubmission = fastBreakGame.getFastBreakSubmissionByAccount(accountAddress: self.owner?.address!)
            if existingSubmission != nil {
                panic("account already submitted to fast break")
            }

            let fastBreakSubmission = FastBreak.FastBreakSubmission(
                accountAddress: self.owner?.address!,
                fastBreakGameID: fastBreakGameID,
                topShots: topShots
            )

            fastBreakGame.submitFastBreak(submission: fastBreakSubmission)


            let fastBreakNFT <- create NFT(
                fastBreakGameID: fastBreakGameID,
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
            submissionDeadline: UInt64,
            numPlayers: UInt64
        )
        pub fun updateFastBreakGame(id: String, status: UInt8, winner: Address?)
        pub fun updateFastBreakScore(fastBreakGameID: String, accountAddress: Address, points: UInt64, win: Bool)
        pub fun addStatToFastBreakGame(fastBreakGameID: String, name: String, rawType: UInt8, valueNeeded: UInt64)
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
            let fastBreakRun = (&FastBreak.fastBreakRunByID[id] as &FastBreak.FastBreakRun?)
                ?? panic("fast break run does not exist")

            let runStatus: FastBreak.RunStatus = FastBreak.RunStatus(rawValue: status)
                ?? panic("run status does not exist")

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
            let fastBreakGame: FastBreak.FastBreakGame = FastBreak.FastBreakGame(
                id: id,
                name: name,
                fastBreakRunID: fastBreakRunID,
                submissionDeadline: submissionDeadline,
                numPlayers: numPlayers
            )
            FastBreak.fastBreakGameByID[fastBreakGame.id] = fastBreakGame
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

            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[fastBreakGameID] as &FastBreak.FastBreakGame?)
                ?? panic("fast break does not exist")

            let statType: FastBreak.StatisticType = FastBreak.StatisticType(rawValue: rawType)
                ?? panic("fast break stat type does not exist")

            let fastBreakStat : FastBreak.FastBreakStat = FastBreak.FastBreakStat(
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
        pub fun updateFastBreakGame(id: String, status: UInt8, winner: Address?) {

            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[id] as &FastBreak.FastBreakGame?)
                ?? panic("fast break does not exist")

            let fastBreakStatus: FastBreak.GameStatus = FastBreak.GameStatus(rawValue: status)
                ?? panic("fast break status does not exist")

            fastBreakGame.update(status: fastBreakStatus, winner: winner)

            emit FastBreakGameStatusChange(id: fastBreakGame.id, newRawStatus: fastBreakGame.status.rawValue)

        }

        /// Updates the submission scores of a Fast Break
        ///
        pub fun updateFastBreakScore(fastBreakGameID: String, accountAddress: Address, points: UInt64, win: Bool) {
            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[fastBreakGameID] as &FastBreak.FastBreakGame?)
                ?? panic("fast break does not exist")

            let isNewWin = fastBreakGame.updateScore(accountAddress: accountAddress, points: points, win: win)

            if isNewWin {
                let fastBreakRun = (&FastBreak.fastBreakRunByID[fastBreakGame.fastBreakRunID] as &FastBreak.FastBreakRun?)
                    ?? panic("could not obtain reference to fast break run")

                fastBreakRun.incrementRunWinCount(accountAddress: accountAddress)

                let submission = fastBreakGame.submissions[accountAddress]!

                emit FastBreakGameWinner(
                    accountAddress: submission.accountAddress,
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
