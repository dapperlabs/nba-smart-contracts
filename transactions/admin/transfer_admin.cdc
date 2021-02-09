import TopShot from 0xTOPSHOTADDRESS
import TopshotAdminReceiver from 0xADMINRECEIVERADDRESS

transaction {

    prepare(acct: AuthAccount) {
        let admin <- acct.load<@TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No topshot admin in storage")

        TopshotAdminReceiver.storeAdmin(newAdmin: <-admin)
    }
}