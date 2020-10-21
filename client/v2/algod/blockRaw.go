package algod

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type BlockRaw struct {
	c     *Client
	round uint64
	p     models.GetBlockParams
}

/*
func (s *BlockRaw) Do(ctx context.Context, headers ...*common.Header) (result []byte, err error) {
	s.p.Format = "msgpack"
	err = s.c.getMsgpack(ctx, &result, fmt.Sprintf("/v2/blocks/%d", s.round), s.p, headers)
	if err != nil {
		return
	}
	return
}
*/

func (s *BlockRaw) Do(ctx context.Context, headers ...*common.Header) (result []byte, err error) {
	s.p.Format = "msgpack"
	common := common.Client(*s.c)

	var response *http.Response
	response, err = common.SubmitFormRaw(ctx, fmt.Sprintf("/v2/blocks/%d", s.round),s.p, "GET", false, headers)
	if err != nil {
		response.Body.Close()
		return
	}

	result, err = ioutil.ReadAll(response.Body)
	if err != nil {
		response.Body.Close()
		return
	}

	response.Body.Close()
	return
}