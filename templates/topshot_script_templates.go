package templates

import (
	"fmt"

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
                collectionRef.borrowNFT(id: UInt64(%d)).id == UInt64(%d),
                message: "ID %d does not exist in the collection"
            )
		}
	`

	return []byte(fmt.Sprintf(template, nftAddr, tokenAddr, ownerAddr, expectedID, expectedID, expectedID))
}
