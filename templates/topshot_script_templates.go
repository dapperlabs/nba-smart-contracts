package templates

import (
	"fmt"
	"strconv"

	"github.com/onflow/flow-go-sdk"
)

// GenerateInspectTopshotFieldScript creates a script that checks
// a field of the topshot contract
func GenerateInspectTopshotFieldScript(nftAddr, tokenAddr flow.Address, fieldName, fieldType string, expectedValue int) []byte {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s

		pub fun main() {
			assert(
                TopShot.%s == %s(%d),
                message: "incorrect %s"
            )
		}
	`

	return []byte(fmt.Sprintf(template, nftAddr, tokenAddr, fieldName, fieldType, expectedValue, fieldName))
}

// GenerateInspectCollectionScript creates a script that checks
// a collection for a certain ID
func GenerateInspectCollectionScript(nftAddr, tokenAddr, ownerAddr flow.Address, expectedID int) []byte {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s

		pub fun main() {
			let collectionRef = getAccount(0x%s).getCapability(/public/MomentCollection)!
				.borrow<&{TopShot.MomentCollectionPublic}>()
				?? panic("Could not get public moment collection reference")

			assert(
                collectionRef.borrowNFT(id: %d).id == UInt64(%d),
                message: "ID %d does not exist in the collection"
            )
		}
	`

	return []byte(fmt.Sprintf(template, nftAddr, tokenAddr, ownerAddr, expectedID, expectedID, expectedID))
}

// GenerateInspectCollectionIDsScript creates a script that checks
// a collection for a certain ID
func GenerateInspectCollectionIDsScript(nftAddr, tokenAddr, ownerAddr flow.Address, momentIDs []uint64) []byte {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s

		pub fun main() {
			let collectionRef = getAccount(0x%s).getCapability(/public/MomentCollection)!
				.borrow<&{TopShot.MomentCollectionPublic}>()
				?? panic("Could not get public moment collection reference")

			let ids = collectionRef.getIDs()

			let expectedIDs = [%s]

			assert(
				ids.length == expectedIDs.length,
				message: "ID array is not the expected length"
			)

			var i = 0
			for element in ids {
				if element != expectedIDs[i] {
					panic("Unexpected ID in the array")
				}
				i = i + 1
			}
		}
	`

	// Stringify moment IDs
	momentIDList := ""
	for _, momentID := range momentIDs {
		id := strconv.Itoa(int(momentID))
		momentIDList = momentIDList + `UInt64(` + id + `), `
	}
	// Remove comma and space from last entry
	if idListLen := len(momentIDList); idListLen > 2 {
		momentIDList = momentIDList[:len(momentIDList)-2]
	}

	return []byte(fmt.Sprintf(template, nftAddr, tokenAddr, ownerAddr, momentIDList))
}
