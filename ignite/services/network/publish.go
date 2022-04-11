package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/ignite-hq/cli/ignite/pkg/cosmoserror"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosutil/genesis"
	"github.com/ignite-hq/cli/ignite/pkg/events"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

// publishOptions holds info about how to create a chain.
type publishOptions struct {
	genesis     *genesis.Genesis
	chainID     string
	campaignID  uint64
	noCheck     bool
	metadata    string
	totalSupply sdk.Coins
	mainnet     bool
}

// PublishOption configures chain creation.
type PublishOption func(*publishOptions)

// WithCampaign add a campaign id.
func WithCampaign(id uint64) PublishOption {
	return func(o *publishOptions) {
		o.campaignID = id
	}
}

// WithChainID use a custom chain id.
func WithChainID(chainID string) PublishOption {
	return func(o *publishOptions) {
		o.chainID = chainID
	}
}

// WithNoCheck disables checking integrity of the chain.
func WithNoCheck() PublishOption {
	return func(o *publishOptions) {
		o.noCheck = true
	}
}

// WithCustomGenesis enables using a custom genesis during publish.
func WithCustomGenesis(gen *genesis.Genesis) PublishOption {
	return func(o *publishOptions) {
		o.genesis = gen
	}
}

// WithMetadata provides a meta data proposal to update the campaign.
func WithMetadata(metadata string) PublishOption {
	return func(c *publishOptions) {
		c.metadata = metadata
	}
}

// WithTotalSupply provides a total supply proposal to update the campaign.
func WithTotalSupply(totalSupply sdk.Coins) PublishOption {
	return func(c *publishOptions) {
		c.totalSupply = totalSupply
	}
}

// Mainnet initialize a published chain into the mainnet
func Mainnet() PublishOption {
	return func(o *publishOptions) {
		o.mainnet = true
	}
}

// Publish submits Genesis to SPN to announce a new network.
func (n Network) Publish(ctx context.Context, c Chain, options ...PublishOption) (launchID, campaignID, mainnetID uint64, err error) {
	o := publishOptions{}
	for _, apply := range options {
		apply(&o)
	}

	var (
		chainID     string
		genesisHash string
		genesisURL  string
	)

	// if the initial genesis is a genesis URL and no check are performed, we simply fetch it and get its hash.
	if o.genesis != nil {
		genesisURL = o.genesis.URL()
		genesisHash, err = o.genesis.Hash()
		if err != nil {
			return 0, 0, 0, err
		}

		n.ev.Send(events.New(events.StatusOngoing, "Fetching custom Genesis from URL"))
		if o.genesis.TarballPath() != "" {
			n.ev.Send(
				events.New(events.StatusDone,
					fmt.Sprintf("Extracted custom Genesis from tarball at %s", o.genesis.TarballPath()),
				),
			)
		} else {
			n.ev.Send(events.New(events.StatusDone, "Custom Genesis JSON from URL fetched"))
		}
		chainID, err = o.genesis.ChainID()
		if err != nil {
			return 0, 0, 0, err
		}
	}

	// use chain id flag always in the highest priority.
	if o.chainID != "" {
		chainID = o.chainID
	}
	// if the chain id is empty, use a default one.
	if chainID == "" {
		chainID, err = c.ChainID()
		if err != nil {
			return 0, 0, 0, err
		}
	}

	coordinatorAddress := n.account.Address(networktypes.SPN)
	campaignID = o.campaignID

	n.ev.Send(events.New(events.StatusOngoing, "Publishing the network"))

	_, err = profiletypes.
		NewQueryClient(n.cosmos.Context).
		CoordinatorByAddress(ctx, &profiletypes.QueryGetCoordinatorByAddressRequest{
			Address: coordinatorAddress,
		})
	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		msgCreateCoordinator := profiletypes.NewMsgCreateCoordinator(
			coordinatorAddress,
			"",
			"",
			"",
		)
		if _, err := n.cosmos.BroadcastTx(n.account.Name, msgCreateCoordinator); err != nil {
			return 0, 0, 0, err
		}
	} else if err != nil {
		return 0, 0, 0, err
	}

	if campaignID != 0 {
		_, err = campaigntypes.
			NewQueryClient(n.cosmos.Context).
			Campaign(ctx, &campaigntypes.QueryGetCampaignRequest{
				CampaignID: o.campaignID,
			})
		if err != nil {
			return 0, 0, 0, err
		}
	} else {
		campaignID, err = n.CreateCampaign(c.Name(), o.metadata, o.totalSupply)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	msgCreateChain := launchtypes.NewMsgCreateChain(
		n.account.Address(networktypes.SPN),
		chainID,
		c.SourceURL(),
		c.SourceHash(),
		genesisURL,
		genesisHash,
		true,
		campaignID,
		nil,
	)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msgCreateChain)
	if err != nil {
		return 0, 0, 0, err
	}

	var createChainRes launchtypes.MsgCreateChainResponse
	if err := res.Decode(&createChainRes); err != nil {
		return 0, 0, 0, err
	}

	if err := c.CacheBinary(createChainRes.LaunchID); err != nil {
		return 0, 0, 0, err

	}

	if o.mainnet {
		mainnetID, err = n.InitializeMainnet(campaignID, c.SourceURL(), c.SourceHash(), chainID)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	return createChainRes.LaunchID, campaignID, mainnetID, nil
}