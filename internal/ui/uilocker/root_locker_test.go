package uilocker

import (
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
	"github.com/code-gorilla-au/rush/internal/ui/uitest"
)

func TestLockerModel_SwitchPage(t *testing.T) {
	group := odize.NewGroup(t, nil)

	state := &uistate.GlobalState{}
	theme := styles.NewIceTheme()
	ts, ps, _ := uitest.SetupServices(t)
	m := NewLockerModel(state, ts, ps, theme)

	group.Test("should update current page on MsgSwitchPage", func(t *testing.T) {
		m.Update(MsgSwitchLockerPage{NewPage: SubPageLockerPlayers})
		odize.AssertEqual(t, SubPageLockerPlayers, m.currentPage)

		m.Update(MsgSwitchLockerPage{NewPage: SubPageLockerPlaybooksList})
		odize.AssertEqual(t, SubPageLockerPlaybooksList, m.currentPage)
	})

	group.Test("should not update current page for non-locker pages", func(t *testing.T) {
		m.currentPage = SubPageLockerRoom
		m.Update(uistate.MsgSwitchPage{NewPage: uistate.PageTitle})
		odize.AssertEqual(t, SubPageLockerRoom, m.currentPage)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
