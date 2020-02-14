import TopShot from 0x02

transaction {
    prepare(acct: Account) {
        let adminRef = acct.storage[&TopShot.Admin] ?? panic("No admin!")
        
        let id1 = adminRef.castMold(metadata: {"Name": "Lebron"}, qualityCounts: {1: UInt32(1), 2: UInt32(2), 3: UInt32(3), 4: UInt32(4), 5: UInt32(5), 6: UInt32(6), 7: UInt32(7), 8: UInt32(8)})
        
        let id2 = adminRef.castMold(metadata: {"Name": "Oladipo"}, qualityCounts: {1: UInt32(0), 2: UInt32(1), 3: UInt32(10), 4: UInt32(20), 5: UInt32(40), 6: UInt32(40), 7: UInt32(40), 8: UInt32(40)})
    }
}
 