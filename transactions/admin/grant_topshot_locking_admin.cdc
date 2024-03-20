import TopShot from 0xTOPSHOTADDRESS
import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

// this transaction takes a TopShot Locking Admin resource and
// saves it to the account storage of the account second authorizer

transaction {
    prepare(acct: auth(BorrowValue) &Account, acct2: auth(SaveValue) &Account) {
        let topShotLockingAdmin = acct.storage.borrow<&TopShotLocking.Admin>(from: TopShotLocking.AdminStoragePath())
          ?? panic("could not borrow admin reference")

        acct2.storage.save(<- topShotLockingAdmin.createNewAdmin(), to: TopShotLocking.AdminStoragePath())
    }
}
