package main

import (
	"errors"
	"testing"

	"github.com/mattermost/mattermost-plugin-msteams-sync/server/msteams"
	mockClient "github.com/mattermost/mattermost-plugin-msteams-sync/server/msteams/mocks"
	"github.com/mattermost/mattermost-plugin-msteams-sync/server/store"
	mockStore "github.com/mattermost/mattermost-plugin-msteams-sync/server/store/mocks"
	"github.com/mattermost/mattermost-plugin-msteams-sync/server/testutils"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecuteUnlinkCommand(t *testing.T) {
	mockAPI := &plugintest.API{}
	p := newTestPlugin()

	for _, testCase := range []struct {
		description   string
		setupAPI      func(*plugintest.API)
		args          *model.CommandArgs
		response      *model.CommandResponse
		setupStore    func(*mockStore.Store)
		setupPlugin   func()
		expectedError string
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
			response: &model.CommandResponse{},
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
			response: &model.CommandResponse{
				ResponseType: "ephemeral",
				Text:         "Unable to get the current channel information.",
				Username:     "MS Teams",
				IconURL:      "https://upload.wikimedia.org/wikipedia/commons/c/c9/Microsoft_Office_Teams_%282018%E2%80%93present%29.svg",
			},
			setupStore:    func(s *mockStore.Store) {},
			expectedError: "Error while getting the current channel.",
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
			response: &model.CommandResponse{
				ResponseType: "ephemeral",
				Text:         "Unable to unlink the channel, you has to be a channel admin to unlink it.",
				Username:     "MS Teams",
				ChannelId:    testutils.GetChannelID(),
				IconURL:      "https://upload.wikimedia.org/wikipedia/commons/c/c9/Microsoft_Office_Teams_%282018%E2%80%93present%29.svg",
			},
			setupStore:    func(s *mockStore.Store) {},
			expectedError: "Error while unlinking the channel as user is not a channel admin.",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			testCase.setupAPI(mockAPI)
			p.SetAPI(mockAPI)

			testCase.setupStore(p.store.(*mockStore.Store))
			resp, err := p.executeUnlinkCommand(&plugin.Context{}, testCase.args)
			if testCase.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, testCase.response, resp)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestExecuteShowCommand(t *testing.T) {
	p := newTestPlugin()
	mockAPI := &plugintest.API{}

	for _, testCase := range []struct {
		description   string
		args          *model.CommandArgs
		response      *model.CommandResponse
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
			response: &model.CommandResponse{},
			setupAPI: func(api *plugintest.API) {
				api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(testutils.GetPost())
			},
			setupStore: func(s *mockStore.Store) {
				s.On("GetLinkByChannelID", testutils.GetChannelID()).Return(&store.ChannelLink{}, nil)
			},
			setupClient: func(c *mockClient.Client) {
				c.On("GetTeam", mock.AnythingOfType("string")).Return(&msteams.Team{}, nil)
				c.On("GetChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&msteams.Channel{}, nil)
			},
		},
		{
			description: "Unable to get the link",
			args:        &model.CommandArgs{},
			response: &model.CommandResponse{
				ResponseType: "ephemeral",
				Text:         "Link doesn't exists.",
				Username:     "MS Teams",
				IconURL:      "https://upload.wikimedia.org/wikipedia/commons/c/c9/Microsoft_Office_Teams_%282018%E2%80%93present%29.svg",
			},
			setupAPI: func(api *plugintest.API) {},
			setupStore: func(s *mockStore.Store) {
				s.On("GetLinkByChannelID", "").Return(nil, errors.New("Error while getting the link"))
			},
			setupClient:   func(c *mockClient.Client) {},
			expectedError: "Error while getting link.",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			testCase.setupAPI(mockAPI)
			p.SetAPI(mockAPI)

			testCase.setupStore(p.store.(*mockStore.Store))
			testCase.setupClient(p.msteamsAppClient.(*mockClient.Client))
			resp, err := p.executeShowCommand(&plugin.Context{}, testCase.args)

			if testCase.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, testCase.response, resp)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
