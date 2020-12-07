package starportcmd

import (
	"fmt"
	"errors"
	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/numbers"
	"github.com/tendermint/starport/starport/pkg/spn"
)

func NewNetworkProposalReject() *cobra.Command {
	c := &cobra.Command{
		Use:   "reject [chain-id] [number<,...>]",
		Short: "Reject proposals",
		RunE:  networkProposalRejectHandler,
		Args:  cobra.ExactArgs(2),
	}
	return c
}

func networkProposalRejectHandler(cmd *cobra.Command, args []string) error {
	var (
		chainID      = args[0]
		proposalList = args[1]
	)

	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	var reviewals []spn.Reviewal

	ids, err := numbers.ParseList(proposalList)
	if err != nil {
		return err
	}
	for _, id := range ids {
		reviewals = append(reviewals, spn.RejectProposal(id))
	}

	gas, broadcast, err := nb.SubmitReviewals(cmd.Context(), chainID, reviewals...)
	if err != nil {
		return err
	}

	// Prompt for confirmation
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("This operation will cost about %v gas. Confirm the transaction?",
			gas,
		),
		IsConfirm: true,
	}
	if _, err := prompt.Run(); err != nil {
		return errors.New("transaction aborted")
	}

	// Broadcast the transaction
	if err := broadcast(); err != nil {
		return err
	}

	fmt.Printf("Proposal(s) %s rejected ⛔️\n", numbers.List(ids, "#"))
	return nil
}