package nba

import _ "embed"

var (
	// shardedCollection
	//go:embed transactions/shardedCollection/setup_sharded_collection.cdc
	ShardedcollectionSetupShardedCollection []byte

	//go:embed transactions/shardedCollection/transfer_from_sharded.cdc
	ShardedcollectionTransferFromSharded []byte

	//go:embed transactions/shardedCollection/batch_from_sharded.cdc
	ShardedcollectionBatchFromSharded []byte

	// admin
	//go:embed transactions/admin/add_plays_to_set.cdc
	AdminAddPlaysToSet []byte

	//go:embed transactions/admin/fulfill_pack.cdc
	AdminFulfillPack []byte

	//go:embed transactions/admin/retireAll_plays_from_set.cdc
	AdminRetireallPlaysFromSet []byte

	//go:embed transactions/admin/mark_moment_unlockable.cdc
	AdminMarkMomentUnlockable []byte

	//go:embed transactions/admin/create_subedition.cdc
	AdminCreateSubedition []byte

	//go:embed transactions/admin/transfer_admin.cdc
	AdminTransferAdmin []byte

	//go:embed transactions/admin/mint_moment.cdc
	AdminMintMoment []byte

	//go:embed transactions/admin/grant_topshot_locking_admin.cdc
	AdminGrantTopshotLockingAdmin []byte

	//go:embed transactions/admin/batch_mint_moment.cdc
	AdminBatchMintMoment []byte

	//go:embed transactions/admin/set_nft_subedition.cdc
	AdminSetNftSubedition []byte

	//go:embed transactions/admin/create_set_and_play_struct.cdc
	AdminCreateSetAndPlayStruct []byte

	//go:embed transactions/admin/create_new_subedition_admin_resource.cdc
	AdminCreateNewSubeditionAdminResource []byte

	//go:embed transactions/admin/add_play_to_set.cdc
	AdminAddPlayToSet []byte

	//go:embed transactions/admin/set_nfts_lock_expiry.cdc
	AdminSetNftsLockExpiry []byte

	//go:embed transactions/admin/start_new_series.cdc
	AdminStartNewSeries []byte

	//go:embed transactions/admin/update_tagline.cdc
	AdminUpdateTagline []byte

	//go:embed transactions/admin/lock_set.cdc
	AdminLockSet []byte

	//go:embed transactions/admin/batch_mint_moment_with_subedition.cdc
	AdminBatchMintMomentWithSubedition []byte

	//go:embed transactions/admin/retire_play_from_set.cdc
	AdminRetirePlayFromSet []byte

	//go:embed transactions/admin/mint_moment_with_subedition.cdc
	AdminMintMomentWithSubedition []byte

	//go:embed transactions/admin/create_play.cdc
	AdminCreatePlay []byte

	//go:embed transactions/admin/retire_all.cdc
	AdminRetireAll []byte

	//go:embed transactions/admin/create_set.cdc
	AdminCreateSet []byte

	//go:embed transactions/admin/unlock_all_moments.cdc
	AdminUnlockAllMoments []byte

	// marketV3
	//go:embed transactions/marketV3/purchase_both_markets.cdc
	Marketv3PurchaseBothMarkets []byte

	//go:embed transactions/marketV3/purchase_moment.cdc
	Marketv3PurchaseMoment []byte

	//go:embed transactions/marketV3/change_receiver.cdc
	Marketv3ChangeReceiver []byte

	//go:embed transactions/marketV3/create_sale.cdc
	Marketv3CreateSale []byte

	//go:embed transactions/marketV3/change_price.cdc
	Marketv3ChangePrice []byte

	//go:embed transactions/marketV3/stop_sale.cdc
	Marketv3StopSale []byte

	//go:embed transactions/marketV3/mint_and_purchase.cdc
	Marketv3MintAndPurchase []byte

	//go:embed transactions/marketV3/upgrade_sale.cdc
	Marketv3UpgradeSale []byte

	//go:embed transactions/marketV3/purchase_group_of_moments.cdc
	Marketv3PurchaseGroupOfMoments []byte

	// marketV3/scripts
	//go:embed transactions/marketV3/scripts/get_sale_percentage.cdc
	Marketv3ScriptsGetSalePercentage []byte

	//go:embed transactions/marketV3/scripts/get_sale_set_id.cdc
	Marketv3ScriptsGetSaleSetId []byte

	//go:embed transactions/marketV3/scripts/get_sale_price.cdc
	Marketv3ScriptsGetSalePrice []byte

	//go:embed transactions/marketV3/scripts/get_sale_len.cdc
	Marketv3ScriptsGetSaleLen []byte

	// marketV3
	//go:embed transactions/marketV3/start_sale.cdc
	Marketv3StartSale []byte

	//go:embed transactions/marketV3/create_start_sale.cdc
	Marketv3CreateStartSale []byte

	// user
	//go:embed transactions/user/transfer_moment.cdc
	UserTransferMoment []byte

	//go:embed transactions/user/setup_collection.cdc
	UserSetupCollection []byte

	//go:embed transactions/user/transfer_moment_v3_sale.cdc
	UserTransferMomentV3Sale []byte

	//go:embed transactions/user/batch_transfer.cdc
	UserBatchTransfer []byte

	//go:embed transactions/user/lock_fake_nft.cdc
	UserLockFakeNft []byte

	//go:embed transactions/user/setup_up_all_collections.cdc
	UserSetupUpAllCollections []byte

	//go:embed transactions/user/unlock_moment.cdc
	UserUnlockMoment []byte

	//go:embed transactions/user/setup_switchboard_account.cdc
	UserSetupSwitchboardAccount []byte

	//go:embed transactions/user/lock_moment.cdc
	UserLockMoment []byte

	//go:embed transactions/user/destroy_moments.cdc
	UserDestroyMoments []byte

	//go:embed transactions/user/destroy_moments_v2.cdc
	UserDestroyMomentsV2 []byte

	//go:embed transactions/user/batch_unlock_moments.cdc
	UserBatchUnlockMoments []byte

	//go:embed transactions/user/batch_lock_moments.cdc
	UserBatchLockMoments []byte

	// scripts
	//go:embed transactions/scripts/get_nft_metadata.cdc
	ScriptsGetNftMetadata []byte

	//go:embed transactions/scripts/get_currentSeries.cdc
	ScriptsGetCurrentseries []byte

	// scripts/plays
	//go:embed transactions/scripts/plays/get_play_metadata.cdc
	ScriptsPlaysGetPlayMetadata []byte

	//go:embed transactions/scripts/plays/get_play_metadata_field.cdc
	ScriptsPlaysGetPlayMetadataField []byte

	//go:embed transactions/scripts/plays/get_all_plays.cdc
	ScriptsPlaysGetAllPlays []byte

	//go:embed transactions/scripts/plays/get_nextPlayID.cdc
	ScriptsPlaysGetNextplayid []byte

	// scripts/sets
	//go:embed transactions/scripts/sets/get_setName.cdc
	ScriptsSetsGetSetname []byte

	//go:embed transactions/scripts/sets/get_set_locked.cdc
	ScriptsSetsGetSetLocked []byte

	//go:embed transactions/scripts/sets/get_edition_retired.cdc
	ScriptsSetsGetEditionRetired []byte

	//go:embed transactions/scripts/sets/get_set_data.cdc
	ScriptsSetsGetSetData []byte

	//go:embed transactions/scripts/sets/get_setIDs_by_name.cdc
	ScriptsSetsGetSetidsByName []byte

	//go:embed transactions/scripts/sets/get_numMoments_in_edition.cdc
	ScriptsSetsGetNummomentsInEdition []byte

	//go:embed transactions/scripts/sets/get_plays_in_set.cdc
	ScriptsSetsGetPlaysInSet []byte

	//go:embed transactions/scripts/sets/get_setSeries.cdc
	ScriptsSetsGetSetseries []byte

	//go:embed transactions/scripts/sets/get_nextSetID.cdc
	ScriptsSetsGetNextsetid []byte

	// scripts/users
	//go:embed transactions/scripts/users/is_account_all_set_up.cdc
	ScriptsUsersIsAccountAllSetUp []byte

	// scripts/subeditions
	//go:embed transactions/scripts/subeditions/get_subedition_by_id.cdc
	ScriptsSubeditionsGetSubeditionById []byte

	//go:embed transactions/scripts/subeditions/get_all_subeditions.cdc
	ScriptsSubeditionsGetAllSubeditions []byte

	//go:embed transactions/scripts/subeditions/get_nft_subedition.cdc
	ScriptsSubeditionsGetNftSubedition []byte

	//go:embed transactions/scripts/subeditions/get_nextSubeditionID.cdc
	ScriptsSubeditionsGetNextsubeditionid []byte

	// scripts/collections
	//go:embed transactions/scripts/collections/get_moment_isLocked.cdc
	ScriptsCollectionsGetMomentIslocked []byte

	//go:embed transactions/scripts/collections/get_moment_setName.cdc
	ScriptsCollectionsGetMomentSetname []byte

	//go:embed transactions/scripts/collections/get_moment_playID.cdc
	ScriptsCollectionsGetMomentPlayid []byte

	//go:embed transactions/scripts/collections/get_metadata_field.cdc
	ScriptsCollectionsGetMetadataField []byte

	//go:embed transactions/scripts/collections/get_locked_nfts_length.cdc
	ScriptsCollectionsGetLockedNftsLength []byte

	//go:embed transactions/scripts/collections/get_id_in_Collection.cdc
	ScriptsCollectionsGetIdInCollection []byte

	//go:embed transactions/scripts/collections/get_collection_ids.cdc
	ScriptsCollectionsGetCollectionIds []byte

	//go:embed transactions/scripts/collections/get_moment_lockExpiry.cdc
	ScriptsCollectionsGetMomentLockexpiry []byte

	//go:embed transactions/scripts/collections/get_moment_series.cdc
	ScriptsCollectionsGetMomentSeries []byte

	//go:embed transactions/scripts/collections/get_setplays_are_owned.cdc
	ScriptsCollectionsGetSetplaysAreOwned []byte

	//go:embed transactions/scripts/collections/borrow_nft_safe.cdc
	ScriptsCollectionsBorrowNftSafe []byte

	//go:embed transactions/scripts/collections/get_moment_setID.cdc
	ScriptsCollectionsGetMomentSetid []byte

	//go:embed transactions/scripts/collections/get_metadata.cdc
	ScriptsCollectionsGetMetadata []byte

	//go:embed transactions/scripts/collections/get_moment_serialNum.cdc
	ScriptsCollectionsGetMomentSerialnum []byte

	// scripts
	//go:embed transactions/scripts/get_totalSupply.cdc
	ScriptsGetTotalsupply []byte

	//go:embed transactions/scripts/get_topshot_metadata.cdc
	ScriptsGetTopshotMetadata []byte

	// fastbreak/oracle
	//go:embed transactions/fastbreak/oracle/update_fast_break_game.cdc
	FastbreakOracleUpdateFastBreakGame []byte

	//go:embed transactions/fastbreak/oracle/create_game.cdc
	FastbreakOracleCreateGame []byte

	//go:embed transactions/fastbreak/oracle/add_stat_to_game.cdc
	FastbreakOracleAddStatToGame []byte

	//go:embed transactions/fastbreak/oracle/score_fast_break_submission.cdc
	FastbreakOracleScoreFastBreakSubmission []byte

	//go:embed transactions/fastbreak/oracle/create_run.cdc
	FastbreakOracleCreateRun []byte

	// fastbreak/scripts
	//go:embed transactions/fastbreak/scripts/get_fast_break_stats.cdc
	FastbreakScriptsGetFastBreakStats []byte

	//go:embed transactions/fastbreak/scripts/get_token_count.cdc
	FastbreakScriptsGetTokenCount []byte

	//go:embed transactions/fastbreak/scripts/get_player_score.cdc
	FastbreakScriptsGetPlayerScore []byte

	//go:embed transactions/fastbreak/scripts/get_current_player.cdc
	FastbreakScriptsGetCurrentPlayer []byte

	//go:embed transactions/fastbreak/scripts/get_fast_break.cdc
	FastbreakScriptsGetFastBreak []byte

	// fastbreak/player
	//go:embed transactions/fastbreak/player/update_submission.cdc
	FastbreakPlayerUpdateSubmission []byte

	//go:embed transactions/fastbreak/player/create_player.cdc
	FastbreakPlayerCreatePlayer []byte

	//go:embed transactions/fastbreak/player/play.cdc
	FastbreakPlayerPlay []byte

	// market
	//go:embed transactions/market/purchase_moment.cdc
	MarketPurchaseMoment []byte

	//go:embed transactions/market/change_receiver.cdc
	MarketChangeReceiver []byte

	//go:embed transactions/market/create_sale.cdc
	MarketCreateSale []byte

	//go:embed transactions/market/change_price.cdc
	MarketChangePrice []byte

	//go:embed transactions/market/stop_sale.cdc
	MarketStopSale []byte

	//go:embed transactions/market/change_percentage.cdc
	MarketChangePercentage []byte

	//go:embed transactions/market/mint_and_purchase.cdc
	MarketMintAndPurchase []byte

	// market/scripts
	//go:embed transactions/market/scripts/get_sale_percentage.cdc
	MarketScriptsGetSalePercentage []byte

	//go:embed transactions/market/scripts/get_sale_set_id.cdc
	MarketScriptsGetSaleSetId []byte

	//go:embed transactions/market/scripts/get_sale_price.cdc
	MarketScriptsGetSalePrice []byte

	//go:embed transactions/market/scripts/get_sale_len.cdc
	MarketScriptsGetSaleLen []byte

	// market
	//go:embed transactions/market/start_sale.cdc
	MarketStartSale []byte

	//go:embed transactions/market/create_start_sale.cdc
	MarketCreateStartSale []byte
)
