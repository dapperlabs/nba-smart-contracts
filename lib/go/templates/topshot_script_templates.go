package templates

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/onflow/flow-go-sdk"
)

const (
	scriptsPath = "../../../transactions/scripts/"

	// Topshot contract scripts
	currentSeriesFilename = "read_currentSeries.cdc"
	totalSupplyFilename   = "read_totalSupply.cdc"

	// Play related scripts
	getAllPlaysFilename = "plays/get_all_plays.cdc"
	nextPlayIDFilename  = "read_nextPlayID.cdc"
	playMetadata        = "plays/read_play_metadata.cdc"
	playMetadataField   = "plays/read_play_metadata_field.cdc"

	// Set related scripts
	editionRetiredFilename      = "sets/read_edition_retired.cdc"
	numMomentsInEditionFilename = "sets/read_numMoments_in_edition.cdc"
	setIDsByNameFilename        = "sets/read_setIDs_by_name.cdc"
	setSeriesFilename           = "sets/read_setSeries.cdc"
	nextSetIDFilename           = "sets/read_nextSetID.cdc"
	playsInSetFilename          = "sets/read_plays_in_set.cdc"
	setNameFilename             = "sets/read_setName.cdc"
	setLockedFilename           = "sets/read_set_locked.cdc"

	// collections scripts
	collectionIDsFilename   = "collections/get_collection_ids.cdc"
	metadataFieldFilename   = "collections/get_metadata_field.cdc"
	momentSeriesFilename    = "collections/get_moment_series.cdc"
	idInCollectionFilename  = "collections/get_id_in_Collection.cdc"
	momentPlayIDFilename    = "collections/get_moment_playID.cdc"
	momentSetIDFilename     = "collections/get_moment_setID.cdc"
	metadataFilename        = "collections/get_metadata.cdc"
	momentSerialNumFilename = "collections/get_moment_serialNum.cdc"
	momentSetNameFilename   = "collections/get_moment_setName.cdc"
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

// GenerateInspectCollectionDataScript creates a script that checks
// a collection for a certain ID
func GenerateInspectCollectionDataScript(nftAddr, tokenAddr, ownerAddr flow.Address, expectedID, expectedSet int) []byte {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s

		pub fun main() {
			let collectionRef = getAccount(0x%s).getCapability(/public/MomentCollection)!
				.borrow<&{TopShot.MomentCollectionPublic}>()
				?? panic("Could not get public moment collection reference")

			let token = collectionRef.borrowMoment(id: %d)
				?? panic("Could not borrow a reference to the specified moment")

			let data = token.data

			assert(
				data.setID == UInt32(%d),
				message: "ID %d does not have the expected Set ID %d"
			)
		}
	`

	return []byte(fmt.Sprintf(template, nftAddr, tokenAddr, ownerAddr, expectedID, expectedSet, expectedID, expectedSet))
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

			if ids.length != 0 {
				var i = 0
				for element in ids {
					if element != expectedIDs[i] {
						panic("Unexpected ID in the array")
					}
					i = i + 1
				}
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
			return TopShot.getAllPlays()
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr))
}

// GenerateReturnPlayMetadataScript creates a script that returns the metadata of a play
func GenerateReturnPlayMetadataScript(tokenAddr flow.Address, playID int, expectedKey, expectedValue string) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): {String: String} {
			let metadata = TopShot.getPlayMetaData(playID: %d)!

			assert (
				metadata["%s"] == "%s",
				message: "Key Value is incorrect"
			)

			return metadata
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, playID, expectedKey, expectedValue))
}

// GenerateReturnPlayMetadataByFieldScript creates a script that returns the metadata of a play
func GenerateReturnPlayMetadataByFieldScript(tokenAddr flow.Address, playID int, expectedKey, expectedValue string) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): String {
			let metadata = TopShot.getPlayMetaDataByField(playID: %d, field: "%s")!

			assert (
				metadata == "%s",
				message: "Field Value is incorrect"
			)

			return metadata
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, playID, expectedKey, expectedValue))
}

// GenerateReturnSetNameScript creates a script that returns the metadata of a play
func GenerateReturnSetNameScript(tokenAddr flow.Address, setID int, expectedName string) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): String {
			let name = TopShot.getSetName(setID: %d)!

			assert (
				name == "%s",
				message: "Set name is incorrect"
			)

			return name
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, setID, expectedName))
}

// GenerateReturnSetIDsByNameScript creates a script that returns the metadata of a play
func GenerateReturnSetIDsByNameScript(tokenAddr flow.Address, setName string, expectedID int) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): [UInt32] {
			let ids = TopShot.getSetIDsByName(setName: "%s")!

			assert (
				ids[0] == UInt32(%d),
				message: "Set id is incorrect"
			)

			return ids
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, setName, expectedID))
}

// GenerateReturnSetSeriesScript creates a script that returns the metadata of a play
func GenerateReturnSetSeriesScript(tokenAddr flow.Address, setID int, expectedSeries int) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): UInt32 {
			let series = TopShot.getSetSeries(setID: %d)!

			assert (
				series == UInt32(%d),
				message: "Set series is incorrect"
			)

			return series
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, setID, expectedSeries))
}

// GenerateReturnPlaysInSetScript creates a script that returns an array of plays in a set
func GenerateReturnPlaysInSetScript(tokenAddr flow.Address, setID int, expectedPlays []int) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): [UInt32] {
			let plays = TopShot.getPlaysInSet(setID: %d)!

			let expectedPlays = [%s]

			assert(
				plays.length == expectedPlays.length,
				message: "Play ID array is not the expected length"
			)

			var i = 0
			for playID in plays {
				if playID != expectedPlays[i] {
					panic("Unexpected ID in the array")
				}
				i = i + 1
			}

			return plays
		}
	`

	// Stringify IDs
	IDList := ""
	for _, ID := range expectedPlays {
		id := strconv.Itoa(int(ID))
		IDList = IDList + `UInt32(` + id + `), `
	}
	// Remove comma and space from last entry
	if idListLen := len(IDList); idListLen > 2 {
		IDList = IDList[:len(IDList)-2]
	}

	return []byte(fmt.Sprintf(template, tokenAddr, setID, IDList))
}

// GenerateReturnIsEditionRetiredScript creates a script that indicates if an edition is retired
func GenerateReturnIsEditionRetiredScript(tokenAddr flow.Address, setID, playID int, expectedResult string) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): Bool {
			let isRetired = TopShot.isEditionRetired(setID: %d, playID: %d)!

			assert (
				isRetired == %s,
				message: "isRetired is incorrect"
			)

			return isRetired
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, setID, playID, expectedResult))
}

// GenerateReturnIsSetLockedScript creates a script that indicates if a set is locked
func GenerateReturnIsSetLockedScript(tokenAddr flow.Address, setID int, expectedResult string) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): Bool {
			let isLocked = TopShot.isSetLocked(setID: %d)!

			assert (
				isLocked == %s,
				message: "isLocked is incorrect"
			)

			return isLocked
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, setID, expectedResult))
}

// GenerateGetNumMomentsInEditionScript creates a script
// that returns the number of moments that have been minted in an edition
func GenerateGetNumMomentsInEditionScript(tokenAddr flow.Address, setID, playID int, expectedMoments int) []byte {
	template := `
		import TopShot from 0x%s

		pub fun main(): UInt32 {
			let numMoments = TopShot.getNumMomentsInEdition(setID: %d, playID: %d)!

			assert (
				numMoments == UInt32(%d),
				message: "Number of moments in the edition is incorrect"
			)

			return numMoments
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, setID, playID, expectedMoments))
}

func GenerateChallengeCompletedScript(tokenAddr, userAddress flow.Address, setIDs []uint32, playIDs []uint32) ([]byte, error) {
	if len(setIDs) != len(playIDs) {
		return nil, errors.New("set and play ID arrays have mismatched lengths")
	}
	if len(setIDs) == 0 {
		return nil, errors.New("no SetPlays specified")
	}

	template := `
		import TopShot from 0x%s
		
		pub fun main(): Bool {
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
					let token = collectionRef.borrowMoment(id: momentID)
						?? panic("Could not borrow a reference to the specified moment")

					let momentData = token.data
					let setID = momentData.setID
					let playID = momentData.playID
					if setID == setIDs[i] && playID == playIDs[i] {
						numMatchingMoments = numMatchingMoments + 1
						break
					}
				}
				i = i + 1
			}
			
			return numMatchingMoments == setIDs.length
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
