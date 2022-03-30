package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmonaut/loan/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RepayLoan(goCtx context.Context, msg *types.MsgRepayLoan) (*types.MsgRepayLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf(" key %d not exist", msg.Id))
	}

	if loan.State != "approved" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("state %v is not approved", loan.State))
	}

	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	lender, _ := sdk.AccAddressFromBech32(loan.Lender)
	if msg.Creator != loan.Borrower {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("Can not repay %v is not borrow", borrower))
	}

	amount, _ := sdk.ParseCoinsNormalized(loan.Amount)
	fee, _ := sdk.ParseCoinsNormalized(loan.Fee)
	collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)

	err := k.bankKeeper.SendCoins(ctx, borrower, lender, amount)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Can not send coins")
	}

	err = k.bankKeeper.SendCoins(ctx, borrower, lender, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Can not send fee coins")
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrower, collateral)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Can not send collateral")
	}

	loan.State = types.Repay
	k.SetLoan(ctx, loan)

	return &types.MsgRepayLoanResponse{}, nil
}
