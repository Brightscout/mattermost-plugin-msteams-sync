package main

import (
	"testing"

	"github.com/mattermost/mattermost-plugin-msteams-sync/server/msteams"
	mockClient "github.com/mattermost/mattermost-plugin-msteams-sync/server/msteams/mocks"
	mockStore "github.com/mattermost/mattermost-plugin-msteams-sync/server/store/mocks"
	"github.com/mattermost/mattermost-plugin-msteams-sync/server/store/storemodels"
	"github.com/mattermost/mattermost-plugin-msteams-sync/server/testutils"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/mock"
)

func TestExecuteUnlinkCommand(t *testing.T) {
	mockAPI := &plugintest.API{}
	p := newTestPlugin()

	for _, testCase := range []struct {
		description string
		setupAPI    func(*plugintest.API)
		args        *model.CommandArgs
		setupStore  func(*mockStore.Store)
		setupPlugin func()
	}{
		{
			description: "Successfully executed unlinked command",
			setupAPI: func(api *plugintest.API) {
				api.On("GetChannel", testutils.GetChannelID()).Return(&model.Channel{
					Type: model.ChannelTypeOpen,
				}, nil)
				api.On("HasPermissionToChannel", testutils.GetUserID(), testutils.GetChannelID(), model.PermissionManagePublicChannelProperties).Return(true)
				api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(testutils.GetPost())
			},
			args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
			},
			setupStore: func(s *mockStore.Store) {
				s.On("DeleteLinkByChannelID", mock.AnythingOfType("string")).Return(nil)
			},
		},
		{
			description: "Unable to get the current channel",
			setupAPI: func(api *plugintest.API) {
				api.On("GetChannel", "").Return(nil, testutils.GetInternalServerAppError("Error while getting the current channel."))
			},
			args: &model.CommandArgs{
				UserId:    "",
				ChannelId: "",
			},
			setupStore: func(s *mockStore.Store) {},
		},
		{
			description: "Unable to unlink channel as user is not a channel admin.",
			setupAPI: func(api *plugintest.API) {
				api.On("GetChannel", testutils.GetChannelID()).Return(&model.Channel{
					Type: model.ChannelTypeOpen,
				}, nil)
				api.On("HasPermissionToChannel", "", testutils.GetChannelID(), model.PermissionManagePublicChannelProperties).Return(false)
			},
			args: &model.CommandArgs{
				ChannelId: testutils.GetChannelID(),
			},
			setupStore: func(s *mockStore.Store) {},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			testCase.setupAPI(mockAPI)
			p.SetAPI(mockAPI)

			testCase.setupStore(p.store.(*mockStore.Store))
			p.executeUnlinkCommand(&plugin.Context{}, testCase.args)
		})
	}
}

func TestExecuteShowCommand(t *testing.T) {
	p := newTestPlugin()
	mockAPI := &plugintest.API{}

	for _, testCase := range []struct {
		description   string
		args          *model.CommandArgs
		setupAPI      func(*plugintest.API)
		setupStore    func(*mockStore.Store)
		setupClient   func(*mockClient.Client)
		setupPlugin   func()
		expectedError string
	}{
		{
			description: "Successfully executed show command",
			args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
			},
			setupAPI: func(api *plugintest.API) {
				api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(testutils.GetPost())
			},
			setupStore: func(s *mockStore.Store) {
				s.On("GetLinkByChannelID", testutils.GetChannelID()).Return(&storemodels.ChannelLink{
					MSTeamsTeam: "Valid-MSTeamsTeam",
				}, nil)
			},
			setupClient: func(c *mockClient.Client) {
				c.On("GetTeam", "Valid-MSTeamsTeam").Return(&msteams.Team{}, nil)
				c.On("GetChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&msteams.Channel{}, nil)
			},
		},
		{
			description: "Unable to get the link",
			args:        &model.CommandArgs{},
			setupAPI: func(api *plugintest.API) {
				api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(testutils.GetPost())
			},
			setupStore: func(s *mockStore.Store) {
				s.On("GetLinkByChannelID", "").Return(nil, errors.New("Error while getting the link"))
			},
			setupClient:   func(c *mockClient.Client) {},
			expectedError: "Error while getting link.",
		},
		{
			description: "Unable to get the MS Teams team information",
			args: &model.CommandArgs{
				ChannelId: "Invalid-ChannelID",
			},
			setupAPI: func(api *plugintest.API) {
				api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(testutils.GetPost())
			},
			setupStore: func(s *mockStore.Store) {
				s.On("GetLinkByChannelID", "Invalid-ChannelID").Return(&storemodels.ChannelLink{
					MSTeamsTeam: "Invalid-MSTeamsTeam",
				}, nil)
			},
			setupClient: func(c *mockClient.Client) {
				c.On("GetTeam", "Invalid-MSTeamsTeam").Return(nil, errors.New("Error while getting the MS Teams team information."))
			},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			testCase.setupAPI(mockAPI)
			p.SetAPI(mockAPI)

			testCase.setupStore(p.store.(*mockStore.Store))
			testCase.setupClient(p.msteamsAppClient.(*mockClient.Client))
			p.executeShowCommand(&plugin.Context{}, testCase.args)
		})
	}
}
