import Market from 0x04

// this is the transacion a user would run if they want
// to change the price of one of their tokens that is for sale

// values that would be configurable are
// id of the token for sale 
// new price of the token
transaction {

    prepare(acct: AuthAccount) {

        // remove the sale collection from storage
        if acct.borrow<&Market.SaleCollection>(from: /storage/MomentSale) == nil {

            // create a temporary reference to the sale
            let saleRef = acct.borrow<&Market.SaleCollection>(from: /storage/MomentSale)!

            // put the token up for sale
            saleRef.changePrice(tokenID: 1, newPrice: 40.00)

            log("Token price changed")
            log(1)

        } else {
            panic("No sale!")
        }
    }
}
 