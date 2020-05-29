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

// GenerateReturnAllPlaysScript creates a script that returns an array
// of all the plays that have been created
func GenerateReturnAllPlaysScript(tokenAddr flow.Address) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): [TopShot.Play] {
			return TopShot.playDatas.values
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr))
}

// GenerateReturnPlayMetadataScript creates a script that returns the metadata of a play
func GenerateReturnPlayMetadataScript(tokenAddr flow.Address, playID int) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): {String: String} {
			return TopShot.getPlayMetaData(playID: %d)!
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, playID))
}

// GenerateReturnSetSeriesScript creates a script that returns the metadata of a play
func GenerateReturnSetSeriesScript(tokenAddr flow.Address, playID int) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): {String: String} {
			return TopShot.getPlayMetaData(playID: %d)!
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, playID))
}

// GenerateReturnPlaysInSetScript creates a script that returns an array of plays in a set
func GenerateReturnPlaysInSetScript(tokenAddr flow.Address, setID int) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): [UInt32] {
			return TopShot.getPlaysInSet(setID: %d)!
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, setID))
}

// GenerateReturnIsEditionRetiredScript creates a script that indicates if an edition is retired
func GenerateReturnIsEditionRetiredScript(tokenAddr flow.Address, setID, playID int) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): Bool {
			return TopShot.isEditionRetired(setID: %d playID: %d)!
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, setID, playID))
}

// GenerateReturnIsSetLockedScript creates a script that indicates if a set is locked
func GenerateReturnIsSetLockedScript(tokenAddr flow.Address, setID int) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): Bool {
			return TopShot.isSetLocked(setID: %d)!
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, setID))
}

// GenerateGetNumMomentsInEditionScript creates a script
// that returns the number of moments that have been minted in an edition
func GenerateGetNumMomentsInEditionScript(tokenAddr flow.Address, setID, playID int) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): UInt32 {
			return TopShot.getNumMomentsInEdition(setID: %d, playID: %d)!
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, setID, playID))
}
