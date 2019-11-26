package sctest

import (
	"fmt"

	"github.com/dapperlabs/flow-go/model/flow"
)

// GenerateCreateMoldCollectionScript Creates a script that instantiates a new
// MoldCollection instance
func GenerateCreateMoldCollectionScript(tokenAddr flow.Address) []byte {
	template := `
		import MoldCollection, createMoldCollection from 0x%s

		pub fun main(acct: Account) {
			var collection: <-MoldCollection <- createMoldCollection()
			
			let oldCollection <- acct.storage[MoldCollection] <- collection

			acct.storage[&MoldCollection] = &acct.storage[MoldCollection] as MoldCollection

			destroy oldCollection
		}`
	return []byte(fmt.Sprintf(template, tokenAddr))
}

// GenerateMintMoldScript creates a script that mints a new mold in the collection
func GenerateMintMoldScript(tokenCodeAddr flow.Address, name string, numInEdition int) []byte {
	template := `
		import Mold, MoldCollection from 0x%s

		pub fun main(acct: Account) {

			let collectionRef = acct.storage[&MoldCollection] ?? panic("missing mold collection reference")

			collectionRef.castMold(name: %s, rarityCounts: {"Uncommon": %d, "Rare": 10, "Epic": 1, "Elite": 0, "Legendary": 0})
		}`

	return []byte(fmt.Sprintf(template, tokenCodeAddr.String(), name, numInEdition))
}

// GenerateMintMomentFactoryScript creates a script that creates a new moment factory resource
func GenerateMintMomentFactoryScript(tokenCodeAddr flow.Address) []byte {
	template := `
		import Mold, MoldCollection, MomentFactory, createMomentFactory from 0x%s

		pub fun main(acct: Account) {

			let collectionRef = acct.storage[&MoldCollection] ?? panic("missing mold collection reference")

			var oldfactory <- acct.storage[MomentFactory] <- createMomentFactory(ref: collectionRef)
			destroy oldfactory

			acct.storage[&MomentFactory] = &acct.storage[MomentFactory] as MomentFactory
		}`

	return []byte(fmt.Sprintf(template, tokenCodeAddr.String()))
}

// GenerateCreateMomentCollectionScript creates a new collection in the signers acount using
// the moment factory collection creation method
func GenerateCreateMomentCollectionScript(tokenCodeAddr flow.Address, userAddr flow.Address) []byte {
	template := `
		import Mold, MoldCollection, MomentFactory, MomentCollection from 0x%s

		pub fun main(acct: Account) {

			let factoryAddr = getAccount("%s")
			let factoryRef = factoryAddr.storage[&MomentFactory] ?? panic("missing factory reference")

			var oldCollection <- acct.storage[MomentCollection] <- factoryRef.createMomentCollection()
			destroy oldCollection

			acct.storage[&MomentCollection] = &acct.storage[MomentCollection] as MomentCollection
		}`

	return []byte(fmt.Sprintf(template, tokenCodeAddr.String(), userAddr))
}

// GenerateMintMomentScript creates a script that mints a new moment in the collection
func GenerateMintMomentScript(tokenCodeAddr flow.Address, moldID int, rarity string, recipientAddress flow.Address) []byte {
	template := `
		import Mold, MoldCollection, MomentFactory, MomentCollection from 0x%s

		pub fun main(acct: Account) {
			let collectionAddr = getAccount("%s")
			let collectionRef = collectionAddr.storage[&MomentCollection] ?? panic("missing collection reference")

			let factoryRef = acct.storage[&MomentFactory] ?? panic("missing moment factory reference")

			factoryRef.mintMoment(moldID: %d, rarity: "%s", recipient: collectionRef)
		}`

	return []byte(fmt.Sprintf(template, tokenCodeAddr.String(), recipientAddress, moldID, rarity))
}

// GenerateInspectMoldCollectionScript creates a script that retrieves a mold collection
// from storage and makes assertions about the mold IDs that it contains
func GenerateInspectMoldCollectionScript(nftCodeAddr, userAddr flow.Address, nftID int, shouldExist bool) []byte {
	template := `
		import Mold, MoldCollection from 0x%s

		pub fun main() {
			let acct = getAccount("%s")
			let collectionRef = acct.storage[&MoldCollection] ?? panic("missing collection reference")

			if collectionRef.molds[%d] == nil {
				if %v {
					panic("Token ID doesn't exist!")
				}
			}
		}`

	return []byte(fmt.Sprintf(template, nftCodeAddr, userAddr, nftID, shouldExist))
}

// GenerateInspectMomentScript creates a script that retrieves a moment collection
// from storage and checks to see if an ID exists there
func GenerateInspectMomentScript(nftCodeAddr flow.Address, momentID int, name string) []byte {
	template := `
	import Mold, MoldCollection, MomentFactory, MomentCollection, Moment from 0x%s

	pub fun main(acct: Account) {

		let collectionRef = acct.storage[&MomentCollection] ?? panic("missing moment collection")
		
		let moldRef = collectionRef.moments[%d]?.moldReference ?? panic("missing moment mold reference!")
		let moldID = collectionRef.moments[%d]?.moldID ?? panic("missing moment mold ID!")

		if moldRef.getMoldName(moldID: moldID) != "%s" {
			panic("mold name is incorrect!")
		}
	}`

	return []byte(fmt.Sprintf(template, nftCodeAddr, momentID, momentID, name))
}
