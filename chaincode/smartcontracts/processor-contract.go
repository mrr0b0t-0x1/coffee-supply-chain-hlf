package smartcontracts

import (
	"encoding/json"
	"fmt"

	"coffee-supply-chain-hlf/chaincode/models"
	"coffee-supply-chain-hlf/chaincode/utils"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ProcessorContract struct {
	contractapi.Contract
}

func (pc *ProcessorContract) ProcessBatch(ctx contractapi.TransactionContextInterface, batchID string) (bool, error) {
	txnID := ctx.GetStub().GetTxID()

	utils.LogMessage("ProcessorContract.ProcessBatch", "Recieved transaction to proces a batch", batchID, txnID)

	if batchID == "" {
		return false, fmt.Errorf("ProcessorContract.ProcessBatch: Batch ID is empty")
	}

	batchStr, err := utils.GetBatchFromTransient(ctx)
	if err != nil {
		return false, fmt.Errorf("error while getting batch data from transient: %s", err)
	}

	var batch models.ProcessorBatch
	err = json.Unmarshal(batchStr, &batch)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %s", err)
	}

	err = utils.PutState(ctx, batchID, batch)
	if err != nil {
		return false, fmt.Errorf("error while creating batch. BatchId: %s, error: %w", batchID, err)
	}

	batchProcessedEvent := models.Batch{
		BatchID: batchID,
		TxnID:   txnID,
		Org:     utils.PROCESSOR_ROLE,
	}

	err = utils.SetEvent(ctx, utils.BATCH_PROCESSED_EVENT, batchProcessedEvent)
	if err != nil {
		return false, fmt.Errorf("error while setting event, %w", err)
	}

	utils.LogMessage("ProcessorContract.ProcessBatch", "Stored batch processed data on ledger", batchID, ctx.GetStub().GetTxID())

	return true, nil
}
