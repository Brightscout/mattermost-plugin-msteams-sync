package main

import (
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
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

func TestExecuteCommand(t *testing.T) {
	defer monkey.UnpatchAll()
	p := Plugin{}

	for _, testCase := range []struct {
		description string
		setupPlugin func()
		args        *model.CommandArgs
	}{
		{
			description: "Invalid command",
			args: &model.CommandArgs{
				Command: "/invalid",
			},
			setupPlugin: func() {},
		},
		{
			description: "Link command",
			args: &model.CommandArgs{
				Command: "/msteams-sync link",
			},
			setupPlugin: func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(&p), "ExecuteLinkCommand", func(*Plugin, *plugin.Context, *model.CommandArgs, []string) (*model.CommandResponse, *model.AppError) {
					return &model.CommandResponse{}, nil
				})
			},
		},
		{
			description: "UnLink command",
			args: &model.CommandArgs{
				Command: "/msteams-sync unlink",
			},
			setupPlugin: func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(&p), "ExecuteUnlinkCommand", func(*Plugin, *plugin.Context, *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
					return &model.CommandResponse{}, nil
				})
			},
		},
		{
			description: "Show command",
			args: &model.CommandArgs{
				Command: "/msteams-sync show",
			},
			setupPlugin: func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(&p), "ExecuteShowCommand", func(*Plugin, *plugin.Context, *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
					return &model.CommandResponse{}, nil
				})
			},
		},
		{
			description: "Connect command",
			args: &model.CommandArgs{
				Command: "/msteams-sync connect",
			},
			setupPlugin: func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(&p), "ExecuteConnectCommand", func(*Plugin, *plugin.Context, *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
					return &model.CommandResponse{}, nil
				})
			},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			assert := assert.New(t)
			testCase.setupPlugin()

			resp, err := p.ExecuteCommand(&plugin.Context{}, testCase.args)
			assert.EqualValues(&model.CommandResponse{}, resp)
			assert.Nil(err)
		})
	}
}

func TestExecuteUnlinkCommand(t *testing.T) {
	mockAPI := &plugintest.API{}
	p := newTestPlugin()

	for _, testCase := range []struct {
		description   string
		setupAPI      func(*plugintest.API)
		args          *model.CommandArgs
		response      *model.CommandResponse
		setupStore    func(*mockStore.Store) *mockStore.Store
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
			setupStore: func(s *mockStore.Store) *mockStore.Store {
				s.On("DeleteLinkByChannelID", mock.AnythingOfType("string")).Return(nil)
				return s
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
			setupStore: func(s *mockStore.Store) *mockStore.Store {
				return s
			},
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
			setupStore: func(s *mockStore.Store) *mockStore.Store {
				return s
			},
			expectedError: "Error while unlinking the channel as user is not a channel admin.",
		},
		{
			description: "Unable to delete link.",
			setupAPI: func(api *plugintest.API) {
				api.On("GetChannel", testutils.GetChannelID()).Return(&model.Channel{
					Type: model.ChannelTypeOpen,
				}, nil)
				api.On("HasPermissionToChannel", testutils.GetUserID(), testutils.GetChannelID(), model.PermissionManagePublicChannelProperties).Return(true)
			},
			args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
			},
			response: &model.CommandResponse{
				ResponseType: "ephemeral",
				Text:         "Unable to delete link.",
				Username:     "MS Teams",
				ChannelId:    testutils.GetChannelID(),
				IconURL:      "https://upload.wikimedia.org/wikipedia/commons/c/c9/Microsoft_Office_Teams_%282018%E2%80%93present%29.svg",
			},
			setupStore: func(s *mockStore.Store) *mockStore.Store {
				s.On("DeleteLinkByChannelID", mock.AnythingOfType("string")).Return(errors.New("unable to delete a channel"))
				return s
			},
			expectedError: "Error while deleting a channel.",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			testCase.setupAPI(mockAPI)
			p.SetAPI(mockAPI)
			mockStore.NewStore(t)
			store := testCase.setupStore(mockStore.NewStore(t))
			p.store = store
			resp, err := p.ExecuteUnlinkCommand(&plugin.Context{}, testCase.args)

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
		setupStore    func(*mockStore.Store) *mockStore.Store
		setupClient   func(*mockClient.Client) *mockClient.Client
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
			setupStore: func(s *mockStore.Store) *mockStore.Store {
				s.On("GetLinkByChannelID", testutils.GetChannelID()).Return(&store.ChannelLink{}, nil)
				return s
			},
			setupClient: func(c *mockClient.Client) *mockClient.Client {
				c.On("GetTeam", mock.AnythingOfType("string")).Return(&msteams.Team{}, nil)
				c.On("GetChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&msteams.Channel{}, nil)
				return c
			},
		},
		{
			description: "Unable to get the link",
			args: &model.CommandArgs{
				UserId:    "",
				ChannelId: "",
			},
			response: &model.CommandResponse{
				ResponseType: "ephemeral",
				Text:         "Link doesn't exists.",
				Username:     "MS Teams",
				IconURL:      "https://upload.wikimedia.org/wikipedia/commons/c/c9/Microsoft_Office_Teams_%282018%E2%80%93present%29.svg",
			},
			setupAPI: func(api *plugintest.API) {},
			setupStore: func(s *mockStore.Store) *mockStore.Store {
				s.On("GetLinkByChannelID", "").Return(nil, errors.New("Error while getting the link"))
				return s
			},
			setupClient: func(c *mockClient.Client) *mockClient.Client {
				return c
			},
			expectedError: "Error while getting link.",
		},
		{
			description: "Unable to get the MS Teams team information.",
			args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
			},
			response: &model.CommandResponse{
				ResponseType: "ephemeral",
				Text:         "Unable to get the MS Teams team information.",
				Username:     "MS Teams",
				ChannelId:    testutils.GetChannelID(),
				IconURL:      "https://upload.wikimedia.org/wikipedia/commons/c/c9/Microsoft_Office_Teams_%282018%E2%80%93present%29.svg",
			},
			setupAPI: func(api *plugintest.API) {},
			setupStore: func(s *mockStore.Store) *mockStore.Store {
				s.On("GetLinkByChannelID", testutils.GetChannelID()).Return(&store.ChannelLink{}, nil)
				return s
			},
			setupClient: func(c *mockClient.Client) *mockClient.Client {
				c.On("GetTeam", mock.Anything).Return(nil, errors.New("Error while getting the MS Teams team information"))
				return c
			},
			expectedError: "Error while getting the MS Teams team information.",
		},
		{
			description: "Unable to get the MS Teams channel information.",
			args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
			},
			response: &model.CommandResponse{
				ResponseType: "ephemeral",
				Text:         "Unable to get the MS Teams channel information.",
				Username:     "MS Teams",
				ChannelId:    testutils.GetChannelID(),
				IconURL:      "https://upload.wikimedia.org/wikipedia/commons/c/c9/Microsoft_Office_Teams_%282018%E2%80%93present%29.svg",
			},
			setupAPI: func(api *plugintest.API) {},
			setupStore: func(s *mockStore.Store) *mockStore.Store {
				s.On("GetLinkByChannelID", testutils.GetChannelID()).Return(&store.ChannelLink{}, nil)
				return s
			},
			setupClient: func(c *mockClient.Client) *mockClient.Client {
				c.On("GetTeam", mock.AnythingOfType("string")).Return(&msteams.Team{}, nil)
				c.On("GetChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Error while getting the MS Teams channel information"))
				return c
			},
			expectedError: "Error while getting the MS Teams channel information.",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			testCase.setupAPI(mockAPI)
			p.SetAPI(mockAPI)
			mockStore.NewStore(t)
			mockClient.NewClient(t)
			store := testCase.setupStore(mockStore.NewStore(t))
			client := testCase.setupClient(mockClient.NewClient(t))

			p.store = store
			p.msteamsAppClient = client
			resp, err := p.ExecuteShowCommand(&plugin.Context{}, testCase.args)

			if testCase.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, testCase.response, resp)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
