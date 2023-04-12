package handlers

import (
	"errors"
	"fmt"
	"testing"

	mocksPlugin "github.com/mattermost/mattermost-plugin-msteams-sync/server/handlers/mocks"
	"github.com/mattermost/mattermost-plugin-msteams-sync/server/msteams"
	mocksClient "github.com/mattermost/mattermost-plugin-msteams-sync/server/msteams/mocks"
	mocksStore "github.com/mattermost/mattermost-plugin-msteams-sync/server/store/mocks"
	"github.com/mattermost/mattermost-plugin-msteams-sync/server/store/storemodels"
	"github.com/mattermost/mattermost-plugin-msteams-sync/server/testutils"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"
	"github.com/stretchr/testify/mock"
)

func TestHandleCreatedActivity(t *testing.T) {
	ah := ActivityHandler{}
	client := mocksClient.NewClient(t)
	store := mocksStore.NewStore(t)
	mockAPI := &plugintest.API{}

	for _, testCase := range []struct {
		description string
		activityIds msteams.ActivityIds
		setupPlugin func(plugin *mocksPlugin.PluginIface)
		setupAPI    func()
		setupClient func()
	}{
		{
			description: "Unable to get original message",
			activityIds: msteams.ActivityIds{
				ChatID: "invalid-ChatID",
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetClientForApp").Return(client)
				p.On("GetAPI").Return(mockAPI)
			},
			setupClient: func() {
				client.On("GetChat", "invalid-ChatID").Return(nil, errors.New("Error while getting original chat"))
				client.On("GetReply", testutils.GetTeamUserID(), testutils.GetChannelID(), testutils.GetMessageID(), testutils.GetReplyID()).Return(&msteams.Message{
					UserID:    testutils.GetUserID(),
					TeamID:    testutils.GetTeamUserID(),
					ChannelID: testutils.GetChannelID(),
				}, nil)
			},
			setupAPI: func() {},
		},
		{
			description: "Skipping not user event",
			activityIds: msteams.ActivityIds{
				ChatID:    testutils.GetChatID(),
				MessageID: "mock-MessageID1",
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetClientForApp").Return(client)
				p.On("GetClientForTeamsUser", testutils.GetUserID()).Return(client, nil)
				p.On("GetAPI").Return(mockAPI)
			},
			setupClient: func() {
				client.On("GetChat", testutils.GetChatID()).Return(&msteams.Chat{
					Members: []msteams.ChatMember{
						{
							UserID: testutils.GetUserID(),
						},
					},
					ID: testutils.GetChatID(),
				}, nil)
				client.On("GetChatMessage", testutils.GetChatID(), "mock-MessageID1").Return(&msteams.Message{}, nil)
			},
			setupAPI: func() {},
		},
		{
			description: "Skipping messages from bot user",
			activityIds: msteams.ActivityIds{
				ChatID:    testutils.GetChatID(),
				MessageID: "mock-MessageID2",
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetClientForApp").Return(client)
				p.On("GetClientForTeamsUser", testutils.GetUserID()).Return(client, nil)
				p.On("GetAPI").Return(mockAPI)
				p.On("GetStore").Return(store)
				p.On("GetBotUserID").Return("mock-BotUserID")
				store.On("MattermostToTeamsUserID", "mock-BotUserID").Return("mock-UserID1", nil)
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil)
			},
			setupClient: func() {
				client.On("GetChat", testutils.GetChatID()).Return(&msteams.Chat{
					Members: []msteams.ChatMember{
						{
							UserID: testutils.GetUserID(),
						},
					},
					ID: testutils.GetChatID(),
				}, nil)
				client.On("GetChatMessage", testutils.GetChatID(), "mock-MessageID2").Return(&msteams.Message{
					UserID: "mock-UserID1",
				}, nil)
			},
			setupAPI: func() {},
		},
		{
			description: "Successfully handled created activity",
			activityIds: msteams.ActivityIds{
				ReplyID:   testutils.GetReplyID(),
				MessageID: testutils.GetMessageID(),
				TeamID:    testutils.GetTeamUserID(),
				ChannelID: testutils.GetChannelID(),
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetStore").Return(store)
				store.On("GetLinkByMSTeamsChannelID", testutils.GetTeamUserID(), testutils.GetChannelID()).Return(&storemodels.ChannelLink{
					Creator:           testutils.GetUserID(),
					MattermostChannel: "mock-MattermostChannel",
				}, nil)
				p.On("GetClientForUser", testutils.GetUserID()).Return(client, nil)
				p.On("GetBotUserID").Return("mock-BotUserID")
				store.On("TeamsToMattermostUserID", testutils.GetUserID()).Return("", nil)
				p.On("GetURL").Return("https://example.com/")
				p.On("GetAPI").Return(mockAPI)
				store.On("GetPostInfoByMSTeamsID", testutils.GetChannelID(), "").Return(&storemodels.PostInfo{}, nil)
			},
			setupClient: func() {
				client.On("GetReply", testutils.GetTeamUserID(), testutils.GetChannelID(), testutils.GetMessageID(), testutils.GetReplyID()).Return(&msteams.Message{
					UserID:    testutils.GetUserID(),
					TeamID:    testutils.GetTeamUserID(),
					ChannelID: testutils.GetChannelID(),
				}, nil)
			},
			setupAPI: func() {},
		},
		// {
		// 	description: "Unable to get original channel id",
		// 	activityIds: msteams.ActivityIds{
		// 		ChatID:    "mock-ChatID",
		// 		MessageID: "mock-MessageID3",
		// 	},
		// 	setupPlugin: func(p *mocksPlugin.PluginIface) {
		// 		p.On("GetClientForApp").Return(client)
		// 		p.On("GetClientForTeamsUser", testutils.GetUserID()).Return(client, nil)
		// 		p.On("GetAPI").Return(mockAPI)
		// 		p.On("GetStore").Return(store)
		// 		p.On("GetBotUserID").Return("mock-BotUserID1")
		// 		store.On("TeamsToMattermostUserID", "mock-UserID2").Return(testutils.GetSenderID(), nil)
		// 		store.On("MattermostToTeamsUserID", "mock-BotUserID1").Return("mock-UserID1", nil)
		// 		mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		// 		// p.On("GetSyncDirectMessages").Return(false)
		// 	},
		// 	setupClient: func() {
		// 		client.On("GetChat", "mock-ChatID").Return(nil, nil)
		// 		client.On("GetChatMessage", testutils.GetChatID(), "mock-MessageID3").Return(&msteams.Message{
		// 			UserID: "mock-UserID2",
		// 		}, nil)
		// 	},
		// 	setupAPI: func() {},
		// },

		// {
		// 	description: "Unable to get the message",
		// 	activityIds: msteams.ActivityIds{
		// 		MessageID: testutils.GetMessageID(),
		// 		ChatID:    testutils.GetChatID(),
		// 	},
		// 	setupPlugin: func(p *mocksPlugin.PluginIface) {
		// 		client := mocksClient.NewClient(t)
		// 		p.On("GetClientForApp").Return(client)
		// 		client.On("GetChat", testutils.GetChatID()).Return(&msteams.Chat{}, nil)
		// 		p.On("GetAPI").Return(mockAPI)
		// 	},
		// 	setupAPI: func() {
		// 		mockAPI.On("LogDebug", "Unable to get the message (probably because belongs to private chate in not-linked users)")
		// 	},
		// },
		// {
		// 	description: "Skipping not user event",
		// 	activityIds: msteams.ActivityIds{
		// 		MessageID: testutils.GetMessageID(),
		// 		ChatID:    testutils.GetChatID(),
		// 	},
		// 	setupPlugin: func(p *mocksPlugin.PluginIface) {
		// 		client := mocksClient.NewClient(t)
		// 		p.On("GetClientForApp").Return(client)
		// 		client.On("GetChat", testutils.GetChatID()).Return(&msteams.Chat{
		// 			Members: []msteams.ChatMember{
		// 				{
		// 					UserID: testutils.GetUserID(),
		// 				},
		// 			},
		// 		}, nil)
		// 		p.On("GetClientForTeamsUser", testutils.GetUserID()).Return(client, nil)
		// 		client.On("GetChatMessage", "", testutils.GetMessageID()).Return(nil, errors.New("Error while getting the original Post"))
		// 		p.On("GetAPI").Return(mockAPI)
		// 	},
		// 	setupAPI: func() {
		// 		mockAPI.On("LogError", "Unable to get original post", "error", errors.New("Error while getting the original Post"))
		// 	},
		// },
		// {
		// 	description: "Skipping messages from bot user",
		// 	activityIds: msteams.ActivityIds{
		// 		MessageID: testutils.GetMessageID(),
		// 		ChatID:    testutils.GetChatID(),
		// 	},
		// 	setupPlugin: func(p *mocksPlugin.PluginIface) {
		// 		client := mocksClient.NewClient(t)
		// 		store := mocksStore.NewStore(t)
		// 		p.On("GetClientForApp").Return(client)
		// 		client.On("GetChat", testutils.GetChatID()).Return(&msteams.Chat{
		// 			ID: testutils.GetChatID(),
		// 			Members: []msteams.ChatMember{
		// 				{
		// 					UserID: testutils.GetUserID(),
		// 				},
		// 			},
		// 		}, nil)
		// 		p.On("GetClientForTeamsUser", testutils.GetUserID()).Return(client, nil)
		// 		client.On("GetChatMessage", testutils.GetChatID(), testutils.GetMessageID()).Return(&msteams.Message{
		// 			UserID: testutils.GetUserID(),
		// 		}, nil)
		// 		p.On("GetAPI").Return(mockAPI)
		// 		p.On("GetStore").Return(store)
		// 		p.On("GetBotUserID").Return("")
		// 		store.On("MattermostToTeamsUserID", "").Return(testutils.GetUserID(), nil)
		// 	},
		// 	setupAPI: func() {
		// 		mockAPI.On("KVSet", lastReceivedChangeKey, mock.Anything).Return(nil)
		// 		mockAPI.On("LogDebug", "Skipping messages from bot user")
		// 	},
		// },
		// {
		// 	description: "Unable to get original channel id",
		// 	activityIds: msteams.ActivityIds{
		// 		MessageID: testutils.GetMessageID(),
		// 		ChatID:    testutils.GetChatID(),
		// 	},
		// 	setupPlugin: func(p *mocksPlugin.PluginIface) {
		// 		client := mocksClient.NewClient(t)
		// 		store := mocksStore.NewStore(t)
		// 		p.On("GetClientForApp").Return(client)
		// 		client.On("GetChat", testutils.GetChatID()).Return(&msteams.Chat{
		// 			ID: testutils.GetChatID(),
		// 			Members: []msteams.ChatMember{
		// 				{
		// 					UserID: testutils.GetUserID(),
		// 				},
		// 			},
		// 		}, nil)
		// 		p.On("GetClientForTeamsUser", testutils.GetUserID()).Return(client, nil)
		// 		client.On("GetChatMessage", testutils.GetChatID(), testutils.GetMessageID()).Return(&msteams.Message{
		// 			UserID: testutils.GetUserID(),
		// 		}, nil)
		// 		p.On("GetAPI").Return(mockAPI)
		// 		p.On("GetStore").Return(store)
		// 		p.On("GetBotUserID").Return("mock-BotUserID")
		// 		p.On("GetSyncDirectMessages").Return(true)
		// 		var token *oauth2.Token
		// 		store.On("SetUserInfo", "mock-UserID", testutils.GetUserID(), token).Return(nil)
		// 		store.On("MattermostToTeamsUserID", "mock-BotUserID").Return("mock-msteamsUserID", nil)
		// 		store.On("TeamsToMattermostUserID", testutils.GetUserID()).Return("", nil)
		// 	},
		// 	setupAPI: func() {
		// 		mockAPI.On("GetUserByEmail", testutils.GetUserID()+"@msteamssync").Return(&model.User{
		// 			Id: "mock-UserID",
		// 		}, nil)
		// 	},
		// },
		// {
		// 	description: "Channel not set",
		// 	activityIds: msteams.ActivityIds{
		// 		MessageID: testutils.GetMessageID(),
		// 		ReplyID:   testutils.GetAPIReplyID(),
		// 	},
		// 	setupPlugin: func(p *mocksPlugin.PluginIface) {
		// client := mocksClient.NewClient(t)
		// store := mocksStore.NewStore(t)
		// p.On("GetClientForApp").Return(client)
		// client.On("GetChat", "mock-ChatID").Return(&msteams.Chat{
		// 	ID: "mock-ChatID",
		// 	Members: []msteams.ChatMember{
		// 		{
		// 			UserID: testutils.GetUserID(),
		// 		},
		// 	},
		// }, nil)
		// p.On("GetClientForTeamsUser", testutils.GetUserID()).Return(client, nil)
		// client.On("GetChatMessage", "mock-ChatID", testutils.GetMessageID()).Return(&msteams.Message{
		// 	UserID: testutils.GetUserID(),
		// }, nil)
		// p.On("GetAPI").Return(mockAPI)
		// p.On("GetStore").Return(store)
		// p.On("GetBotUserID").Return("mock-BotUserID")
		// p.On("GetSyncDirectMessages").Return(true)
		// var token *oauth2.Token
		// store.On("SetUserInfo", "mock-UserID", testutils.GetUserID(), token).Return(nil)
		// store.On("MattermostToTeamsUserID", "mock-BotUserID").Return("mock-msteamsUserID", nil)
		// store.On("TeamsToMattermostUserID", testutils.GetUserID()).Return("", errors.New("Error while getting mmUserID"))
		// 	},
		// 	setupAPI: func() {
		// 		mockAPI.On("GetUserByEmail", testutils.GetUserID()+"@msteamssync").Return(&model.User{
		// 			Id: "mock-UserID",
		// 		}, nil)
		// 	},
		// },
		// {
		// 	description: "Unable to transform teams post in mattermost post",
		// 	activityIds: msteams.ActivityIds{
		// 		TeamID:    testutils.GetTeamUserID(),
		// 		ChatID:    testutils.GetChatID(),
		// 		ChannelID: testutils.GetChannelID(),
		// 	},
		// 	setupPlugin: func(p *mocksPlugin.PluginIface) {
		// 		client := mocksClient.NewClient(t)
		// 		store := mocksStore.NewStore(t)
		// 		p.On("GetClientForApp").Return(client)
		// 		client.On("GetChat", testutils.GetChatID()).Return(&msteams.Chat{
		// 			ID: testutils.GetChatID(),
		// 			Members: []msteams.ChatMember{
		// 				{
		// 					UserID: testutils.GetUserID(),
		// 				},
		// 			},
		// 		}, nil)
		// 		p.On("GetClientForTeamsUser", testutils.GetUserID()).Return(client, nil)
		// 		client.On("GetChatMessage", testutils.GetChatID(), testutils.GetMessageID()).Return(&msteams.Message{
		// 			UserID: testutils.GetUserID(),
		// 		}, nil)
		// 		p.On("GetAPI").Return(mockAPI)
		// 		p.On("GetStore").Return(store)
		// 		p.On("GetBotUserID").Return("mock-BotUserID")
		// 		p.On("GetSyncDirectMessages").Return(true)
		// 		var token *oauth2.Token
		// 		store.On("SetUserInfo", "mock-UserID", testutils.GetUserID(), token).Return(nil)
		// 		store.On("MattermostToTeamsUserID", "mock-BotUserID").Return("mock-msteamsUserID", nil)
		// 		store.On("TeamsToMattermostUserID", testutils.GetUserID()).Return("", errors.New("Error while getting mmUserID"))
		// 	},
		// 	setupAPI: func() {
		// 		mockAPI.On("GetUserByEmail", testutils.GetUserID()+"@msteamssync").Return(&model.User{
		// 			Id: "mock-UserID",
		// 		}, nil)
		// 	},
		// },
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.Anything, mock.Anything, mock.Anything)
			mockAPI.On("LogDebug", mock.Anything, mock.Anything, mock.Anything)
			p := mocksPlugin.NewPluginIface(t)
			testCase.setupPlugin(p)
			testCase.setupAPI()
			testCase.setupClient()
			ah.plugin = p

			ah.handleCreatedActivity(testCase.activityIds)
		})
	}
}

func TestHandleUpdatedActivity(t *testing.T) {
	ah := ActivityHandler{}
	store := mocksStore.NewStore(t)
	mockAPI := &plugintest.API{}
	client := mocksClient.NewClient(t)

	for _, testCase := range []struct {
		description string
		activityIds msteams.ActivityIds
		setupPlugin func(plugin *mocksPlugin.PluginIface)
		setupClient func()
		setupAPI    func()
	}{
		{
			description: "Unable to get original message",
			activityIds: msteams.ActivityIds{
				ChatID: "invalid-ChatID",
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetClientForApp").Return(client)
				p.On("GetAPI").Return(mockAPI)
			},
			setupClient: func() {
				client.On("GetChat", "invalid-ChatID").Return(nil, errors.New("Error while getting original chat"))
			},
			setupAPI: func() {},
		},
		{
			description: "Skipping not user event",
			activityIds: msteams.ActivityIds{
				ChatID:    testutils.GetChatID(),
				MessageID: "mock-MessageID1",
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetClientForApp").Return(client)
				p.On("GetClientForTeamsUser", testutils.GetUserID()).Return(client, nil)
				p.On("GetAPI").Return(mockAPI)
			},
			setupClient: func() {
				client.On("GetChat", testutils.GetChatID()).Return(&msteams.Chat{
					Members: []msteams.ChatMember{
						{
							UserID: testutils.GetUserID(),
						},
					},
					ID: testutils.GetChatID(),
				}, nil)
				client.On("GetChatMessage", testutils.GetChatID(), "mock-MessageID1").Return(&msteams.Message{}, nil)
			},
			setupAPI: func() {},
		},
		{
			description: "Skipping messages from bot user",
			activityIds: msteams.ActivityIds{
				ChatID:    testutils.GetChatID(),
				MessageID: "mock-MessageID2",
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetClientForApp").Return(client)
				p.On("GetClientForTeamsUser", testutils.GetUserID()).Return(client, nil)
				p.On("GetAPI").Return(mockAPI)
				p.On("GetStore").Return(store)
				p.On("GetBotUserID").Return("mock-BotUserID")
				store.On("MattermostToTeamsUserID", "mock-BotUserID").Return("mock-UserID1", nil)
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil)
			},
			setupClient: func() {
				client.On("GetChat", testutils.GetChatID()).Return(&msteams.Chat{
					Members: []msteams.ChatMember{
						{
							UserID: testutils.GetUserID(),
						},
					},
					ID: testutils.GetChatID(),
				}, nil)
				client.On("GetChatMessage", testutils.GetChatID(), "mock-MessageID2").Return(&msteams.Message{
					UserID: "mock-UserID1",
				}, nil)
			},
			setupAPI: func() {},
		},
		{
			description: "Successfully update last received change date",
			activityIds: msteams.ActivityIds{
				ReplyID:   testutils.GetReplyID(),
				MessageID: testutils.GetMessageID(),
				TeamID:    testutils.GetTeamUserID(),
				ChannelID: testutils.GetChannelID(),
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetStore").Return(store)
				store.On("GetLinkByMSTeamsChannelID", testutils.GetTeamUserID(), testutils.GetChannelID()).Return(&storemodels.ChannelLink{
					Creator:           testutils.GetUserID(),
					MattermostChannel: "mock-MattermostChannel",
				}, nil)
				p.On("GetClientForUser", testutils.GetUserID()).Return(client, nil)
				p.On("GetBotUserID").Return("mock-BotUserID")
				store.On("GetPostInfoByMSTeamsID", testutils.GetChannelID(), "").Return(nil, nil)
			},
			setupClient: func() {
				client.On("GetReply", testutils.GetTeamUserID(), testutils.GetChannelID(), testutils.GetMessageID(), testutils.GetReplyID()).Return(&msteams.Message{
					UserID:    testutils.GetUserID(),
					TeamID:    testutils.GetTeamUserID(),
					ChannelID: testutils.GetChannelID(),
				}, nil)
			},
			setupAPI: func() {},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			mockAPI.On("LogError", mock.Anything, mock.Anything, mock.Anything)
			mockAPI.On("LogDebug", mock.Anything, mock.Anything, mock.Anything)
			p := mocksPlugin.NewPluginIface(t)
			testCase.setupPlugin(p)
			testCase.setupClient()
			testCase.setupAPI()

			ah.plugin = p
			ah.handleUpdatedActivity(testCase.activityIds)
		})
	}
}

func TestHandleDeletedActivity(t *testing.T) {
	ah := ActivityHandler{}
	store := mocksStore.NewStore(t)
	mockAPI := &plugintest.API{}

	for _, testCase := range []struct {
		description string
		activityIds msteams.ActivityIds
		setupPlugin func(plugin *mocksPlugin.PluginIface)
		setupAPI    func()
		setupStore  func()
	}{
		{
			description: "Successfully deleted post",
			activityIds: msteams.ActivityIds{
				ChatID:    testutils.GetChatID(),
				ChannelID: testutils.GetChannelID(),
				MessageID: testutils.GetMessageID(),
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetStore").Return(store)
				p.On("GetAPI").Return(mockAPI)
			},
			setupAPI: func() {
				mockAPI.On("DeletePost", testutils.GetMattermostID()).Return(nil)
				mockAPI.On("LogError", "Unable to to delete post", "msgID", "", "error", &model.AppError{
					Message: "Error while deleting a post",
				})
			},
			setupStore: func() {
				store.On("GetPostInfoByMSTeamsID", fmt.Sprintf("%s%s", testutils.GetChatID(), testutils.GetChannelID()), testutils.GetMessageID()).Return(&storemodels.PostInfo{
					MattermostID: testutils.GetMattermostID(),
				}, nil)
			},
		},
		{
			description: "Unable to get post info by MS teams ID",
			activityIds: msteams.ActivityIds{
				ChannelID: testutils.GetChannelID(),
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetStore").Return(store)
			},
			setupAPI: func() {},
			setupStore: func() {
				store.On("GetPostInfoByMSTeamsID", testutils.GetChannelID(), "").Return(nil, errors.New("Error while getting post info by MS teams ID"))
			},
		},
		{
			description: "Unable to to delete post",
			activityIds: msteams.ActivityIds{
				ChannelID: testutils.GetChannelID(),
				MessageID: testutils.GetMessageID(),
			},
			setupPlugin: func(p *mocksPlugin.PluginIface) {
				p.On("GetStore").Return(store)
				p.On("GetAPI").Return(mockAPI)
			},
			setupAPI: func() {
				mockAPI.On("DeletePost", "").Return(&model.AppError{
					Message: "Error while deleting a post",
				})
				mockAPI.On("LogError", "Unable to to delete post", "msgID", "", "error", &model.AppError{
					Message: "Error while deleting a post",
				})
			},
			setupStore: func() {
				store.On("GetPostInfoByMSTeamsID", testutils.GetChannelID(), testutils.GetMessageID()).Return(&storemodels.PostInfo{}, nil)
			},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			p := mocksPlugin.NewPluginIface(t)
			testCase.setupPlugin(p)
			testCase.setupAPI()
			testCase.setupStore()

			ah.plugin = p
			ah.handleDeletedActivity(testCase.activityIds)
		})
	}
}
