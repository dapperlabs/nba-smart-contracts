import TopShot from 0x02

transaction {

    let adminRef: &TopShot.Admin

    prepare(acct: Account) {
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin
    }

    execute {
        
        let id1 = self.adminRef.castMold(metadata: {"Name": "Lebron"}, 
                                         qualityCounts: [UInt32(3000000000),
                                                         UInt32(1000000000),
                                                         UInt32(0), UInt32(0), 
                                                         UInt32(0), UInt32(0),
                                                         UInt32(0), UInt32(0), 
                                                         UInt32(0), UInt32(0),
                                                         UInt32(100), 
                                                         UInt32(0), UInt32(0), 
                                                         UInt32(10), UInt32(0),
                                                         UInt32(3)])

        let id2 = self.adminRef.castMold(metadata: {"Name": "Oladipo"}, 
                                         qualityCounts: [UInt32(3000000000),
                                                         UInt32(1000000000),
                                                         UInt32(0), UInt32(0), 
                                                         UInt32(0), UInt32(0),
                                                         UInt32(0), UInt32(0), 
                                                         UInt32(0), UInt32(0),
                                                         UInt32(100),
                                                         UInt32(0), UInt32(0),
                                                         UInt32(10), UInt32(0),
                                                         UInt32(3)])

        log("Molds 1 and 2 Succcesfully cast!")
    }
}
 