import TopShot from 0xTOPSHOTADDRESS

pub fun main(account: Address, nftID: UInt64): UInt32 {

     let publicSubEditionRef = getAccount(account).getCapability(/public/PublicSubEdition)
        .borrow<&{TopShot.PublicSubEdition}>()
        ?? panic("Could not get public subEdition reference")

    let subEdition = publicSubEditionRef.getMomentsSubEdition(nftID: nftID)
                ?? panic("Could not find the specified moment")
    return subEdition
}