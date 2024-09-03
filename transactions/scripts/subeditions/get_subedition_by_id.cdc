import TopShot from 0xTOPSHOTADDRESS

// This script returns the full Subedition entity from
// the TopShot smart contract

// Parameters:
//
// subeditionID: The unique ID for the subedition whose data needs to be read

// Returns: Subedition
// struct from TopShot contract

access(all) fun main(subeditionID: UInt32): &TopShot.Subedition {

    let subedititon = TopShot.getSubeditionByID(subeditionID: subeditionID)

    return subedititon
}