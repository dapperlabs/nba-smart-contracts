import TopShot from 0x02

transaction {
    prepare(acct: Account) {
        let adminRef = acct.storage[&TopShot.Admin] ?? panic("No admin!")
        
        let id1 = adminRef.castMold(metadata: {"Name": "Lebron"}, qualityCounts: [UInt32(1), UInt32(2), UInt32(3), UInt32(4), UInt32(5), UInt32(6), UInt32(7), UInt32(8)])
        
        let id2 = adminRef.castMold(metadata: {"Name": "Oladipo"}, qualityCounts: [UInt32(0), UInt32(1), UInt32(10), UInt32(20), UInt32(40), UInt32(40), UInt32(40), UInt32(40)])
    }
}
 