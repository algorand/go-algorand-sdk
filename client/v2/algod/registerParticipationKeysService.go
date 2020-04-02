package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type RegisterParticipationKeysService struct {
	c       *Client
	account string
	p       models.RegisterParticipationKeysAccountIdParams
}

func (s *RegisterParticipationKeysService) Fee(fee uint64) *RegisterParticipationKeysService {
	s.p.Fee = fee
	return s
}

func (s *RegisterParticipationKeysService) KeyDilution(dilution uint64) *RegisterParticipationKeysService {
	s.p.KeyDilution = dilution
	return s
}

func (s *RegisterParticipationKeysService) RoundLastValid(round uint64) *RegisterParticipationKeysService {
	s.p.RoundLastValid = round
	return s
}

func (s *RegisterParticipationKeysService) NoWait(nowait bool) *RegisterParticipationKeysService {
	s.p.NoWait = nowait
	return s
}

func (s *RegisterParticipationKeysService) Do(ctx context.Context, headers ...*common.Header) error {
	return s.c.post(ctx, nil, fmt.Sprintf("/register-participation-keys/%s", s.account), s.p, headers)
}
