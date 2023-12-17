/*
    Fast Break Game Contract
    Author: Jeremy Ahrens jer.ahrens@dapperlabs.com
*/

import NonFungibleToken from 0xNFTADDRESS
//import TopShot from 0xTOPSHOTADDRESS
//import NonFungibleToken from 0xf8d6e0586b0a20c7

pub contract FastBreak: NonFungibleToken {

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
    pub event FastBreakRunStatusChange(id: String, newStatus: String)
    pub event FastBreakGameCreated(id: String, name: String)
    pub event FastBreakGameStatusChange(id: String, newStatus: String)
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

    pub let CollectionStoragePath:  StoragePath
    pub let CollectionPublicPath:   PublicPath
    pub let OracleStoragePath:       StoragePath
    pub let OraclePrivatePath:      PrivatePath

    pub var totalSupply:        UInt64
    access(self) let fastBreakRunByID:        {String: FastBreakRun}
    access(self) let fastBreakGameByID:           {String: FastBreakGame}

    pub struct FastBreakRun {
        pub let id: String
        pub let name: String
        pub var status: String
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
                self.status = "FAST_BREAK_RUN_OPEN"
                self.runStart = runStart
                self.runEnd = runEnd
                self.leaderboard = {}
                self.fatigueModeOn = fatigueModeOn
            }
        }

        access(contract) fun updateStatus(status: String) { self.status = status }

        access(contract) fun updateLeaderboard(wallet: Address, wins: UInt64) {
            let leaderboard = self.leaderboard
            leaderboard[wallet] = wins
        }
    }

    pub fun getFastBreakRun(id: String): FastBreak.FastBreakRun? {
        return FastBreak.fastBreakRunByID[id]
    }

    pub struct FastBreakGame {
        pub let id: String
        pub let name: String
        pub let fatigueModeOn: Bool
        pub let isPublic: Bool
        pub let submissionDeadline: UInt64
        pub let numPlayers: UInt64
        pub var status: String
        pub var winner: Address?
        pub let submissions: {Address: FastBreak.FastBreakSubmission}
        pub let fastBreakRunID: String

        init (
            id: String,
            name: String,
            fastBreakRunID: String,
            fatigueModeOn: Bool,
            isPublic: Bool,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        ) {
            if let fb = FastBreak.fastBreakGameByID[id] {
                self.id = fb.id
                self.name = fb.name
                self.fatigueModeOn = fb.fatigueModeOn
                self.isPublic = fb.isPublic
                self.submissionDeadline = fb.submissionDeadline
                self.numPlayers = fb.numPlayers
                self.status = fb.status
                self.winner = fb.winner
                self.submissions = fb.submissions
                self.fastBreakRunID = fb.fastBreakRunID
            } else {
                self.id = id
                self.name = name
                self.fatigueModeOn = fatigueModeOn
                self.isPublic = isPublic
                self.submissionDeadline = submissionDeadline
                self.numPlayers = numPlayers
                self.status = "FAST_BREAK_OPEN"
                self.winner = 0x0000000000000000
                self.submissions = {}
                self.fastBreakRunID = fastBreakRunID
            }
        }

        pub fun getFastBreakSubmissionByWallet(wallet: Address): FastBreak.FastBreakSubmission? {
            let fastBreakSubmissions = self.submissions!

            return fastBreakSubmissions[wallet]
        }

        access(contract) fun update(status: String, winner: Address) {
            self.status = status
            self.winner = winner
        }

        access(contract) fun submitFastBreak(submission: FastBreak.FastBreakSubmission) {
            pre {
                self.submissionDeadline > UInt64(getCurrentBlock().timestamp): "submission missed deadline"
            }

            self.submissions[submission.wallet] = submission
        }

        access(contract) fun updateScore(wallet: Address, points: UInt64, win: Bool) {
            let submissions = self.submissions
            let submission: FastBreak.FastBreakSubmission = submissions[wallet]!
            submission.setPoints(points: points, win: win)
            self.submissions[wallet] = submission
        }
    }

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

        access(contract) fun setPoints(points: UInt64, win: Bool) {
            self.points = points
            self.win = win
        }

    }

    pub fun getFastBreakGame(id: String): FastBreak.FastBreakGame? {
        return FastBreak.fastBreakGameByID[id]
    }

    pub fun validatePlaySubmission(fastBreakGame: FastBreakGame, topShots: [UInt64]): Bool {
        if Int(fastBreakGame.numPlayers) == topShots.length {
            return true
        }
        return false
    }

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

    pub resource interface FastBreakPlayer {
        pub fun play(fastBreakGameID: String, topShots: [UInt64]): @FastBreak.NFT
    }

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

        pub fun play(fastBreakGameID: String, topShots: [UInt64]): @FastBreak.NFT {

            pre {
                FastBreak.fastBreakGameByID.containsKey(fastBreakGameID): "no such fast break game"
            }

            let fastBreakGame = (&FastBreak.fastBreakGameByID[fastBreakGameID] as &FastBreak.FastBreakGame?)!

            let fastBreakNFT <- create NFT(
                fastBreakGameID: fastBreakGame.id,
                serialNumber: self.numMinted + 1,
                topShots: topShots,
                mintedTo: self.owner?.address!
            )

            let fastBreakSubmission = FastBreak.FastBreakSubmission(
                wallet: self.owner?.address!,
                fastBreakGameID: fastBreakNFT.fastBreakGameID,
                topShots: fastBreakNFT.topShots
            )

            fastBreakGame.submitFastBreak(submission: fastBreakSubmission)

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

    pub resource interface GameOracle {
        pub fun createFastBreakRun(id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool)
        pub fun updateFastBreakRunStatus(id: String, status: String)
        pub fun createFastBreakGame(
            id: String,
            name: String,
            fastBreakRunID: String,
            fatigueModeOn: Bool,
            isPublic: Bool,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        )
        pub fun updateFastBreakGame(id: String, status: String, winner: Address)
        pub fun submitFastBreak(wallet: Address, submission: FastBreak.FastBreakSubmission)
        pub fun updateFastBreakScore(fastBreakGameID: String, wallet: Address, points: UInt64, win: Bool)
    }

    pub resource FastBreakDaemon: GameOracle {

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

        pub fun updateFastBreakRunStatus(id: String, status: String) {
            let fastBreakRun = (&FastBreak.fastBreakRunByID[id] as &FastBreak.FastBreakRun?)!

            fastBreakRun.updateStatus(status: status)
            emit FastBreakRunStatusChange(id: fastBreakRun.id, newStatus: fastBreakRun.status)
        }

        pub fun createFastBreakGame(
            id: String,
            name: String,
            fastBreakRunID: String,
            fatigueModeOn: Bool,
            isPublic: Bool,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        ) {
            let fastBreakGame: FastBreak.FastBreakGame = FastBreak.FastBreakGame(
                id: id,
                name: name,
                fastBreakRunID: fastBreakRunID,
                fatigueModeOn: fatigueModeOn,
                isPublic: isPublic,
                submissionDeadline: submissionDeadline,
                numPlayers: numPlayers
            )
            FastBreak.fastBreakGameByID[fastBreakGame.id] = fastBreakGame
            emit FastBreakGameCreated(
                id: fastBreakGame.id,
                name: fastBreakGame.name
            )
        }

        pub fun updateFastBreakGame(id: String, status: String, winner: Address) {
            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[id] as &FastBreak.FastBreakGame?)!

            fastBreakGame.update(status: status, winner: winner)
            emit FastBreakGameStatusChange(id: fastBreakGame.id, newStatus: fastBreakGame.status)
        }

        pub fun submitFastBreak(wallet: Address, submission: FastBreak.FastBreakSubmission) {
            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[submission.fastBreakGameID] as &FastBreak.FastBreakGame?)!
            fastBreakGame.submitFastBreak(submission: submission)
        }

        pub fun updateFastBreakScore(fastBreakGameID: String, wallet: Address, points: UInt64, win: Bool) {
            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[fastBreakGameID] as &FastBreak.FastBreakGame?)!
            fastBreakGame.updateScore(wallet: wallet, points: points, win: win)
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
