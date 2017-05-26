package handler

import (
	"burrow/word256"
	"fmt"
	"github.com/golang/protobuf/proto"
	. "sawtooth_burrow_evm/protobuf/evm_pb2"
	"sawtooth_sdk/client"
	"sawtooth_sdk/processor"
)

// -- AppState --

// StateManager simplifies accessing EVM related data stored in state
type StateManager struct {
	state  *processor.State
	prefix string
}

func NewStateManager(prefix string, state *processor.State) *StateManager {
	return &StateManager{
		prefix: prefix,
		state:  state,
	}
}

// NewEntry creates a new entry in state. If an entry already exists at the
// given address or the entry cannot be created, an error is returned.
func (mgr *StateManager) NewEntry(vmAddress []byte) (*EvmEntry, error) {
	entry, err := mgr.GetEntry(vmAddress)
	if err != nil {
		return nil, err
	}

	if entry != nil {
		return nil, fmt.Errorf("Address already in use")
	}

	entry = &EvmEntry{
		Account: &EvmStateAccount{
			Address: vmAddress,
			Balance: 0,
			Code:    make([]byte, 0),
			Nonce:   0,
		},
		Storage: make([]*EvmStorage, 0),
	}

	err = mgr.SetEntry(vmAddress, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// DelEntry removes the given entry from state. An error is returned if the
// entry does not exist.
func (mgr *StateManager) DelEntry(vmAddress []byte) error {
	entry, err := mgr.GetEntry(vmAddress)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("Entry does not exist %v", vmAddress)
	}
	err = mgr.SetEntry(vmAddress, &EvmEntry{})
	if err != nil {
		return err
	}
	return nil
}

// GetEntry retrieve the entry from state at the given address. If the entry
// does not exist, nil is returned.
func (mgr *StateManager) GetEntry(vmAddress []byte) (*EvmEntry, error) {
	address := toStateAddress(mgr.prefix, vmAddress)

	// Retrieve the account from global state
	entries, err := mgr.state.Get([]string{address})
	if err != nil {
		return nil, err
	}
	entryData, exists := entries[address]
	if !exists {
		return nil, nil
	}

	// Deserialize the entry
	entry := &EvmEntry{}
	err = proto.Unmarshal(entryData, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// MustGetEntry wraps GetEntry and panics if the entry does not exist of there
// is an error.
func (mgr *StateManager) MustGetEntry(vmAddress []byte) *EvmEntry {
	entry, err := mgr.GetEntry(vmAddress)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to GetEntry(%v): %v", vmAddress, err,
		))
	}

	if entry == nil {
		panic(fmt.Sprintf(
			"Tried to GetEntry(%v) but nothing exists there", vmAddress,
		))
	}

	return entry
}

// SetEntry writes the entry to the given address. Returns an error if it fails
// to set the address.
func (mgr *StateManager) SetEntry(vmAddress []byte, entry *EvmEntry) error {
	address := toStateAddress(mgr.prefix, vmAddress)

	entryData, err := proto.Marshal(entry)
	if err != nil {
		return err
	}

	// Store the account in global state
	addresses, err := mgr.state.Set(map[string][]byte{
		address: entryData,
	})
	if err != nil {
		return err
	}

	for _, a := range addresses {
		if a == address {
			return nil
		}
	}
	return fmt.Errorf("Address not set: %v", address)
}

// MustSetEntry wraps set entry and panics if there is an error.
func (mgr *StateManager) MustSetEntry(vmAddress []byte, entry *EvmEntry) {
	err := mgr.SetEntry(vmAddress, entry)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to SetEntry(%v, %v): %v", vmAddress, entry, err,
		))
	}
}

// Convert the byte representation of an address used by the EVM to the string
// representation used by global state.
func toStateAddress(prefix string, b []byte) string {
	// Make sure the address is padded correctly
	b = word256.RightPadBytes(b, 32)

	return prefix + client.MustEncode(b)
}
