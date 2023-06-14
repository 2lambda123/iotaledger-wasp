package chains

import (
	"time"

	"github.com/iotaledger/hive.go/app"
)

type ParametersChains struct {
	BroadcastUpToNPeers              int           `default:"2" usage:"number of peers an offledger request is broadcasted to"`
	BroadcastInterval                time.Duration `default:"5s" usage:"time between re-broadcast of offledger requests"`
	APICacheTTL                      time.Duration `default:"300s" usage:"time to keep processed offledger requests in api cache"`
	PullMissingRequestsFromCommittee bool          `default:"true" usage:"whether or not to pull missing requests from other committee members"`
	DeriveAliasOutputByQuorum        bool          `default:"true" usage:"false means we propose own AliasOutput, true - by majority vote."`
	PipeliningLimit                  int           `default:"-1" usage:"-1 -- infinite, 0 -- disabled, X -- build the chain if there is up to X transactions unconfirmed by L1."`
	ConsensusDelay                   time.Duration `default:"500ms" usage:"Minimal delay between consensus runs."`
}

type ParametersWAL struct {
	Enabled bool   `default:"true" usage:"whether the \"write-ahead logging\" is enabled"`
	Path    string `default:"waspdb/wal" usage:"the path to the \"write-ahead logging\" folder"`
}

type ParametersValidator struct {
	Address string `default:"" usage:"bech32 encoded address to collect validator fee payments"`
}

var (
	ParamsChains    = &ParametersChains{}
	ParamsWAL       = &ParametersWAL{}
	ParamsValidator = &ParametersValidator{}
)

var params = &app.ComponentParams{
	Params: map[string]any{
		"chains":    ParamsChains,
		"wal":       ParamsWAL,
		"validator": ParamsValidator,
	},
	Masked: nil,
}
