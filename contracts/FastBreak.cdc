/*
    Fast Break Game Contract
    Author: Jeremy Ahrens jer.ahrens@dapperlabs.com
*/

import NonFungibleToken from 0xf8d6e0586b0a20c7

pub contract FastBreak: NonFungibleToken {

    pub event ContractInitialized()
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)
    pub event FastBreakRunCreated(id: UInt64, name: String)
    pub event FastBreakRunStatusChange(id: UInt64, newStatus: String)
    pub event FastBreakGameCreated(id: UInt64, name: String)
    pub event FastBreakGameStatusChange(id: UInt64, newStatus: String)
    pub event FastBreakNFTBurned(id: UInt64, serialNumber: UInt64)
    pub event FastBreakNFTMinted(
        id: UInt64,
        fastBreakGameID: UInt64,
        serialNumber: UInt64
    )


    pub let CollectionStoragePath:  StoragePath
    pub let CollectionPublicPath:   PublicPath
    pub let OracleStoragePath:       StoragePath
    pub let OraclePrivatePath:      PrivatePath

    pub var totalSupply:        UInt64
    pub var nextFastBreakRunID:       UInt64
    pub var nextFastBreakGameID:          UInt64

    access(self) let fastBreakRunByID:        {UInt64: FastBreakRun}
    access(self) let fastBreakGameByID:           {UInt64: FastBreakGame}

    pub struct FastBreakRun {
        pub let id: UInt64
        pub let name: String
        pub var status: String
        pub let runStart: UInt64
        pub let runEnd: UInt64

        init (id: UInt64, name: String, runStart: UInt64, runEnd: UInt64) {
            if let fastBreakRun = FastBreak.fastBreakRunByID[id] {

                self.id = fastBreakRun.id
                self.name = fastBreakRun.name
                self.status = fastBreakRun.status
                self.runStart = fastBreakRun.runStart
                self.runEnd = fastBreakRun.runEnd
            } else {
                self.id = id
                self.name = name
                self.status = "FAST_BREAK_RUN_OPEN"
                self.runStart = runStart
                self.runEnd = runEnd
            }
        }

        pub fun updateStatus(status: String) { self.status = status }
    }

    pub fun getFastBreakRun(id: UInt64): FastBreak.FastBreakRun? {
        return FastBreak.fastBreakRunByID[id]
    }

    pub struct FastBreakGame {
        pub let id: UInt64
        pub let name: String
        pub let fatigueModeOn: Bool
        pub let isPublic: Bool
        pub let submissionDeadline: UInt64
        pub let numPlayers: UInt64
        pub var status: String
        pub var winner: Address?
        pub let submissions: {Address: AnyStruct}?

        init (
            id: UInt64,
            name: String,
            fastBreakRunId: UInt64,
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
            }
        }

        pub fun updateStatus(status: String) { self.status = status }

        pub fun updateWinner(winner: Address) { self.winner = winner }
    }

    pub fun getFastBreakGame(id: UInt64): FastBreak.FastBreakGame? {
        return FastBreak.fastBreakGameByID[id]
    }


    pub resource NFT: NonFungibleToken.INFT {
        pub let id: UInt64
        pub let fastBreakGameID: UInt64
        pub let serialNumber: UInt64
        pub let mintingDate: UFix64
        //pub let mintedTo: Address
        pub let topShots: [UInt64]

        destroy() {
            emit FastBreakNFTBurned(id: self.id, serialNumber: self.serialNumber)
        }

        init(
            fastBreakGameID: UInt64,
            serialNumber: UInt64,
            topShots: [UInt64]
        ) {
            pre {
                FastBreak.fastBreakGameByID[fastBreakGameID] != nil: "no such fast break"
            }

            self.id = self.uuid
            self.fastBreakGameID = fastBreakGameID
            self.serialNumber = serialNumber
            self.mintingDate = getCurrentBlock().timestamp
            self.topShots = topShots

            emit FastBreakNFTMinted(
                id: self.id,
                fastBreakGameID: self.fastBreakGameID,
                serialNumber: self.serialNumber
            )
        }


        pub fun isWinner() {
            //TODO return fastbreak.submissions[address].win
        }

        pub fun points() {
            //TODO return fastbreak.submissions[address].points
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

    pub resource Collection:
        NonFungibleToken.Provider,
        NonFungibleToken.Receiver,
        NonFungibleToken.CollectionPublic,
        FastBreakNFTCollectionPublic
    {

        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}

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

    pub resource interface FastBreakGamePlayer {
        pub fun play(fastBreakGameID: UInt64, topShots: [UInt64]): @FastBreak.NFT
    }

    pub resource interface GameOracle {
        pub fun createFastBreakRun(name: String, runStart: UInt64, runEnd: UInt64): UInt64
        pub fun updateFastBreakRunStatus(id: UInt64, status: String): UInt64
        pub fun createFastBreakGame(
            name: String,
            fastBreakRunId: UInt64,
            fatigueModeOn: Bool,
            isPublic: Bool,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        ): UInt64
        pub fun updateFastBreakGameStatus(id: UInt64, status: String): UInt64
    }

    pub resource FastBreakPlayer: FastBreakGamePlayer {

        pub fun play(fastBreakGameID: UInt64, topShots: [UInt64]): @FastBreak.NFT {
            pre {
                FastBreak.fastBreakGameByID.containsKey(fastBreakGameID): "no such fast break game"
            }

            let fastBreakGame = (&FastBreak.fastBreakGameByID[fastBreakGameID] as &FastBreak.FastBreakGame?)!

            let fastBreakNFT <- create NFT(
                fastBreakGameID: fastBreakGame.id,
                serialNumber: 0, // TODO not sure what I want to do here yet
                topShots: topShots
            )

            FastBreak.totalSupply = FastBreak.totalSupply + 1
            return <- fastBreakNFT
        }
    }

    pub resource FastBreakDaemon: GameOracle {

        pub fun createFastBreakRun(name: String, runStart: UInt64, runEnd: UInt64): UInt64 {

            let fastBreakRun = FastBreak.FastBreakRun(
                id: FastBreak.nextFastBreakRunID,
                name: name,
                runStart: runStart,
                runEnd: runEnd
            )
            FastBreak.fastBreakRunByID[fastBreakRun.id] = fastBreakRun
            emit FastBreakRunCreated(
                id: fastBreakRun.id,
                name: fastBreakRun.name
            )
            FastBreak.nextFastBreakRunID = fastBreakRun.id + 1 as UInt64
            return fastBreakRun.id
        }

        pub fun updateFastBreakRunStatus(id: UInt64, status: String): UInt64 {
            let fastBreakRun = (&FastBreak.fastBreakRunByID[id] as &FastBreak.FastBreakRun?)!

            fastBreakRun.updateStatus(status: status)
            emit FastBreakRunStatusChange(id: fastBreakRun.id, newStatus: fastBreakRun.status)
            return fastBreakRun.id
        }

        pub fun createFastBreakGame(
            name: String,
            fastBreakRunId: UInt64,
            fatigueModeOn: Bool,
            isPublic: Bool,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        ): UInt64 {

            let fastBreakGame: FastBreak.FastBreakGame = FastBreak.FastBreakGame(
                id: FastBreak.nextFastBreakGameID,
                name: name,
                fastBreakRunId: fastBreakRunId,
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
            FastBreak.nextFastBreakGameID = fastBreakGame.id + 1 as UInt64
            return fastBreakGame.id
        }

        pub fun updateFastBreakGameStatus(id: UInt64, status: String): UInt64 {
            let fastBreakGame: &FastBreak.FastBreakGame = (&FastBreak.fastBreakGameByID[id] as &FastBreak.FastBreakGame?)!

            fastBreakGame.updateStatus(status: status)
            emit FastBreakGameStatusChange(id: fastBreakGame.id, newStatus: fastBreakGame.status)
            return fastBreakGame.id
        }
    }

    init() {
        self.CollectionStoragePath = /storage/FastBreakNFTCollection
        self.CollectionPublicPath = /public/FastBreakNFTCollection
        self.OracleStoragePath = /storage/FastBreakDaemon
        self.OraclePrivatePath = /private/FastBreakDaemon

        self.totalSupply = 0
        self.nextFastBreakRunID = 1
        self.nextFastBreakGameID = 1

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