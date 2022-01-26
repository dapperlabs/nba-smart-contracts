import TopShot from 0xTOPSHOTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS

pub struct TopShotMomentMetadataView {

    pub let fullName: String?
    pub let firstName: String?
    pub let lastName: String?
    pub let birthdate: String?
    pub let birthplace: String?
    pub let jerseyNumber: String?
    pub let draftTeam: String?
    pub let draftYear: String?
    pub let draftSelection: String?
    pub let draftRound: String?
    pub let teamAtMomentNBAID: String?
    pub let teamAtMoment: String?
    pub let primaryPosition: String?
    pub let height: String?
    pub let weight: String?
    pub let totalYearsExperience: String?
    pub let nbaSeason: String?
    pub let dateOfMoment: String?
    pub let playCategory: String?
    pub let playType: String?
    pub let homeTeamName: String?
    pub let awayTeamName: String?
    pub let homeTeamScore: String?
    pub let awayTeamScore: String?
    pub let seriesNumber: UInt32?
    pub let setName: String?
    pub let serialNumber: UInt32
    pub let playID: UInt32
    pub let setID: UInt32
    pub let numMomentsInEdition: UInt32?

    init(
        fullName: String?,
        firstName: String?,
        lastName: String?,
        birthdate: String?,
        birthplace: String?,
        jerseyNumber: String?,
        draftTeam: String?,
        draftYear: String?,
        draftSelection: String?,
        draftRound: String?,
        teamAtMomentNBAID: String?,
        teamAtMoment: String?,
        primaryPosition: String?,
        height: String?,
        weight: String?,
        totalYearsExperience: String?,
        nbaSeason: String?,
        dateOfMoment: String?,
        playCategory: String?,
        playType: String?,
        homeTeamName: String?,
        awayTeamName: String?,
        homeTeamScore: String?,
        awayTeamScore: String?,
        seriesNumber: UInt32?,
        setName: String?,
        serialNumber: UInt32,
        playID: UInt32,
        setID: UInt32,
        numMomentsInEdition: UInt32?
    ) {
        self.fullName = fullName
        self.firstName = firstName
        self.lastName = lastName
        self.birthdate = birthdate
        self.birthplace = birthplace
        self.jerseyNumber = jerseyNumber
        self.draftTeam = draftTeam
        self.draftYear = draftYear
        self.draftSelection = draftSelection
        self.draftRound = draftRound
        self.teamAtMomentNBAID = teamAtMomentNBAID
        self.teamAtMoment = teamAtMoment
        self.primaryPosition = primaryPosition
        self.height = height
        self.weight = weight
        self.totalYearsExperience = totalYearsExperience
        self.nbaSeason = nbaSeason
        self.dateOfMoment= dateOfMoment
        self.playCategory = playCategory
        self.playType = playType
        self.homeTeamName = homeTeamName
        self.awayTeamName = awayTeamName
        self.homeTeamScore = homeTeamScore
        self.awayTeamScore = awayTeamScore
        self.seriesNumber = seriesNumber
        self.setName = setName
        self.serialNumber = serialNumber
        self.playID = playID
        self.setID = setID
        self.numMomentsInEdition = numMomentsInEdition
    }
}

pub fun main(address: Address, id: UInt64): TopShotMomentMetadataView {
    let account = getAccount(address)

    let collectionRef = account.getCapability(/public/MomentCollection)
                            .borrow<&{TopShot.MomentCollectionPublic}>()!

    let nft = collectionRef.borrowMoment(id: id)!
    
    // Get the Top Shot specific metadata for this NFT
    let view = nft.resolveView(Type<TopShot.TopShotMomentMetadataView>())!

    let metadata = view as! TopShot.TopShotMomentMetadataView
    
    return TopShotMomentMetadataView(
        fullName: metadata.fullName,
        firstName: metadata.firstName,
        lastName: metadata.lastName,
        birthdate: metadata.lastName,
        birthplace: metadata.birthplace,
        jerseyNumber: metadata.jerseyNumber,
        draftTeam: metadata.draftTeam,
        draftYear: metadata.draftYear,
        draftSelection: metadata.draftSelection,
        draftRound: metadata.draftRound,
        teamAtMomentNBAID: metadata.teamAtMomentNBAID,
        teamAtMoment: metadata.teamAtMoment,
        primaryPosition: metadata.primaryPosition,
        height: metadata.height,
        weight: metadata.weight,
        totalYearsExperience: metadata.totalYearsExperience,
        nbaSeason: metadata.nbaSeason,
        dateOfMoment: metadata.dateOfMoment,
        playCategory: metadata.playCategory,
        playType: metadata.playType,
        homeTeamName: metadata.homeTeamName,
        awayTeamName: metadata.awayTeamName,
        homeTeamScore: metadata.homeTeamScore,
        awayTeamScore: metadata.awayTeamScore,
        seriesNumber: metadata.seriesNumber,
        setName: metadata.setName,
        serialNumber: metadata.serialNumber,
        playID: metadata.playID,
        setID: metadata.setID,
        numMomentsInEdition: metadata.numMomentsInEdition
    )
}