package kv

import (
	"sort"
	"testing"

	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/interfaces"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v5/testing/assert"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
	"github.com/prysmaticlabs/prysm/v5/testing/util"
)

func TestKV_Forkchoice_CanSaveRetrieve(t *testing.T) {
	cache := NewAttCaches()

	att1 := util.HydrateAttestation(&ethpb.Attestation{Data: &ethpb.AttestationData{Slot: 1}, AggregationBits: bitfield.Bitlist{0b1101}})
	att2 := util.HydrateAttestation(&ethpb.Attestation{Data: &ethpb.AttestationData{Slot: 2}, AggregationBits: bitfield.Bitlist{0b1101}})
	att3 := util.HydrateAttestation(&ethpb.Attestation{Data: &ethpb.AttestationData{Slot: 3}, AggregationBits: bitfield.Bitlist{0b1101}})
	atts := []interfaces.Attestation{att1, att2, att3}

	for _, att := range atts {
		require.NoError(t, cache.SaveForkchoiceAttestation(att))
	}

	returned := cache.ForkchoiceAttestations()

	sort.Slice(returned, func(i, j int) bool {
		return returned[i].GetData().Slot < returned[j].GetData().Slot
	})

	assert.DeepEqual(t, atts, returned)
}

func TestKV_Forkchoice_CanDelete(t *testing.T) {
	cache := NewAttCaches()

	att1 := util.HydrateAttestation(&ethpb.Attestation{Data: &ethpb.AttestationData{Slot: 1}, AggregationBits: bitfield.Bitlist{0b1101}})
	att2 := util.HydrateAttestation(&ethpb.Attestation{Data: &ethpb.AttestationData{Slot: 2}, AggregationBits: bitfield.Bitlist{0b1101}})
	att3 := util.HydrateAttestation(&ethpb.Attestation{Data: &ethpb.AttestationData{Slot: 3}, AggregationBits: bitfield.Bitlist{0b1101}})
	atts := []interfaces.Attestation{att1, att2, att3}

	for _, att := range atts {
		require.NoError(t, cache.SaveForkchoiceAttestation(att))
	}

	require.NoError(t, cache.DeleteForkchoiceAttestation(att1))
	require.NoError(t, cache.DeleteForkchoiceAttestation(att3))

	returned := cache.ForkchoiceAttestations()
	wanted := []interfaces.Attestation{att2}
	assert.DeepEqual(t, wanted, returned)
}

func TestKV_Forkchoice_CanCount(t *testing.T) {
	cache := NewAttCaches()

	att1 := util.HydrateAttestation(&ethpb.Attestation{Data: &ethpb.AttestationData{Slot: 1}, AggregationBits: bitfield.Bitlist{0b1101}})
	att2 := util.HydrateAttestation(&ethpb.Attestation{Data: &ethpb.AttestationData{Slot: 2}, AggregationBits: bitfield.Bitlist{0b1101}})
	att3 := util.HydrateAttestation(&ethpb.Attestation{Data: &ethpb.AttestationData{Slot: 3}, AggregationBits: bitfield.Bitlist{0b1101}})
	atts := []*ethpb.Attestation{att1, att2, att3}

	for _, att := range atts {
		require.NoError(t, cache.SaveForkchoiceAttestation(att))
	}

	require.Equal(t, 3, cache.ForkchoiceAttestationCount())
}
