import TopShot from 0xTOPSHOTADDRESS
import TopshotAdminReceiver from 0xADMINRECEIVERADDRESS

transaction {

    // Local variable for the topshot Admin object
    let adminRef: @TopShot.Admin

    prepare(acct: AuthAccount) {

        self.adminRef <- acct.load<@TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No topshot admin in storage")
    }

    execute {

        TopshotAdminReceiver.storeAdmin(newAdmin: <-self.adminRef)
        
    }
}