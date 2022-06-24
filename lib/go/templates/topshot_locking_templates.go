package templates

import "github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"

const (
	lockMomentFilename          = "user/lock_moment.cdc"
	unlockMomentFilename        = "user/unlock_moment.cdc"
	batchLockMomentsFilename    = "user/batch_lock_moments.cdc"
	batchUnlockMomentsFilename  = "user/batch_unlock_moments.cdc"
	isLockedScriptFilename      = "collections/get_moment_isLocked.cdc"
	getLockExpiryScriptFilename = "collections/get_moment_lockExpiry.cdc"
	lockFakeNFTFilename         = "user/lock_fake_nft.cdc"

	adminMarkMomentUnlockableFilename = "admin/mark_moment_unlockable.cdc"
)

// GenerateTopShotLockingLockMomentScript creates a script that locks a moment.
func GenerateTopShotLockingLockMomentScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + lockMomentFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateTopShotLockingUnlockMomentScript creates a script that unlocks a moment.
func GenerateTopShotLockingUnlockMomentScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + unlockMomentFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateBatchLockMomentScript creates a script that locks multiple moments.
func GenerateBatchLockMomentScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + batchLockMomentsFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateBatchUnlockMomentScript creates a script that unlocks multiple moments.
func GenerateBatchUnlockMomentScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + batchUnlockMomentsFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetMomentIsLockedScript creates a script that checks if a moment is locked
func GenerateGetMomentIsLockedScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + isLockedScriptFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetMomentLockExpiryScript creates a script that returns the expiry timestamp of a moment
func GenerateGetMomentLockExpiryScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + getLockExpiryScriptFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateLockFakeNFTScript creates a script that tries to lock a NonFungibleToken.NFT
func GenerateLockFakeNFTScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + lockFakeNFTFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateAdminMarkMomentUnlockableScript creates a script that marks a moment as unlockable
func GenerateAdminMarkMomentUnlockableScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + adminMarkMomentUnlockableFilename)

	return []byte(replaceAddresses(code, env))
}
