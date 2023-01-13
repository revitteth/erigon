// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package types

import (
	"encoding/json"

	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/common/hexutil"
)

var _ = (*withdrawalMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (w Withdrawal) MarshalJSON() ([]byte, error) {
	type Withdrawal struct {
		Index     hexutil.Uint64 `json:"index"`
		Validator hexutil.Uint64 `json:"validatorIndex"`
		Address   common.Address `json:"address"`
		Amount    hexutil.Uint64 `json:"amount"`
	}
	var enc Withdrawal
	enc.Index = hexutil.Uint64(w.Index)
	enc.Validator = hexutil.Uint64(w.Validator)
	enc.Address = w.Address
	enc.Amount = hexutil.Uint64(w.Amount)
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (w *Withdrawal) UnmarshalJSON(input []byte) error {
	type Withdrawal struct {
		Index     *hexutil.Uint64 `json:"index"`
		Validator *hexutil.Uint64 `json:"validatorIndex"`
		Address   *common.Address `json:"address"`
		Amount    *hexutil.Uint64 `json:"amount"`
	}
	var dec Withdrawal
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Index != nil {
		w.Index = uint64(*dec.Index)
	}
	if dec.Validator != nil {
		w.Validator = uint64(*dec.Validator)
	}
	if dec.Address != nil {
		w.Address = *dec.Address
	}
	if dec.Amount != nil {
		w.Amount = uint64(*dec.Amount)
	}
	return nil
}
