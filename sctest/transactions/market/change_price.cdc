import Market from 0x03

// this is the transacion a user would run if they want
// to put one of their moments up for sale

// values that would be configurable are
// withdrawID of the token for sale 
// price of the token
transaction {

    prepare(acct: AuthAccount) {

        // remove the sale collection from storage
        if acct.storage[Market.SaleCollection] != nil {

            let saleRef = &acct.storage[Market.SaleCollection] as &Market.SaleCollection

            // put the token up for sale
            saleRef.changePrice(tokenID: 1, newPrice: 30)

            log("Token price changed")
            log(1)

        } else {
            panic("No sale!")
        }
    }
}
 