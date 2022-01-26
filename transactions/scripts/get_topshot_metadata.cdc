import TopShot from 0xTOPSHOTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS


pub fun main(address: Address, id: UInt64): TopShot.TopShotMomentMetadataView {
    let account = getAccount(address)

    let collectionRef = account.getCapability(/public/MomentCollection)
                            .borrow<&{TopShot.MomentCollectionPublic}>()!

    let nft = collectionRef.borrowMoment(id: id)!
    
    // Get the Top Shot specific metadata for this NFT
    let view = nft.resolveView(Type<TopShot.TopShotMomentMetadataView>())!

    let metadata = view as! TopShot.TopShotMomentMetadataView
    
    return TopShot.TopShotMomentMetadataView(
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