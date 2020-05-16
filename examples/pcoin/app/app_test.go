package app

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/examples/pcoin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/cosmos/cosmos-sdk/client/context"
	"fmt"
)

func setGenesis(baseApp *BasecoinApp, accounts ...*types.AppAccount) (types.GenesisState, error) {
	genAccts := make([]*types.GenesisAccount, len(accounts))
	for i, appAct := range accounts {
		genAccts[i] = types.NewGenesisAccount(appAct)
	}

	genesisState := types.GenesisState{Accounts: genAccts}
	stateBytes, err := wire.MarshalJSONIndent(baseApp.cdc, genesisState)
	if err != nil {
		return types.GenesisState{}, err
	}

	// initialize and commit the chain
	baseApp.InitChain(abci.RequestInitChain{
		Validators: []abci.Validator{}, AppStateBytes: stateBytes,
	})
	baseApp.Commit()

	return genesisState, nil
}

var db = dbm.NewMemDB()

func TestGenesis(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	baseApp := NewBasecoinApp(logger, db)

	// construct a pubkey and an address for the test account
	pubkey := ed25519.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())

	// construct some test coins
	coins, err := sdk.ParseCoins("77foocoin,99barcoin")
	require.Nil(t, err)

	// create an auth.BaseAccount for the given test account and set it's coins
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	err = baseAcct.SetCoins(coins)
	require.Nil(t, err)

	// create a new test AppAccount with the given auth.BaseAccount
	appAcct := types.NewAppAccount("foobar", baseAcct)
	genState, err := setGenesis(baseApp, appAcct)
	require.Nil(t, err)

	// create a context for the BaseApp
	ctx := baseApp.BaseApp.NewContext(true, abci.Header{})
	res := baseApp.accountMapper.GetAccount(ctx, baseAcct.Address)
	require.Equal(t, appAcct, res)

	// reload app and ensure the account is still there
	baseApp = NewBasecoinApp(logger, db)

	stateBytes, err := wire.MarshalJSONIndent(baseApp.cdc, genState)
	require.Nil(t, err)

	// initialize the chain with the expected genesis state
	baseApp.InitChain(abci.RequestInitChain{
		Validators: []abci.Validator{}, AppStateBytes: stateBytes,
	})

	ctx = baseApp.BaseApp.NewContext(true, abci.Header{})
	res = baseApp.accountMapper.GetAccount(ctx, baseAcct.Address)
	require.Equal(t, appAcct, res)
}

func TestSend(t *testing.T) {
	//TestGenesis(t)

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	baseApp := NewBasecoinApp(logger, db)

	//kb, err := keys.GetKeyBase()
	ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(baseApp.cdc))
	passphrase, err := ctx.GetPassphraseFromStdin("bob")
	if err != nil {
		panic(err)
	}
	fmt.Println(passphrase)

	//// construct a pubkey and an address for the test account
	//pubkey := ed25519.GenPrivKey().PubKey()
	//addr := sdk.AccAddress(pubkey.Address())
  //
	//// create an auth.BaseAccount for the given test account and set it's coins
	//baseAcct := auth.NewBaseAccountWithAddress(addr)
	//err = baseAcct.SetCoins(coins)
	//require.Nil(t, err)
  //
	//// create a new test AppAccount with the given auth.BaseAccount
	//appAcct := types.NewAppAccount("alice", baseAcct)
	//genState, err := setGenesis(baseApp, appAcct)
	//require.Nil(t, err)
  //
	//// create a context for the BaseApp
	//ctx := baseApp.BaseApp.NewContext(true, abci.Header{})
	//res := baseApp.accountMapper.GetAccount(ctx, baseAcct.Address)
	//require.Equal(t, appAcct, res)
  //
	//// reload app and ensure the account is still there
	//baseApp = NewBasecoinApp(logger, db)
  //
	//stateBytes, err := wire.MarshalJSONIndent(baseApp.cdc, genState)
	//require.Nil(t, err)
  //
	//// initialize the chain with the expected genesis state
	//baseApp.InitChain(abci.RequestInitChain{
	//	Validators: []abci.Validator{}, AppStateBytes: stateBytes,
	//})
  //
	//ctx = baseApp.BaseApp.NewContext(true, abci.Header{})
	//res = baseApp.accountMapper.GetAccount(ctx, baseAcct.Address)
	//require.Equal(t, appAcct, res)
}
