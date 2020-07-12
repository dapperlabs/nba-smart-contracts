package test

// type Client struct {
// 	flowClient     *client.Client
// 	signer         crypto.Signer
// 	serviceAccount *sdk.Account
// 	serviceKey     emulator.ServiceKey
// 	txQueue        []sdk.Transaction
// }

// // NewClient returns a client that conforms to BlockchainAPI
// func NewClient(flowAccessAddress, privateKeyHex string) (BlockchainAPI, error) {
// 	flowClient, err := client.New(flowAccessAddress, grpc.WithInsecure())
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Service address is the first generated account address
// 	addr := sdk.ServiceAddress(sdk.Testnet)

// 	acc, err := flowClient.GetAccount(context.Background(), addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	accountKey := acc.Keys[0]

// 	privateKey, err := crypto.DecodePrivateKeyHex(accountKey.SigAlgo, privateKeyHex)
// 	if err != nil {
// 		return nil, err
// 	}
// 	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)

// 	serviceKey := emulator.ServiceKey{
// 		ID:             accountKey.ID,
// 		Address:        addr,
// 		SequenceNumber: accountKey.SequenceNumber,
// 		PrivateKey:     &privateKey,
// 		PublicKey:      &accountKey.PublicKey,
// 		HashAlgo:       accountKey.HashAlgo,
// 		SigAlgo:        accountKey.SigAlgo,
// 		Weight:         accountKey.Weight,
// 	}

// 	return &Client{
// 		flowClient:     flowClient,
// 		signer:         signer,
// 		serviceAccount: acc,
// 		serviceKey:     serviceKey,
// 		txQueue:        []sdk.Transaction{},
// 	}, nil
// }

// func (c *Client) AddTransaction(tx sdk.Transaction) error {
// 	c.txQueue = append(c.txQueue, tx)
// 	return nil
// }
// func (c *Client) ExecuteNextTransaction() (*types.TransactionResult, error) {
// 	ctx := context.Background()
// 	tx := c.txQueue[0]
// 	c.txQueue = c.txQueue[1:]

// 	err := c.flowClient.SendTransaction(ctx, tx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	txResp := examples.WaitForSeal(ctx, c.flowClient, tx.ID())

// 	// If service account was the proposer, we have to manage the sequence number here
// 	if txResp.Error == nil && // TODO: remove once https://github.com/dapperlabs/flow-go/issues/4107 is done
// 		tx.ProposalKey.Address == c.serviceKey.Address &&
// 		tx.ProposalKey.KeyID == c.serviceKey.ID {
// 		c.serviceKey.SequenceNumber++
// 	}

// 	return &types.TransactionResult{
// 		TransactionID: tx.ID(),
// 		Error:         txResp.Error,
// 		Events:        txResp.Events,
// 	}, nil
// }
// func (c *Client) CreateAccount(publicKeys []*flow.AccountKey, code []byte) (flow.Address, error) {
// 	ctx := context.Background()

// 	for _, key := range publicKeys {
// 		// Reset IDs and Sequence Numbers
// 		key.ID = 0
// 		key.SequenceNumber = 0
// 	}

// 	accountScript := templates.CreateAccount(publicKeys, code)
// 	addr := flow.Address{}

// 	finalizedBlock, err := c.flowClient.GetLatestBlockHeader(ctx, false)
// 	if err != nil {
// 		return addr, err
// 	}

// 	accountTx := flow.NewTransaction().
// 		AddAuthorizer(c.serviceKey.Address).
// 		SetReferenceBlockID(finalizedBlock.ID).
// 		SetScript(accountScript).
// 		SetProposalKey(c.serviceKey.Address, c.serviceKey.ID, c.serviceKey.SequenceNumber).
// 		SetPayer(c.serviceKey.Address)

// 	err = accountTx.SignEnvelope(
// 		c.serviceKey.Address,
// 		c.serviceKey.ID,
// 		c.signer,
// 	)
// 	if err != nil {
// 		return addr, err
// 	}
// 	err = c.flowClient.SendTransaction(ctx, *accountTx)
// 	if err != nil {
// 		return addr, err
// 	}
// 	accountTxResp := examples.WaitForSeal(ctx, c.flowClient, accountTx.ID())
// 	if accountTxResp.Error != nil {
// 		return addr, accountTxResp.Error
// 	}
// 	// Successful Tx, increment sequence number
// 	c.serviceKey.SequenceNumber++

// 	for _, event := range accountTxResp.Events {
// 		if event.Type == flow.EventAccountCreated {
// 			accountCreatedEvent := flow.AccountCreatedEvent(event)
// 			addr = accountCreatedEvent.Address()
// 		}
// 	}
// 	return addr, nil
// }
// func (c *Client) ExecuteBlock() ([]*types.TransactionResult, error) {
// 	panic("not implemented")
// }
// func (c *Client) CommitBlock() (*sdk.Block, error) {
// 	return c.GetLatestBlock()
// }
// func (c *Client) ExecuteAndCommitBlock() (*sdk.Block, []*types.TransactionResult, error) {
// 	panic("not implemented")
// }
// func (c *Client) GetLatestBlock() (*sdk.Block, error) {
// 	ctx := context.Background()
// 	return c.flowClient.GetLatestBlock(ctx, true)
// }
// func (c *Client) GetBlockByID(id sdk.Identifier) (*sdk.Block, error) {
// 	ctx := context.Background()
// 	return c.flowClient.GetBlockByID(ctx, id)
// }
// func (c *Client) GetBlockByHeight(height uint64) (*sdk.Block, error) {
// 	ctx := context.Background()
// 	return c.flowClient.GetBlockByHeight(ctx, height)
// }
// func (c *Client) GetCollection(colID sdk.Identifier) (*sdk.Collection, error) {
// 	panic("not implemented")
// }
// func (c *Client) GetTransaction(txID sdk.Identifier) (*sdk.Transaction, error) {
// 	panic("not implemented")
// }
// func (c *Client) GetTransactionResult(txID sdk.Identifier) (*sdk.TransactionResult, error) {
// 	panic("not implemented")
// }
// func (c *Client) GetAccount(address sdk.Address) (*sdk.Account, error) {
// 	panic("not implemented")
// }
// func (c *Client) GetAccountAtBlock(address sdk.Address, blockHeight uint64) (*sdk.Account, error) {
// 	panic("not implemented")
// }
// func (c *Client) GetEventsByHeight(blockHeight uint64, eventType string) ([]sdk.Event, error) {
// 	panic("not implemented")
// }
// func (c *Client) ExecuteScript(script []byte) (*types.ScriptResult, error) {
// 	ctx := context.Background()
// 	res, err := c.flowClient.ExecuteScriptAtLatestBlock(ctx, script)
// 	return &types.ScriptResult{
// 		Value: res,
// 		Error: err,
// 	}, nil
// }
// func (c *Client) ExecuteScriptAtBlock(script []byte, blockHeight uint64) (*types.ScriptResult, error) {
// 	panic("not implemented")
// }
// func (c *Client) ServiceKey() emulator.ServiceKey {
// 	return c.serviceKey
// }
