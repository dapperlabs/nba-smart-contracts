pub contract TopShot {

    access(self) var sets: @{Int: Set}

    pub resource Set {
        pub let setID: Int

        init() {
            self.setID = 0
        }
    }

    pub fun borrowSet(setID: Int): &Set {
        return &TopShot.sets[setID] as &Set
    }

    init() {
        self.sets <- {}
        self.sets[0] <-! create Set()

        let setRef = self.borrowSet(setID: 0)

        log(setRef.setID)
    }
}
 