//
// SPDX-License-Identifier: BSD-3-Clause
//

package redfish

import (
	"encoding/json"
	"strconv"

	"github.com/stmcginnis/gofish/common"
)

// SubProcessor is used to represent a single subprocessor contained within a
// processor.
type SubProcessor struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// MaxSpeedMHz shall indicate the maximum rated clock
	// speed of the processor in MHz.
	MaxSpeedMHz float32
	// ProcessorType shall contain the string which
	// identifies the type of processor contained in this Socket.
	ProcessorType ProcessorType
	// TotalThreads shall indicate the total count of
	// independent execution threads supported by this processor.
	TotalThreads int
	// Status shall contain any status or health properties
	// of the resource.
	Status common.Status
	// Chassis shall be a reference to a
	// resource of type Chassis that represent the physical container
	// associated with this Processor.
	chassis string
	// ConnectedProcessors shall be an array of
	// references of type Processor that are directly connected to this
	// Processor.
	connectedProcessors []string
}

// UnmarshalJSON unmarshals a Processor object from the raw JSON.
func (subProcessor *SubProcessor) UnmarshalJSON(b []byte) error {
	type temp SubProcessor
	type t1 struct {
		temp
		Links struct {
			Chassis             common.Link
			ConnectedProcessors common.Links
		}
	}
	var t t1

	err := json.Unmarshal(b, &t)
	if err != nil {
		// Handle invalid data type returned for MaxSpeedMHz
		var t2 struct {
			t1
			MaxSpeedMHz string
		}
		err2 := json.Unmarshal(b, &t2)

		if err2 != nil {
			// Return the original error
			return err
		}

		// Extract the real Processor struct and replace its MaxSpeedMHz with
		// the parsed string version
		t = t2.t1
		if t2.MaxSpeedMHz != "" {
			bitSize := 32
			mhz, err := strconv.ParseFloat(t2.MaxSpeedMHz, bitSize)
			if err != nil {
				t.MaxSpeedMHz = float32(mhz)
			}
		}
	}

	*subProcessor = SubProcessor(t.temp)

	// Extract the links to other entities for later
	subProcessor.chassis = t.Links.Chassis.String()
	subProcessor.connectedProcessors = t.Links.ConnectedProcessors.ToStrings()

	return nil
}

// GetSubProcessor will get a SubProcessor instance from the processor
func GetSubProcessor(c common.Client, uri string) (*SubProcessor, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var subProcessor SubProcessor
	err = json.NewDecoder(resp.Body).Decode(&subProcessor)
	if err != nil {
		return nil, err
	}

	subProcessor.SetClient(c)
	return &subProcessor, nil
}
