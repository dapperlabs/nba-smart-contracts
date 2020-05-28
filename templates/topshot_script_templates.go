package templates

import (
	"errors"
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

func GenerateChallengeCompletedScript(tokenAddr, userAddress flow.Address, setIDs []uint32, playIDs []uint32) ([]byte, error) {
	if len(setIDs) != len(playIDs) {
		return nil, errors.New("set and play ID arrays have mismatched lengths")
	}

	template := `
		import TopShot from 0x%s
		
		fun main(): Int {
			let setIDs = [%s]
			let playIDs = [%s]
			assert(
				setIDs.length == playIDs.length,
				message: "set and play ID arrays have mismatched lengths"
			)
		
			let collectionRef = getAccount(0x%s).getCapability(/public/MomentCollection)!
						.borrow<&{TopShot.MomentCollectionPublic}>()
						?? panic("Could not get public moment collection reference")
			let momentIDs = collectionRef.getIDs()
		
			var numMatchingMoments = 0
			var i = 0
			while i < setIDs.length {
				for momentID in momentIDs {
					let moment = collectionRef.borrowNFT(id: momentID)
					let setID = moment.data.setID
					let playID = moment.data.playID
					if setID == setIDs[i] && playID == playIDs[i] {
						numMatchingMoments = numMatchingMoments + 1
						break
					}
				}
				i = i + 1
			}
}`
	return []byte(fmt.Sprintf(template, tokenAddr, stringifyUint32Slice(setIDs), stringifyUint32Slice(playIDs), userAddress)), nil
}

func stringifyUint32Slice(ints []uint32) string {
	intArrayStr := ""
	for _, i := range ints {
		intStr := strconv.Itoa(int(i))
		intArrayStr = intArrayStr + `UInt32(` + intStr + `), `
	}
	// Remove comma and space from last entry
	if arrayLen := len(intArrayStr); arrayLen > 2 {
		intArrayStr = intArrayStr[:len(intArrayStr)-2]
	}
	return intArrayStr
}
