import TopShot from 0x01

transaction {
    prepare(acct: Account) {
        
    }

    execute {
        TopShot.castMold(metadata: {"Name": "Lebron"}, qualityCounts: {1: 1, 2: 2, 3: 3, 4: 4, 5: 5})
    }
}