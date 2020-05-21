/*

  AdminReceiver.cdc

  This contract defines a resource that can hold 
  a Topshot Admin resource capability
  Allows another account to transfer 
  an Admin resource capability to a different account
  without requiring multiple transaction signers.

 */

import TopShot from 0x03

pub contract TopshotAdminReceiver {

    // The owner will publish the AdminHolder with this interface
    // to be able to receive an Admin from a different account
    pub resource interface Receiver {
        pub fun setAdmin(newAdminCapability: Capability)
    }

    // Resource object that holds the Admin resource capability
    pub resource AdminHolder: Receiver {

        // Field to hold the admin resource capability
        pub var admin: Capability?

        // Function that the other account will call 
        // to transfer the Admin capability
        pub fun setAdmin(newAdminCapability: Capability) {
            pre {
                self.admin == nil: "Admin already is set"
                newAdminCapability.borrow<&TopShot.Admin>() != nil: "Admin capability is invalid"
            }
            self.admin = newAdminCapability
        }

        // The owner can call this function to remove the admin capability
        // from the resource so it can be transferred elsewhere
        pub fun removeAdmin(): Capability {
            pre {
                self.admin != nil: "No Admin to remove"
            }

            let admin = self.admin!
            self.admin = nil

            return admin
        }

        init() {
            self.admin = nil
        }
    }

    // allows anyone to create a new AdminHolder resource
    pub fun createEmptyAdminHolder(): @AdminHolder {
        return <-create AdminHolder()
    }

    init() {

        // Save the adminholder resource to storage
        self.account.save(<-create AdminHolder(), to: /storage/topshotAdminHolder)

        // Publish a Receiver capability to the AdminHolder resource
        self.account.link<&AdminHolder{Receiver}>(/public/topshotAdminReceiver, target: /storage/topshotAdminHolder)
    }

}
