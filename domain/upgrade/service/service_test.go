// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"database/sql"

	"github.com/golang/mock/gomock"
	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v3"
	"github.com/juju/version/v2"
	"github.com/mattn/go-sqlite3"
	gc "gopkg.in/check.v1"
)

type serviceSuite struct {
	testing.IsolationSuite
	state *MockState
}

var _ = gc.Suite(&serviceSuite{})

var (
	testUUID1 = utils.MustNewUUID().String()
	testUUID2 = utils.MustNewUUID().String()
)

func (s *serviceSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)
	s.state = NewMockState(ctrl)
	return ctrl
}

func (s *serviceSuite) TestCreateUpgrade(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.state.EXPECT().CreateUpgrade(gomock.Any(), version.MustParse("3.0.0"), version.MustParse("3.0.1")).Return(testUUID1, nil)

	upgradeUUID, err := NewService(s.state).CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(upgradeUUID, gc.Equals, testUUID1)
}

func (s *serviceSuite) TestCreateUpgradeInvalidVersions(c *gc.C) {
	_, err := NewService(s.state).CreateUpgrade(context.Background(), version.MustParse("3.0.1"), version.MustParse("3.0.0"))
	c.Assert(errors.IsNotValid(err), jc.IsTrue)

	_, err = NewService(s.state).CreateUpgrade(context.Background(), version.MustParse("3.0.1"), version.MustParse("3.0.1"))
	c.Assert(errors.IsNotValid(err), jc.IsTrue)
}

func (s *serviceSuite) TestSetControllerReady(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.state.EXPECT().SetControllerReady(gomock.Any(), testUUID1, testUUID2).Return(nil)

	err := NewService(s.state).SetControllerReady(context.Background(), testUUID1, testUUID2)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *serviceSuite) TestSetControllerReadyForiegnKey(c *gc.C) {
	defer s.setupMocks(c).Finish()

	fkErr := sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintForeignKey}
	s.state.EXPECT().SetControllerReady(gomock.Any(), testUUID1, testUUID2).Return(fkErr)

	err := NewService(s.state).SetControllerReady(context.Background(), testUUID1, testUUID2)
	c.Log(err)
	c.Assert(errors.IsNotFound(err), jc.IsTrue)
}

func (s *serviceSuite) TestAllProvisioneddControllerReadyTrue(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.state.EXPECT().AllProvisionedControllersReady(gomock.Any(), testUUID1).Return(true, nil)

	allProvisioned, err := NewService(s.state).AllProvisionedControllersReady(context.Background(), testUUID1)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(allProvisioned, jc.IsTrue)
}

func (s *serviceSuite) TestStartUpgrade(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.state.EXPECT().StartUpgrade(gomock.Any(), testUUID1).Return(nil)

	err := NewService(s.state).StartUpgrade(context.Background(), testUUID1)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *serviceSuite) TestStartUpgradeBeforeCreated(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.state.EXPECT().StartUpgrade(gomock.Any(), testUUID1).Return(sql.ErrNoRows)

	err := NewService(s.state).StartUpgrade(context.Background(), testUUID1)
	c.Assert(errors.IsNotFound(err), jc.IsTrue)
}

func (s *serviceSuite) TestActiveUpgrades(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.state.EXPECT().ActiveUpgrades(gomock.Any()).Return([]string{testUUID1, testUUID2}, nil)

	activeUpgrades, err := NewService(s.state).ActiveUpgrades(context.Background())
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(activeUpgrades, gc.DeepEquals, []string{testUUID1, testUUID2})
}
