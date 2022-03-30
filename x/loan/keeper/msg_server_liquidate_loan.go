package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"

	"github.com/cosmonaut/loan/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) LiquidateLoan(goCtx context.Context, msg *types.MsgLiquidateLoan) (*types.MsgLiquidateLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid ID %d", msg.Id))
	}

	if msg.Creator != loan.Lender {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("user %d is not the lender", loan.Lender))
	}

	if loan.State != types.Approve {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("state %v is invalid", loan.Lender))
	}

	deadLine, err := strconv.ParseInt(loan.Deadline, 10, 64)
	if err != nil {
		panic(err)
	}

	if ctx.BlockHeight() < deadLine {
		return nil, sdkerrors.Wrapf(types.ErrDeadline, "Can not liquidate before deadline")
	}

	lender, _ := sdk.AccAddressFromBech32(loan.Lender)
	collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lender, collateral)
	loan.State = types.Liquidate

	k.SetLoan(ctx, loan)

	return &types.MsgLiquidateLoanResponse{}, nil
}
