package smartcontracts

import (
	"encoding/json"
	"fmt"

	"coffee-supply-chain-hlf/chaincode/models"
	"coffee-supply-chain-hlf/chaincode/utils"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type FarmerContract struct {
	contractapi.Contract
}

func (tc *FarmerContract) GrowBatch(ctx contractapi.TransactionContextInterface, batchID string) (bool, error) {
	txnID := ctx.GetStub().GetTxID()

	utils.LogMessage("FarmerContract.GrowBatch", "Recieved transaction to store a grown batch", batchID, txnID)

	if batchID == "" {
		return false, fmt.Errorf("FarmerContract.GrowBatch: Batch ID is empty")
	}

	batchStr, err := utils.GetBatchFromTransient(ctx)
	if err != nil {
		return false, fmt.Errorf("error while getting batch data from transient: %s", err)
	}

	var batch models.FarmerBatch
	err = json.Unmarshal(batchStr, &batch)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %s", err)
	}

	err = utils.PutState(ctx, batchID, batch)
	if err != nil {
		return false, fmt.Errorf("error while creating batch. BatchId: %s, error: %w", batchID, err)
	}

	batchGrownEvent := models.BatchGrown{
		BatchID: batchID,
		TxnID:   txnID,
	}

	err = utils.SetEvent(ctx, utils.BATCH_GROWN_EVENT, batchGrownEvent)
	if err != nil {
		return false, fmt.Errorf("error while setting event, %w", err)
	}

	utils.LogMessage("FarmerContract.GrowBatch", "Stored batch data on ledger", batchID, ctx.GetStub().GetTxID())

	return true, nil
}
