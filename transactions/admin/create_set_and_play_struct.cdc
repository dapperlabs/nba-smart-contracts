import TopShot from 0xTOPSHOTADDRESS

transaction() {
    
    prepare(acct: AuthAccount) {

        let metadata: {String: String} = {"PlayType": "Shoe becomes untied"}

        let newPlay = TopShot.Play(metadata: metadata)

        let newSet = TopShot.SetData(name: "Sneaky Sneakers")
    }
}