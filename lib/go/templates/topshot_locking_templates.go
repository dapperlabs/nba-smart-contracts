package templates

import "github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"

const (
	lockMomentFilename          = "user/lock_moment.cdc"
	unlockMomentFilename        = "user/unlock_moment.cdc"
	isLockedScriptFilename      = "collections/get_moment_isLocked.cdc"
	getLockExpiryScriptFilename = "collections/get_moment_lockExpiry.cdc"

	adminMarkMomentUnlockableFilename = "admin/mark_moment_unlockable.cdc"
)

// GenerateGetTopShotLockingLockMomentScript creates a script that locks a moment.
func GenerateGetTopShotLockingLockMomentScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + lockMomentFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetTopShotLockingUnlockMomentScript creates a script that locks a moment.
func GenerateGetTopShotLockingUnlockMomentScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + unlockMomentFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetMomentIsLockedScript creates a script that checks if a moment is locked
func GenerateGetMomentIsLockedScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + isLockedScriptFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetMomentIsLockExpiryScript creates a script that checks if a moment is locked
func GenerateGetMomentIsLockExpiryScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + getLockExpiryScriptFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetMomentIsLockExpiryScript creates a script that checks if a moment is locked
func GenerateAdminMarkMomentUnlockableScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + adminMarkMomentUnlockableFilename)

	return []byte(replaceAddresses(code, env))
}
