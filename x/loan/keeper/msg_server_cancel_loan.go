package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmonaut/loan/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CancelLoan(goCtx context.Context, msg *types.MsgCancelLoan) (*types.MsgCancelLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	loan,isFound := k.GetLoan(ctx,msg.Id)
	if !isFound {
		return nil,sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,fmt.Sprintf("loan ID:%d is not exist",msg.Id))
	}

	if loan.State != types.Request {
		return nil,sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,fmt.Sprintf("state :%v must be requested",loan.State))
	}

	if loan.Borrower != msg.Creator {
		return nil,sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,fmt.Sprintf("user :%v must be borrower",msg.Creator))
	}

	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx,types.ModuleName,borrower,collateral)
	loan.State = types.Cancel

	k.SetLoan(ctx,loan)

	return &types.MsgCancelLoanResponse{}, nil
}
