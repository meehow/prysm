package sync

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/interfaces"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"google.golang.org/protobuf/proto"
)

func (s *Service) voluntaryExitSubscriber(_ context.Context, msg proto.Message) error {
	ve, ok := msg.(*ethpb.SignedVoluntaryExit)
	if !ok {
		return fmt.Errorf("wrong type, expected: *ethpb.SignedVoluntaryExit got: %T", msg)
	}

	if ve.Exit == nil {
		return errors.New("exit can't be nil")
	}
	s.setExitIndexSeen(ve.Exit.ValidatorIndex)

	s.cfg.exitPool.InsertVoluntaryExit(ve)
	return nil
}

func (s *Service) attesterSlashingSubscriber(ctx context.Context, msg proto.Message) error {
	aSlashing, ok := msg.(interfaces.AttesterSlashing)
	if !ok {
		return fmt.Errorf("wrong type, expected: *ethpb.AttesterSlashing got: %T", msg)
	}
	// Do some nil checks to prevent easy DoS'ing of this handler.
	aSlashing1IsNil := aSlashing == nil || aSlashing.GetFirstAttestation() == nil || aSlashing.GetFirstAttestation().GetAttestingIndices() == nil
	aSlashing2IsNil := aSlashing == nil || aSlashing.GetSecondAttestation() == nil || aSlashing.GetSecondAttestation().GetAttestingIndices() == nil
	if !aSlashing1IsNil && !aSlashing2IsNil {
		headState, err := s.cfg.chain.HeadState(ctx)
		if err != nil {
			return err
		}
		if err := s.cfg.slashingPool.InsertAttesterSlashing(ctx, headState, aSlashing); err != nil {
			return errors.Wrap(err, "could not insert attester slashing into pool")
		}
		s.setAttesterSlashingIndicesSeen(aSlashing.GetFirstAttestation().GetAttestingIndices(), aSlashing.GetSecondAttestation().GetAttestingIndices())
	}
	return nil
}

func (s *Service) proposerSlashingSubscriber(ctx context.Context, msg proto.Message) error {
	pSlashing, ok := msg.(*ethpb.ProposerSlashing)
	if !ok {
		return fmt.Errorf("wrong type, expected: *ethpb.ProposerSlashing got: %T", msg)
	}
	// Do some nil checks to prevent easy DoS'ing of this handler.
	header1IsNil := pSlashing == nil || pSlashing.Header_1 == nil || pSlashing.Header_1.Header == nil
	header2IsNil := pSlashing == nil || pSlashing.Header_2 == nil || pSlashing.Header_2.Header == nil
	if !header1IsNil && !header2IsNil {
		headState, err := s.cfg.chain.HeadState(ctx)
		if err != nil {
			return err
		}
		if err := s.cfg.slashingPool.InsertProposerSlashing(ctx, headState, pSlashing); err != nil {
			return errors.Wrap(err, "could not insert proposer slashing into pool")
		}
		s.setProposerSlashingIndexSeen(pSlashing.Header_1.Header.ProposerIndex)
	}
	return nil
}
