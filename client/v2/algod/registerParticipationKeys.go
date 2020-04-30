package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type RegisterParticipationKeys struct {
	c       *Client
	account string
	p       models.RegisterParticipationKeysAccountIdParams
}

func (s *RegisterParticipationKeys) Fee(fee uint64) *RegisterParticipationKeys {
	s.p.Fee = fee
	return s
}

func (s *RegisterParticipationKeys) KeyDilution(dilution uint64) *RegisterParticipationKeys {
	s.p.KeyDilution = dilution
	return s
}

func (s *RegisterParticipationKeys) RoundLastValid(round uint64) *RegisterParticipationKeys {
	s.p.RoundLastValid = round
	return s
}

func (s *RegisterParticipationKeys) NoWait(nowait bool) *RegisterParticipationKeys {
	s.p.NoWait = nowait
	return s
}

func (s *RegisterParticipationKeys) Do(ctx context.Context, headers ...*common.Header) error {
	return s.c.post(ctx, nil, fmt.Sprintf("/v2/register-participation-keys/%s", s.account), s.p, headers)
}
