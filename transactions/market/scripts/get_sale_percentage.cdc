import Market from 0xMARKETADDRESS

pub fun main(sellerAddress: Address): UFix64 {
    let acct = getAccount(sellerAddress)
    let collectionRef = acct.getCapability(/public/topshotSaleCollection).borrow<&{Market.SalePublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.cutPercentage
}