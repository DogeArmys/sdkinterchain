package inter_tx

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgRegister:
			return handleMsgRegister(ctx, msg, k)
		case types.MsgSend:
			return handleMsgRunTx(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleMsgRegister(ctx sdk.Context, msg types.MsgRegister, k keeper.Keeper) (*sdk.Result, error) {
	err := k.RegisterInterchainAccount(ctx, msg.Sender, msg.SourcePort, msg.SourceChannel)

	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRunTx(ctx sdk.Context, msg types.MsgSend, k keeper.Keeper) (*sdk.Result, error) {
	err := k.TrySendCoins(ctx, msg.SourcePort, msg.SourceChannel, msg.Typ, msg.Sender, msg.ToAddress, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
