/*
      _____                 __    ___.                         __
    _/ ____\____    _______/  |_  \_ |_________   ____ _____  |  | __
    \   __\\__  \  /  ___/\   __\  | __ \_  __ \_/ __ \\__  \ |  |/ /
     |  |   / __ \_\___ \  |  |    | \_\ \  | \/\  ___/ / __ \|    <
     |__|  (____  /____  > |__|    |___  /__|    \___  >____  /__|_ \
                \/     \/              \/            \/     \/     \/

    fast break game contract & oracle

*/

import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS
import TopShotMarketV3, Market from 0xMARKETV3ADDRESS

/// Game & Oracle Contract for Fast Break V1
///
access(all) contract FastBreakV1: NonFungibleToken {

    access(all) entitlement Play
    access(all) entitlement Create
    access(all) entitlement Update

    /// Contract events
    ///

    access(all) event FastBreakPlayerCreated(
        id: UInt64,
        playerName: String
    )

    access(all) event FastBreakRunCreated(
        id: String,
        name: String,
        runStart: UInt64,
        runEnd: UInt64,
        fatigueModeOn: Bool
    )

    access(all) event FastBreakRunStatusChange(id: String, newRawStatus: UInt8)

    access(all) event FastBreakGameCreated(
        id: String,
        name: String,
        fastBreakRunID: String,
        submissionDeadline: UInt64,
        numPlayers: UInt64
    )

    access(all) event FastBreakGameStatusChange(id: String, newRawStatus: UInt8)

    access(all) event FastBreakNFTBurned(id: UInt64, serialNumber: UInt64)

    access(all) event FastBreakGameTokenMinted(
        id: UInt64,
        fastBreakGameID: String,
        serialNumber: UInt64,
        mintingDate: UInt64,
        topShots: [UInt64],
        mintedTo: UInt64
    )

    access(all) event FastBreakGameSubmissionUpdated(
        playerId: UInt64,
        fastBreakGameID: String,
        topShots: [UInt64],
    )

    access(all) event FastBreakGameWinner(
        playerId: UInt64,
        submittedAt: UInt64,
        fastBreakGameID: String,
        topShots: &[UInt64]
    )

    access(all) event FastBreakGameStatAdded(
        fastBreakGameID: String,
        name: String,
        type: UInt8,
        valueNeeded: UInt64
    )

    /// Named Paths
    ///
    access(all) let CollectionStoragePath:      StoragePath
    access(all) let CollectionPublicPath:       PublicPath
    access(all) let OracleStoragePath:          StoragePath
    access(all) let PlayerStoragePath:          StoragePath

    /// Contract variables
    ///
    access(all) var totalSupply:        UInt64
    access(all) var nextPlayerId:        UInt64

    /// Game Enums
    ///

    /// A game of Fast Break has the following status transitions
    ///
    access(all) enum GameStatus: UInt8 {
        access(all) case SCHEDULED /// Game is schedules but closed for submission
        access(all) case OPEN /// Game is open for submission
        access(all) case STARTED /// Game has started
        access(all) case CLOSED /// Game is over and rewards are being distributed
    }

    /// A Fast Break Run has the following status transitions
    ///
    access(all) enum RunStatus: UInt8 {
        access(all) case SCHEDULED
        access(all) case RUNNING /// The first Fast Break game of the run has started
        access(all) case CLOSED /// The last Fast Break game of the run has ended
    }

    /// A Fast Break Statistic can be met by an individual or group of top shots
    ///
    access(all) enum StatisticType: UInt8 {
        access(all) case INDIVIDUAL /// Each top shot must meet or exceed this statistical value
        access(all) case CUMMULATIVE /// All top shots in the submission must meet or exceed this statistical value
    }

    /// Metadata Dictionaries
    /// Player mappings remain in contract storage (bounded by number of users)
    /// Old dictionaries kept for existing data (pre-upgrade)
    ///
    access(self) let fastBreakRunByID:      {String: FastBreakRun}  // Legacy: existing data
    access(self) let fastBreakGameByID:     {String: FastBreakGame}  // Legacy: existing data
    access(self) let fastBreakPlayerByID:   {UInt64: PlayerData}
    access(self) let playerAccountMapping:  {UInt64: Address}
    access(self) let accountPlayerMapping:  {Address: UInt64}

    /// Storage resources for year-based organization
    /// Each year's games/runs stored in separate resources in account storage
    ///
    access(all) resource YearGameStorage {
        access(self) var games: {String: FastBreakGame}
        
        init() {
            self.games = {}
        }
    }

    access(all) resource YearRunStorage {
        access(self) var runs: {String: FastBreakRun}
        
        init() {
            self.runs = {}
        }
    }

    /// Helper function to extract year number from timestamp
    /// Returns year as UInt64 (e.g., 2024)
    ///
    access(self) fun getYearFromTimestamp(timestamp: UInt64): UInt64 {
        // Convert timestamp to date and extract year
        // Unix timestamp: seconds since Jan 1, 1970
        // Approximate: 31536000 seconds per year
        let secondsPerYear: UInt64 = 31536000
        let year1970: UInt64 = 1970
        return year1970 + (timestamp / secondsPerYear)
    }

    /// Generate storage path for game storage based on year
    /// Creates path like "FastBreakGames2024" for year 2024
    ///
    access(self) fun generateGameStoragePathByYear(year: String): StoragePath {
        return StoragePath(identifier: "FastBreakGames".concat(year))!
    }

    /// Generate storage path for run storage based on year
    /// Creates path like "FastBreakRuns2024" for year 2024
    ///
    access(self) fun generateRunStoragePathByYear(year: String): StoragePath {
        return StoragePath(identifier: "FastBreakRuns".concat(year))!
    }

    /// Get or create year's game storage resource from year
    /// If storage doesn't exist for the year, creates and saves it
    /// Returns a reference to the storage resource
    ///
    access(self) fun getOrCreateYearGameStorage(year: String): &YearGameStorage {
        let path = FastBreakV1.generateGameStoragePathByYear(year: year)
        if let storage = self.account.storage.borrow<&YearGameStorage>(from: path) {
            return storage
        }
        // Create and save new storage for this year
        let storage <- create YearGameStorage()
        self.account.storage.save(<-storage, to: path)
        return self.account.storage.borrow<&YearGameStorage>(from: path)!
    }

    /// Get or create year's run storage resource from year
    /// If storage doesn't exist for the year, creates and saves it
    /// Returns a reference to the storage resource
    ///
    access(self) fun getOrCreateYearRunStorage(year: String): &YearRunStorage {
        let path = FastBreakV1.generateRunStoragePathByYear(year: year)
        if let storage = self.account.storage.borrow<&YearRunStorage>(from: path) {
            return storage
        }
        // Create and save new storage for this year
        let storage <- create YearRunStorage()
        self.account.storage.save(<-storage, to: path)
        return self.account.storage.borrow<&YearRunStorage>(from: path)!
    }

    /// Get year's game storage from year string (returns nil if doesn't exist)
    ///
    access(self) fun getYearGameStorageByYear(year: String): &YearGameStorage? {
        let path = FastBreakV1.generateGameStoragePathByYear(year: year)
        return self.account.storage.borrow<&YearGameStorage>(from: path)
    }

    /// Get year's run storage from year string (returns nil if doesn't exist)
    ///
    access(self) fun getYearRunStorageByYear(year: String): &YearRunStorage? {
        let path = FastBreakV1.generateRunStoragePathByYear(year: year)
        return self.account.storage.borrow<&YearRunStorage>(from: path)
    }

    /// Helper to update a game in storage after mutation
    /// Checks if it's in legacy dict first, otherwise updates year-based storage
    /// Requires year parameter to determine correct storage location
    ///
    access(self) fun updateGameInStorage(game: &FastBreakGame, gameID: String, year: String) {
        // Check if it's in legacy dictionary
        if FastBreakV1.fastBreakGameByID[gameID] != nil {
            FastBreakV1.fastBreakGameByID[gameID] = *game
        } else {
            // Update in year-based storage
            if let yearStorage = FastBreakV1.getYearGameStorageByYear(year: year) {
                yearStorage.games[gameID] = *game
            }
        }
    }

    /// Helper to update a run in storage after mutation
    /// Checks if it's in legacy dict first, otherwise updates year-based storage
    /// Requires year parameter to determine correct storage location
    ///
    access(self) fun updateRunInStorage(run: &FastBreakRun, runID: String, year: String) {
        // Check if it's in legacy dictionary
        if FastBreakV1.fastBreakRunByID[runID] != nil {
            FastBreakV1.fastBreakRunByID[runID] = *run
        } else {
            // Update in year-based storage
            if let yearStorage = FastBreakV1.getYearRunStorageByYear(year: year) {
                yearStorage.runs[runID] = *run
            }
        }
    }

    /// A top-level Fast Break Run, the container for Fast Break Games
    /// A Fast Break Run contains many Fast Break games & is a mini-season.
    /// Fatigue mode applies submission limitations for the off-chain version of the game
    /// Fatigue mode limits top shot usage by tier. 4 uses legendary. 2 uses rare. 1 use other.
    ///
    access(all) struct FastBreakRun {
        access(all) let id: String /// The off-chain uuid of the Fast Break Run
        access(all) let name: String /// The name of the Run (R0, R1, etc)
        access(all) var status: FastBreakV1.RunStatus /// The status of the run
        access(all) let runStart: UInt64 /// The block timestamp starting the run
        access(all) let runEnd: UInt64 /// The block timestamp ending the run
        access(all) let runWinCount: {UInt64: UInt64} /// win count by playerId
        access(all) let fatigueModeOn: Bool /// Fatigue mode is a game rule limiting usage of top shots by tier

        init (id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool) {
            // Check if run already exists in account storage
            let year = FastBreakV1.getYearFromTimestamp(timestamp: runStart).toString()
            if let yearStorage = FastBreakV1.getYearRunStorageByYear(year: year) {
                if let existingRun = yearStorage.runs[id] {
                    self.id = existingRun.id
                    self.name = existingRun.name
                    self.status = existingRun.status
                    self.runStart = existingRun.runStart
                    self.runEnd = existingRun.runEnd
                    self.runWinCount = existingRun.runWinCount
                    self.fatigueModeOn = existingRun.fatigueModeOn
                    return
                }
            }
            // Create new run
            self.id = id
            self.name = name
            self.status = FastBreakV1.RunStatus.SCHEDULED
            self.runStart = runStart
            self.runEnd = runEnd
            self.runWinCount = {}
            self.fatigueModeOn = fatigueModeOn
        }

        /// Update status of the Fast Break Run
        ///
        access(contract) fun updateStatus(status: FastBreakV1.RunStatus) { self.status = status }

        /// Write a new win to the Fast Break Run runWinCount
        ///
        access(contract) fun incrementRunWinCount(playerId: UInt64) {
            self.runWinCount[playerId] = (self.runWinCount[playerId] ?? 0) + 1
        }
    }

    /// Get a Fast Break Run by Id
    /// ONLY checks legacy dictionary (for existing data before upgrade)
    /// For new data, use getFastBreakRunByYear with the year parameter
    ///
    access(all) view fun getFastBreakRun(id: String): FastBreakV1.FastBreakRun? {
        // Only check legacy dictionary (existing data before upgrade)
        return FastBreakV1.fastBreakRunByID[id]
    }

    /// Get a Fast Break Run by Id and Year (direct access - O(1))
    /// Checks year-based storage first, then falls back to legacy dictionary
    ///
    access(all) view fun getFastBreakRunByYear(id: String, year: String): FastBreakV1.FastBreakRun? {
        // First check year-based storage using year (direct access)
        if let yearStorage = FastBreakV1.getYearRunStorageByYear(year: year) {
            if let run = yearStorage.runs[id] {
                return run
            }
        }
        
        // Fallback to legacy dictionary (existing data before upgrade)
        return FastBreakV1.fastBreakRunByID[id]
    }

    /// A single Game of Fast Break
    /// A Fast Break is played on any day NBA games are scheduled
    /// It is the intention of this contract to allow private & public Fast Break games
    /// A private Fast Break is visible on-chain but is restricted to private accounts
    /// A public Fast Break can be played by custodial and non-custodial accounts
    ///
    access(all) struct FastBreakGame {
        access(all) let id: String /// The off-chain uuid of the Fast Break
        access(all) let name: String /// The name of the Fast Break (eg FB0, FB1, FB2)
        access(all) var submissionDeadline: UInt64 /// The block timestamp restricting submission to the Fast Break
        access(all) let numPlayers: UInt64 /// The number of top shots a player should submit to the Fast Break
        access(all) var status: FastBreakV1.GameStatus /// The game status
        access(all) var winner: UInt64 /// The playerId of the winner of Fast Break
        access(all) var submissions: {UInt64: FastBreakV1.FastBreakSubmission} /// Map of player submission to the Fast Break
        access(all) let fastBreakRunID: String /// The off-chain uuid of the Fast Break Run containing this Fast Break
        access(all) var stats: [FastBreakStat] /// The NBA statistical requirements for this Fast Break

        init (
            id: String,
            name: String,
            fastBreakRunID: String,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        ) {
            // Check if game already exists in account storage
            let year = FastBreakV1.getYearFromTimestamp(timestamp: submissionDeadline).toString()
            if let yearStorage = FastBreakV1.getYearGameStorageByYear(year: year) {
                if let existingGame = yearStorage.games[id] {
                    self.id = existingGame.id
                    self.name = existingGame.name
                    self.submissionDeadline = existingGame.submissionDeadline
                    self.numPlayers = existingGame.numPlayers
                    self.status = existingGame.status
                    self.winner = existingGame.winner
                    self.submissions = existingGame.submissions
                    self.fastBreakRunID = existingGame.fastBreakRunID
                    self.stats = existingGame.stats
                    return
                }
            }
            // Create new game
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

        /// Get a account's active Fast Break Submission
        ///
        access(all) view fun getFastBreakSubmissionByPlayerId(playerId: UInt64): FastBreakV1.FastBreakSubmission? {
            return self.submissions[playerId]
        }

        /// Add a statistic to the Fast Break during game creation
        ///
        access(contract) fun addStat(stat: FastBreakV1.FastBreakStat) {
            self.stats.append(stat)
        }

        /// Set the submission deadline for a Fast Break
        ///
        access(contract) fun setSubmissionDeadline(deadline: UInt64) {
            self.submissionDeadline = deadline
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
            let submission: FastBreakV1.FastBreakSubmission = self.submissions[playerId]
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
    access(all) view fun isValidSubmission(submissionDeadline: UInt64): Bool {
        return submissionDeadline > UInt64(getCurrentBlock().timestamp) + 60
    }

    /// Internal helper to get a game by ID
    /// Searches current year first, then previous year, then legacy (for oracle functions)
    /// Oracle functions operate in same year, so prioritize new year-based storage
    ///
    access(self) fun getFastBreakGameInternal(id: String): &FastBreakV1.FastBreakGame? {
        let currentTimestamp = UInt64(getCurrentBlock().timestamp)
        let currentYearNum = FastBreakV1.getYearFromTimestamp(timestamp: currentTimestamp)
        
        // Check current year first (most likely - oracle functions operate in same year)
        let currentYearString = currentYearNum.toString()
        if let yearStorage = FastBreakV1.getYearGameStorageByYear(year: currentYearString) {
            if let game = yearStorage.games[id] {
                return &game as &FastBreakV1.FastBreakGame
            }
        }
        
        // Then check previous year
        if currentYearNum > 1970 {
            let previousYearString = (currentYearNum - 1).toString()
            if let yearStorage = FastBreakV1.getYearGameStorageByYear(year: previousYearString) {
                if let game = yearStorage.games[id] {
                    return &game as &FastBreakV1.FastBreakGame
                }
            }
        }
        
        // Finally check legacy dictionary (existing data before upgrade)
        if let game = FastBreakV1.fastBreakGameByID[id] {
            return &game as &FastBreakV1.FastBreakGame
        }
        
        return nil
    }

    /// Get a Fast Break Game by Id
    /// ONLY checks legacy dictionary (for existing data before upgrade)
    /// For new data, use getFastBreakGameByYear with the year parameter
    ///
    access(all) view fun getFastBreakGame(id: String): FastBreakV1.FastBreakGame? {
        // Only check legacy dictionary (existing data before upgrade)
        return FastBreakV1.fastBreakGameByID[id]
    }

    /// Get a Fast Break Game by Id and Year (direct access - O(1))
    /// Checks year-based storage first, then falls back to legacy dictionary
    ///
    access(all) view fun getFastBreakGameByYear(id: String, year: String): FastBreakV1.FastBreakGame? {
        // First check year-based storage using year (direct access)
        if let yearStorage = FastBreakV1.getYearGameStorageByYear(year: year) {
            if let game = yearStorage.games[id] {
                return game
            }
        }
        
        // Fallback to legacy dictionary (existing data before upgrade)
        return FastBreakV1.fastBreakGameByID[id]
    }

    /// Internal helper to get a run by ID
    /// Searches current year first, then previous year, then legacy (for oracle functions)
    /// Oracle functions operate in same year, so prioritize new year-based storage
    ///
    access(self) fun getFastBreakRunInternal(id: String): &FastBreakV1.FastBreakRun? {
        let currentTimestamp = UInt64(getCurrentBlock().timestamp)
        let currentYearNum = FastBreakV1.getYearFromTimestamp(timestamp: currentTimestamp)
        
        // Check current year first (most likely - oracle functions operate in same year)
        let currentYearString = currentYearNum.toString()
        if let yearStorage = FastBreakV1.getYearRunStorageByYear(year: currentYearString) {
            if let run = yearStorage.runs[id] {
                return &run as &FastBreakV1.FastBreakRun
            }
        }
        
        // Then check previous year
        if currentYearNum > 1970 {
            let previousYearString = (currentYearNum - 1).toString()
            if let yearStorage = FastBreakV1.getYearRunStorageByYear(year: previousYearString) {
                if let run = yearStorage.runs[id] {
                    return &run as &FastBreakV1.FastBreakRun
                }
            }
        }
        
        // Finally check legacy dictionary (existing data before upgrade)
        if let run = FastBreakV1.fastBreakRunByID[id] {
            return &run as &FastBreakV1.FastBreakRun
        }
        
        return nil
    }

    /// Get the game stats of a Fast Break
    ///
    access(all) view fun getFastBreakGameStats(id: String): [FastBreakV1.FastBreakStat] {
        if let fastBreak = FastBreakV1.getFastBreakGame(id: id) {
            return fastBreak.stats
        }
        return []
    }

    /// Get a Fast Break account by playerId
    ///
    access(all) view fun getFastBreakPlayer(id: UInt64): Address? {
        return FastBreakV1.playerAccountMapping[id]
    }

    /// A statistical structure used in Fast Break Games
    /// This structure names the NBA statistic top shots must match or exceed
    /// An example is points as the statistic and 30 as the value
    /// A top shot or group of top shots must meet or exceed 30 points
    ///
    access(all) struct FastBreakStat {
        access(all) let name: String
        access(all) let type: FastBreakV1.StatisticType
        access(all) let valueNeeded: UInt64

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
    access(all) struct FastBreakSubmission {
        access(all) let playerId: UInt64
        access(all) var submittedAt: UInt64
        access(all) let fastBreakGameID: String
        access(all) var topShots: [UInt64]
        access(all) var points: UInt64
        access(all) var win: Bool

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
    access(all) resource Player: FastBreakPlayer, NonFungibleToken.NFT {

        access(all) let id: UInt64
        access(all) let playerName: String      /// username
        access(all) var tokensMinted: UInt64    /// num games played

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
        access(Play) fun play(
            fastBreakGameID: String,
            topShots: [UInt64]
        ): @FastBreakV1.NFT {
            pre {
                FastBreakV1.getFastBreakGameInternal(id: fastBreakGameID) != nil: "No such fast break game with gameId: ".concat(fastBreakGameID)
            }

            /// Update player address mapping
            if let ownerAddress = self.owner?.address {
                FastBreakV1.playerAccountMapping[self.id] = ownerAddress
                FastBreakV1.accountPlayerMapping[ownerAddress] = self.id
            }

            /// Validate Top Shots
            let acct = getAccount(self.owner?.address!)
            let collectionRef = acct.capabilities.borrow<&TopShot.Collection>(/public/MomentCollection)
                ?? panic("Player does not have top shot collection")
            let marketV3CollectionRef = acct.capabilities.borrow<&TopShotMarketV3.SaleCollection>(/public/topshotSalev3Collection)
            let marketV1CollectionRef = acct.capabilities.borrow<&Market.SaleCollection>(/public/topshotSaleCollection)

            /// Must own Top Shots to play Fast Break
            /// more efficient to borrow ref than to loop
            ///
            for flowId in topShots {
                let topShotRef = collectionRef.borrowMoment(id: flowId)
                if topShotRef == nil {
                    let hasMarketPlaceV3 = marketV3CollectionRef != nil && marketV3CollectionRef!.borrowMoment(id: flowId) != nil
                    let hasMarketV1 = marketV1CollectionRef != nil && marketV1CollectionRef!.borrowMoment(id: flowId) != nil
                    if !hasMarketPlaceV3 && !hasMarketV1{
                        panic("Top shot not owned in any collection with flowId: ".concat(flowId.toString()))
                    }
                }
            }

            let fastBreakGame = FastBreakV1.getFastBreakGameInternal(id: fastBreakGameID)
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
            
            // Update game in account storage after mutation
            let year = FastBreakV1.getYearFromTimestamp(timestamp: fastBreakGame.submissionDeadline).toString()
            FastBreakV1.updateGameInStorage(game: fastBreakGame, gameID: fastBreakGameID, year: year)

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
        access(Update) fun updateSubmission(
            fastBreakGameID: String,
            topShots: [UInt64]
        ) {
            pre {
                FastBreakV1.getFastBreakGameInternal(id: fastBreakGameID) != nil: "No such fast break game with gameId: ".concat(fastBreakGameID)
            }

            /// Update player address mapping
            if let ownerAddress = self.owner?.address {
                FastBreakV1.playerAccountMapping[self.id] = ownerAddress
                FastBreakV1.accountPlayerMapping[ownerAddress] = self.id
            }

            /// Validate Top Shots
            let acct = getAccount(self.owner?.address!)
            let collectionRef = acct.capabilities.borrow<&TopShot.Collection>(/public/MomentCollection)
                ?? panic("Player does not have top shot collection")
            let marketV3CollectionRef = acct.capabilities.borrow<&TopShotMarketV3.SaleCollection>(/public/topshotSalev3Collection)
            let marketV1CollectionRef = acct.capabilities.borrow<&Market.SaleCollection>(/public/topshotSaleCollection)

            /// Must own Top Shots to play Fast Break
            /// more efficient to borrow ref than to loop
            ///
            for flowId in topShots {
                let topShotRef = collectionRef.borrowMoment(id: flowId)
                if topShotRef == nil {
                    let hasMarketPlaceV3 = marketV3CollectionRef != nil && marketV3CollectionRef!.borrowMoment(id: flowId) != nil
                    let hasMarketV1 = marketV1CollectionRef != nil && marketV1CollectionRef!.borrowMoment(id: flowId) != nil
                    if !hasMarketPlaceV3 && !hasMarketV1{
                        panic("Top shot not owned in any collection with flowId: ".concat(flowId.toString()))
                    }
                }
            }

            let fastBreakGame = FastBreakV1.getFastBreakGameInternal(id: fastBreakGameID)
                ?? panic("Fast break does not exist with gameId: ".concat(fastBreakGameID))

            /// Check that the user has a submission for Fast Break game we can update
            let pastSubmission = fastBreakGame.getFastBreakSubmissionByPlayerId(playerId: self.id)
                ?? panic("Account already with playerID: ".concat(self.id.toString())
                    .concat(" has not played FastBreak with ID: ".concat(fastBreakGameID)))

            fastBreakGame.updateFastBreakTopshots(playerId: self.id, topshotMoments: topShots)
            
            // Update game in account storage after mutation
            let year = FastBreakV1.getYearFromTimestamp(timestamp: fastBreakGame.submissionDeadline).toString()
            FastBreakV1.updateGameInStorage(game: fastBreakGame, gameID: fastBreakGameID, year: year)

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

        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <- FastBreakV1.createEmptyCollection(nftType: Type<@FastBreakV1.Player>())
        }

        access(all) view fun getViews(): [Type] {
            return [
                Type<MetadataViews.NFTCollectionData>(),
                Type<MetadataViews.NFTCollectionDisplay>()
            ]
        }

        access(all) fun resolveView(_ view: Type): AnyStruct? {
            switch view {
                case Type<MetadataViews.NFTCollectionData>():
                    return FastBreakV1.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>())
                case Type<MetadataViews.NFTCollectionDisplay>():
                    return FastBreakV1.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionDisplay>())
            }
            return nil
        }
    }

    access(all) struct PlayerData {

        access(all) let id: UInt64
        access(all) let playerName: String

        init(playerName: String) {
            self.id = FastBreakV1.nextPlayerId
            self.playerName = playerName
        }
    }

    /// Get a player id by account address
    ///
    access(all) view fun getPlayerIdByAccount(accountAddress: Address): UInt64 {
        return FastBreakV1.accountPlayerMapping[accountAddress]!
    }

    /// Validate Fast Break Submission topShots
    ///
    access(all) view fun validatePlaySubmission(fastBreakGame: FastBreakGame, topShots: [UInt64]): Bool {

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
    access(all) resource NFT: NonFungibleToken.NFT {
        access(all) let id: UInt64
        access(all) let fastBreakGameID: String /// The uuid of the Fast Break Game
        access(all) let serialNumber: UInt64 /// Each account mints game tokens from 1 => n
        access(all) let mintingDate: UInt64 /// The block timestamp of the tokens minting
        access(all) let mintedTo: UInt64 /// The playerId of the minter.
        access(all) let topShots: [UInt64] /// The top shot ids of the game tokens submission

        access(all) event ResourceDestroyed(
            id: UInt64 = self.id,
            serialNumber:  UInt64 = self.serialNumber
        )

        init(
            fastBreakGameID: String,
            serialNumber: UInt64,
            topShots: [UInt64],
            mintedTo: UInt64,
        ) {
            pre {
                FastBreakV1.getFastBreakGameInternal(id: fastBreakGameID) != nil: "No such fast break with gameId: ".concat(fastBreakGameID)
            }

            self.id = self.uuid
            self.fastBreakGameID = fastBreakGameID
            self.serialNumber = serialNumber
            self.mintingDate = UInt64(getCurrentBlock().timestamp)
            self.topShots = topShots
            self.mintedTo = mintedTo
        }

        access(all) view fun isWinner(): Bool {
            // Check game submissions (searches current year → previous year → legacy)
            if let fastBreak = FastBreakV1.getFastBreakGameInternal(id: self.fastBreakGameID) {
                if let submission = fastBreak.submissions[self.mintedTo] {
                    return submission.win
                }
            }
            return false
        }

        access(all) view fun points(): UInt64 {
            // Check game submissions (searches current year → previous year → legacy)
            if let fastBreak = FastBreakV1.getFastBreakGameInternal(id: self.fastBreakGameID) {
                if let submission = fastBreak.submissions[self.mintedTo] {
                    return submission.points
                }
            }
            return 0
        }

        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <- FastBreakV1.createEmptyCollection(nftType: Type<@FastBreakV1.NFT>())
        }

        access(all) view fun getViews(): [Type] {
            return [
                Type<MetadataViews.NFTCollectionData>(),
                Type<MetadataViews.NFTCollectionDisplay>()
            ]
        }

        access(all) fun resolveView(_ view: Type): AnyStruct? {
            switch view {
                case Type<MetadataViews.NFTCollectionData>():
                    return FastBreakV1.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>())
                case Type<MetadataViews.NFTCollectionDisplay>():
                    return FastBreakV1.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionDisplay>())
            }
            return nil
        }
    }

    /// The Fast Break game token collection
    ///
    access(all) resource interface FastBreakNFTCollectionPublic : NonFungibleToken.CollectionPublic  {
        access(all) fun batchDeposit(tokens: @{NonFungibleToken.Collection})
        access(all) fun borrowFastBreakNFT(id: UInt64): &FastBreakV1.NFT? {
            post {
                (result == nil) || (result?.id == id):
                    "Cannot borrow Fast Break NFT reference: The ID of the returned reference is incorrect"
            }
        }
    }

    /// Capabilities of Fast Break Players
    ///
    access(all) resource interface FastBreakPlayer {
        access(Play) fun play(
            fastBreakGameID: String,
            topShots: [UInt64]
        ): @FastBreakV1.NFT
    }

    /// Fast Break game collection
    ///
    access(all) resource Collection:
        NonFungibleToken.Collection,
        FastBreakNFTCollectionPublic
    {

        access(all) var ownedNFTs: @{UInt64: {NonFungibleToken.NFT}}

        access(NonFungibleToken.Withdraw) fun withdraw(withdrawID: UInt64): @{NonFungibleToken.NFT} {
            let token <- self.ownedNFTs.remove(key: withdrawID) 
                ?? panic("Could not find a fast break with the given ID in the Fast Break collection. Fast break Id: ".concat(withdrawID.toString()))

            return <-token
        }

        access(all) fun deposit(token: @{NonFungibleToken.NFT}) {
            let token <- token as! @FastBreakV1.NFT
            let id: UInt64 = token.id

            let oldToken <- self.ownedNFTs[id] <- token

            destroy oldToken
        }

        access(all) fun batchDeposit(tokens: @{NonFungibleToken.Collection}) {
            let keys = tokens.getIDs()

            for key in keys {
                self.deposit(token: <-tokens.withdraw(withdrawID: key))
            }

            destroy tokens
        }

        access(all) view fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        access(all) view fun borrowNFT(_ id: UInt64): &{NonFungibleToken.NFT}? {
            return &self.ownedNFTs[id]
        }

        access(all) view fun borrowFastBreakNFT(id: UInt64): &FastBreakV1.NFT? {
            return self.borrowNFT(id) as! &FastBreakV1.NFT?
        }

        access(all) view fun getSupportedNFTTypes(): {Type: Bool} {
            let supportedTypes: {Type: Bool} = {}
            supportedTypes[Type<@FastBreakV1.NFT>()] = true
            return supportedTypes
        }

        // Return whether or not the given type is accepted by the collection
        // A collection that can accept any type should just return true by default
        access(all) view fun isSupportedNFTType(type: Type): Bool {
            if type == Type<@FastBreakV1.NFT>() {
                return true
            }
            return false
        }

        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <- FastBreakV1.createEmptyCollection(nftType: Type<@FastBreakV1.NFT>())
        }

        access(all) view fun getLength(): Int {
            return self.ownedNFTs.length
        }

        init() {
            self.ownedNFTs <- {}
        }
    }

    access(all) fun createEmptyCollection(nftType: Type): @{NonFungibleToken.Collection} {
        if nftType != Type<@FastBreakV1.NFT>() {
            panic("NFT type is not supported")
        }
        return <- create Collection()
    }

    access(all) view fun getContractViews(resourceType: Type?): [Type] {
        return [Type<MetadataViews.NFTCollectionData>(), Type<MetadataViews.NFTCollectionDisplay>()]
    }

    access(all) view fun resolveContractView(resourceType: Type?, viewType: Type): AnyStruct? {
        post {
            result == nil || result!.getType() == viewType: "The returned view must be of the given type or nil"
        }
        switch viewType {
            case Type<MetadataViews.NFTCollectionData>():
                return MetadataViews.NFTCollectionData(
                    storagePath: /storage/FastBreakGameV1,
                    publicPath: /public/FastBreakGameV1,
                    publicCollection: Type<&FastBreakV1.Collection>(),
                    publicLinkedType: Type<&FastBreakV1.Collection>(),
                    createEmptyCollectionFunction: (fun (): @{NonFungibleToken.Collection} {
                        return <-FastBreakV1.createEmptyCollection(nftType: Type<@FastBreakV1.NFT>())
                    })
                )
            case Type<MetadataViews.NFTCollectionDisplay>():
                let bannerImage = MetadataViews.Media(
                    file: MetadataViews.HTTPFile(
                        url: "https://nbatopshot.com/static/fastbreak/fast-break-logo.svg"
                    ),
                    mediaType: "image/svg+xml"
                )
                let squareImage = MetadataViews.Media(
                    file: MetadataViews.HTTPFile(
                        url: "https://nbatopshot.com/static/fastbreak/fast-break-logo.svg"
                    ),
                    mediaType: "image/png"
                )
                return MetadataViews.NFTCollectionDisplay(
                    name: "NBA Top Shot Fast Break",
                    description: "The game of Fast Break is very simple. Collectors will select five players every night for fifteen nights. Each night has different stats and different scores that your team must beat in order to get awarded a win.",
                    externalURL: MetadataViews.ExternalURL("https://nbatopshot.com/fastbreak"),
                    squareImage: squareImage,
                    bannerImage: bannerImage,
                    socials: {
                        "twitter": MetadataViews.ExternalURL("https://twitter.com/nbatopshot"),
                        "discord": MetadataViews.ExternalURL("https://discord.com/invite/nbatopshot"),
                        "instagram": MetadataViews.ExternalURL("https://www.instagram.com/nbatopshot")
                    }
                )
        }
        return nil
    }

    access(all)  fun createPlayer(playerName: String): @FastBreakV1.Player {
        FastBreakV1.nextPlayerId = FastBreakV1.nextPlayerId + UInt64(1)

        emit FastBreakPlayerCreated(
            id: FastBreakV1.nextPlayerId,
            playerName: playerName,
        )

        return <- create FastBreakV1.Player(playerName: playerName)
    }

    /// Capabilities of the Game Oracle
    ///
    access(all) resource interface GameOracle {
        access(Create) fun createFastBreakRun(id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool)
        access(Update) fun updateFastBreakRunStatus(id: String, status: UInt8)
        access(Create) fun createFastBreakGame(
            id: String,
            name: String,
            fastBreakRunID: String,
            submissionDeadline: UInt64,
            numPlayers: UInt64
        )
        access(Update) fun updateFastBreakGame(id: String, status: UInt8, winner: UInt64)
        access(Update) fun updateFastBreakScore(fastBreakGameID: String, playerId: UInt64, points: UInt64, win: Bool)
        access(Update) fun addStatToFastBreakGame(fastBreakGameID: String, name: String, rawType: UInt8, valueNeeded: UInt64)
        access(Update) fun setSubmissionDeadline(fastBreakGameID: String, deadline: UInt64)
    }

    /// Fast Break Daemon game oracle implementation
    ///
    access(all) resource FastBreakDaemon: GameOracle {

        /// Create a Fast Break Run
        ///
        access(Create) fun createFastBreakRun(id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool) {
            let fastBreakRun = FastBreakV1.FastBreakRun(
                id: id,
                name: name,
                runStart: runStart,
                runEnd: runEnd,
                fatigueModeOn: fatigueModeOn
            )
            // Store in year-based account storage
            let year = FastBreakV1.getYearFromTimestamp(timestamp: runStart).toString()
            let yearStorage = FastBreakV1.getOrCreateYearRunStorage(year: year)
            yearStorage.runs[fastBreakRun.id] = fastBreakRun
            emit FastBreakRunCreated(
                id: fastBreakRun.id,
                name: fastBreakRun.name,
                runStart: fastBreakRun.runStart,
                runEnd: fastBreakRun.runEnd,
                fatigueModeOn: fastBreakRun.fatigueModeOn
            )
        }

        /// Update the status of a Fast Break Run
        /// Works with both legacy data (pre-upgrade) and new year-based storage
        /// Searches current year first, then previous year, then legacy
        ///
        access(Update) fun updateFastBreakRunStatus(id: String, status: UInt8) {
            // Search current year first, then previous year, then legacy
            let fastBreakRun = FastBreakV1.getFastBreakRunInternal(id: id)
            let run = fastBreakRun ?? panic("Fast break run does not exist with Id: ".concat(id))

            let runStatus: FastBreakV1.RunStatus = FastBreakV1.RunStatus(rawValue: status)
                ?? panic("Run status does not exist with rawValue: ".concat(status.toString()))

            run.updateStatus(status: runStatus)
            
            // Update in appropriate storage (updateRunInStorage handles legacy vs year-based)
            let year = FastBreakV1.getYearFromTimestamp(timestamp: run.runStart).toString()
            FastBreakV1.updateRunInStorage(run: run, runID: id, year: year)

            emit FastBreakRunStatusChange(id: run.id, newRawStatus: run.status.rawValue)
        }

        /// Create a game of Fast Break
        ///
        access(Create) fun createFastBreakGame(
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
            // Store in year-based account storage
            let year = FastBreakV1.getYearFromTimestamp(timestamp: submissionDeadline).toString()
            let yearStorage = FastBreakV1.getOrCreateYearGameStorage(year: year)
            yearStorage.games[fastBreakGame.id] = fastBreakGame
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
         access(Update) fun addStatToFastBreakGame(fastBreakGameID: String, name: String, rawType: UInt8, valueNeeded: UInt64) {

            let fastBreakGame = FastBreakV1.getFastBreakGameInternal(id: fastBreakGameID)
                ?? panic("Fast break does not exist with Id: ".concat(fastBreakGameID))

            let statType: FastBreakV1.StatisticType = FastBreakV1.StatisticType(rawValue: rawType)
                ?? panic("Fast break stat type does not exist with rawType: ".concat(rawType.toString()))

            let fastBreakStat : FastBreakV1.FastBreakStat = FastBreakV1.FastBreakStat(
                name: name,
                type: statType,
                valueNeeded: valueNeeded
            )

            fastBreakGame.addStat(stat: fastBreakStat)
            
            // Update game in account storage after mutation
            let year = FastBreakV1.getYearFromTimestamp(timestamp: fastBreakGame.submissionDeadline).toString()
            FastBreakV1.updateGameInStorage(game: fastBreakGame, gameID: fastBreakGameID, year: year)
            
            emit FastBreakGameStatAdded(
                fastBreakGameID: fastBreakGame.id,
                name: fastBreakStat.name,
                type: fastBreakStat.type.rawValue,
                valueNeeded: fastBreakStat.valueNeeded
            )

        }

        /// Set the submission deadline for a Fast Break
        ///
        access(Update) fun setSubmissionDeadline(fastBreakGameID: String, deadline: UInt64) {
            let fastBreakGame = FastBreakV1.getFastBreakGameInternal(id: fastBreakGameID)
                ?? panic("Fast break does not exist with Id: ".concat(fastBreakGameID))

            fastBreakGame.setSubmissionDeadline(deadline: deadline)
            
            // Update game in account storage after mutation
            let year = FastBreakV1.getYearFromTimestamp(timestamp: fastBreakGame.submissionDeadline).toString()
            FastBreakV1.updateGameInStorage(game: fastBreakGame, gameID: fastBreakGameID, year: year)
        }

        /// Update the status of a Fast Break
        ///
         access(Update) fun updateFastBreakGame(id: String, status: UInt8, winner: UInt64) {

            let fastBreakGame = FastBreakV1.getFastBreakGameInternal(id: id)
                ?? panic("Fast break does not exist with Id: ".concat(id))

            let fastBreakStatus: FastBreakV1.GameStatus = FastBreakV1.GameStatus(rawValue: status)
                ?? panic("Fast break status does not exist with rawValue: ".concat(status.toString()))

            fastBreakGame.update(status: fastBreakStatus, winner: winner)
            
            // Update game in account storage after mutation
            let year = FastBreakV1.getYearFromTimestamp(timestamp: fastBreakGame.submissionDeadline).toString()
            FastBreakV1.updateGameInStorage(game: fastBreakGame, gameID: id, year: year)

            emit FastBreakGameStatusChange(id: fastBreakGame.id, newRawStatus: fastBreakGame.status.rawValue)

        }

        /// Updates the submission scores of a Fast Break
        ///
        access(Update) fun updateFastBreakScore(fastBreakGameID: String, playerId: UInt64, points: UInt64, win: Bool) {
            let fastBreakGame = FastBreakV1.getFastBreakGameInternal(id: fastBreakGameID)
                ?? panic("Fast break does not exist with Id: ".concat(fastBreakGameID))

            let isNewWin = fastBreakGame.updateScore(playerId: playerId, points: points, win: win)
            
            // Update game in account storage after mutation
            let year = FastBreakV1.getYearFromTimestamp(timestamp: fastBreakGame.submissionDeadline).toString()
            FastBreakV1.updateGameInStorage(game: fastBreakGame, gameID: fastBreakGameID, year: year)

            if isNewWin {
                // Find run - searches current year first, then previous year, then legacy
                let fastBreakRun = FastBreakV1.getFastBreakRunInternal(id: fastBreakGame.fastBreakRunID)
                let run = fastBreakRun ?? panic("Could not obtain reference to fast break run with Id: ".concat(fastBreakGame.fastBreakRunID))

                run.incrementRunWinCount(playerId: playerId)
                
                // Update run in appropriate storage (updateRunInStorage handles legacy vs year-based)
                let runYear = FastBreakV1.getYearFromTimestamp(timestamp: run.runStart).toString()
                FastBreakV1.updateRunInStorage(run: run, runID: fastBreakGame.fastBreakRunID, year: runYear)
                
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
        self.PlayerStoragePath = /storage/FastBreakPlayerV1

        self.totalSupply = 0
        self.nextPlayerId = 0
        self.fastBreakRunByID = {}
        self.fastBreakGameByID = {}
        self.fastBreakPlayerByID = {}
        self.playerAccountMapping = {}
        self.accountPlayerMapping = {}

        let oracle <- create FastBreakDaemon()
        self.account.storage.save(<-oracle, to: self.OracleStoragePath)

    }
}
