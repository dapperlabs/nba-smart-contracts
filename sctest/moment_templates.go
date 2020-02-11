package sctest

import (
	"fmt"

	"github.com/dapperlabs/flow-go-sdk/model/flow"
)

// GenerateMintMoldScript creates a script that mints a new mold in the collection
func GenerateCastMoldScript(tokenCodeAddr flow.Address, name string, numInEdition int) []byte {
	template := `
`

	return []byte(fmt.Sprintf(template, tokenCodeAddr.String(), name, numInEdition))
}

// GenerateMintMomentFactoryScript creates a script that creates a new moment factory resource
func GenerateMintMomentFactoryScript(tokenCodeAddr flow.Address) []byte {
	template := ``

	return []byte(fmt.Sprintf(template, tokenCodeAddr.String()))
}

// GenerateCreateMomentCollectionScript creates a new collection in the signers acount using
// the moment factory collection creation method
func GenerateCreateMomentCollectionScript(tokenCodeAddr flow.Address, userAddr flow.Address) []byte {
	template := ``

	return []byte(fmt.Sprintf(template, tokenCodeAddr.String(), userAddr))
}

// GenerateMintMomentScript creates a script that mints a new moment in the collection
func GenerateMintMomentScript(tokenCodeAddr flow.Address, moldID int, rarity string, recipientAddress flow.Address) []byte {
	template := ``

	return []byte(fmt.Sprintf(template, tokenCodeAddr.String(), recipientAddress, moldID, rarity))
}

// GenerateInspectMomentScript creates a script that retrieves a moment collection
// from storage and checks to see if an ID exists there
func GenerateInspectMomentScript(nftCodeAddr flow.Address, momentID int, name string) []byte {
	template := ``

	return []byte(fmt.Sprintf(template, nftCodeAddr, momentID, momentID, name))
}
