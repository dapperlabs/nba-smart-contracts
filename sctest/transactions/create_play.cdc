import TopShot from 0x02

transaction {

    let adminRef: &TopShot.Admin

    prepare(acct: Account) {
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin
    }

    execute {
        
        let id1 = self.adminRef.createPlayData(metadata: {"Name": "Lebron"})

        let id2 = self.adminRef.createPlayData(metadata: {"Name": "Oladipo"})

        log("PlayData 1 and 2 Succcesfully created!")
    }
}
 