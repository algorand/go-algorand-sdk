package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ConfigurableConsensusProtocolsFilename defines a set of consensus protocols that
// are to be loaded from the data directory ( if present ), to override the
// built-in supported consensus protocols.
const ConfigurableConsensusProtocolsFilename = "consensus.json"

// SaveConfigurableConsensus saves the configurable protocols file to the provided data directory.
// if the params contains zero protocols, the existing consensus.json file will be removed if exists.
func SaveConfigurableConsensus(dataDirectory string, params ConsensusProtocols) error {
	consensusProtocolPath := filepath.Join(dataDirectory, ConfigurableConsensusProtocolsFilename)

	if len(params) == 0 {
		// we have no consensus params to write. In this case, just delete the existing file
		// ( if any )
		err := os.Remove(consensusProtocolPath)
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	encodedConsensusParams, err := json.Marshal(params)
	if err != nil {
		return err
	}
	err = os.WriteFile(consensusProtocolPath, encodedConsensusParams, 0644)
	return err
}

// PreloadConfigurableConsensusProtocols loads the configurable protocols from the data directory
// and merge it with a copy of the Consensus map. Then, it returns it to the caller.
func PreloadConfigurableConsensusProtocols(dataDirectory string) (ConsensusProtocols, error) {
	consensusProtocolPath := filepath.Join(dataDirectory, ConfigurableConsensusProtocolsFilename)
	file, err := os.Open(consensusProtocolPath)

	if err != nil {
		if os.IsNotExist(err) {
			// this file is not required, only optional. if it's missing, no harm is done.
			return Consensus, nil
		}
		return nil, err
	}
	defer file.Close()

	configurableConsensus := make(ConsensusProtocols)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configurableConsensus)
	if err != nil {
		return nil, err
	}
	return Consensus.Merge(configurableConsensus), nil
}

// SetConfigurableConsensusProtocols sets the configurable protocols.
func SetConfigurableConsensusProtocols(newConsensus ConsensusProtocols) ConsensusProtocols {
	oldConsensus := Consensus
	Consensus = newConsensus
	// Set allocation limits
	// checkSetAllocBounds not ported to sdk https://github.com/algorand/go-algorand/blob/e68b54e90cd9dc1b52c3a9df85e0aeb56e8206d5/config/consensus.go#L723
	// for _, p := range Consensus {
	// 	checkSetAllocBounds(p)
	// }
	return oldConsensus
}

// LoadConfigurableConsensusProtocols loads the configurable protocols from the data directory
func LoadConfigurableConsensusProtocols(dataDirectory string) error {
	newConsensus, err := PreloadConfigurableConsensusProtocols(dataDirectory)
	if err != nil {
		return err
	}
	if newConsensus != nil {
		SetConfigurableConsensusProtocols(newConsensus)
	}
	return nil
}
