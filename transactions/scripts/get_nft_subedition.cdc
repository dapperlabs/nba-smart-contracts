import TopShot from 0xTOPSHOTADDRESS

pub fun main(account: Address, nftID: UInt64): UInt32 {

     let publicSubeditionRef = getAccount(account).getCapability(/public/PublicSubedition)
        .borrow<&{TopShot.PublicSubedition}>()
        ?? panic("Could not get public subedition reference")

    let subedition = publicSubeditionRef.getMomentsSubedition(nftID: nftID)
                ?? panic("Could not find the specified moment")
    return subedition
}