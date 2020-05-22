package templates

import (
	"fmt"

	"github.com/onflow/flow-go-sdk"
)

// GenerateInspectTopshotFieldScript creates a script that checks
// a field of the topshot contract
func GenerateInspectTopshotFieldScript(nftAddr, tokenAddr flow.Address, fieldName string, expectedSeries int) []byte {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s

		pub fun main() {
			assert(
                TopShot.%s == UInt64(%d),
                message: "incorrect %s"
            )
		}
	`

	return []byte(fmt.Sprintf(template, nftAddr, tokenAddr, fieldName, expectedSeries, fieldName))
}
