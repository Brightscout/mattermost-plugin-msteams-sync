package main

import (
	"bytes"
	"net/http"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-plugin-msteams-sync/server/links"
	"github.com/mattermost/mattermost-plugin-msteams-sync/server/msteams"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

const (
	botUsername    = "msteams"
	botDisplayName = "MS Teams"
	pluginID       = "com.mattermost.msteams-sync-plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	msteamsAppClientMutex sync.Mutex
	msteamsAppClient      msteams.Client
	msteamsBotClientMutex sync.Mutex
	msteamsBotClient      msteams.Client

	botID  string
	userID string

	links *links.LinksService
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	api := NewAPI(p)
	api.ServeHTTP(w, r)
}

func (p *Plugin) getURL() string {
	config := p.API.GetConfig()
	if strings.HasSuffix(*config.ServiceSettings.SiteURL, "/") {
		return *config.ServiceSettings.SiteURL + "plugins/" + pluginID
	}
	return *config.ServiceSettings.SiteURL + "/plugins/" + pluginID
}

func (p *Plugin) connectTeamsAppClient() error {
	p.msteamsAppClientMutex.Lock()
	defer p.msteamsAppClientMutex.Unlock()

	if p.msteamsAppClient == nil {
		p.msteamsAppClient = msteams.NewApp(
			p.configuration.TenantId,
			p.configuration.ClientId,
			p.configuration.ClientSecret,
		)
	}
	err := p.msteamsAppClient.Connect()
	if err != nil {
		p.API.LogError("Unable to connect to the app client", "error", err)
		return err
	}
	return nil
}

func (p *Plugin) connectTeamsBotClient() error {
	p.msteamsBotClientMutex.Lock()
	defer p.msteamsBotClientMutex.Unlock()
	if p.msteamsBotClient == nil {
		p.msteamsBotClient = msteams.NewBot(
			p.configuration.TenantId,
			p.configuration.ClientId,
			p.configuration.ClientSecret,
			p.configuration.BotUsername,
			p.configuration.BotPassword,
		)
	}
	err := p.msteamsBotClient.Connect()
	if err != nil {
		p.API.LogError("Unable to connect to the bot client", "error", err)
		return err
	}
	return nil
}

func (p *Plugin) start() {
	err := p.connectTeamsAppClient()
	if err != nil {
		p.API.LogError("Unable to connect to the msteams", "error", err)
		return
	}
	err = p.connectTeamsBotClient()
	if err != nil {
		p.API.LogError("Unable to connect to the msteams", "error", err)
		return
	}
	p.links = links.New(p.API, p.msteamsAppClient)
	if err := p.links.Load(); err != nil {
		p.API.LogError("Unable to start the links service", "error", err)
		p.links = nil
		return
	}
	p.links.UpdateEnabledTeams(p.configuration.EnabledTeams)
	p.links.UpdateWebhookSecret(p.configuration.WebhookSecret)
	p.links.UpdateNotificationURL(p.getURL() + "/")
	if err = p.links.Start(); err != nil {
		p.API.LogError("Unable to start the links service", "error", err)
		return
	}

}

func (p *Plugin) stop() {
	if p.links != nil {
		p.links.Stop()
	}
}

func (p *Plugin) restart() {
	p.stop()
	p.start()
}

func (p *Plugin) OnActivate() error {
	bot, appErr := p.API.CreateBot(&model.Bot{
		Username:    botUsername,
		DisplayName: botDisplayName,
		Description: "Created by the MS Teams Sync plugin.",
	})
	if appErr != nil {
		bot, appErr := p.API.GetUserByUsername(botUsername)
		if appErr != nil {
			return appErr
		}
		p.userID = bot.Id
	} else {
		p.userID = bot.UserId
	}

	err := p.API.RegisterCommand(createMsteamsSyncCommand())
	if err != nil {
		return err
	}

	go p.start()
	return nil
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	if post.Props != nil {
		if _, ok := post.Props["msteams_sync_"+p.userID].(bool); ok {
			return
		}
	}

	link := p.links.GetLinkByChannelID(post.ChannelId)
	if link == nil {
		return
	}

	user, _ := p.API.GetUser(post.UserId)

	p.Send(link, user, post)
}

func (p *Plugin) MessageHasBeenUpdated(c *plugin.Context, newPost, oldPost *model.Post) {
	p.API.LogError("UPDATING POST", "oldPost", oldPost.Id, "newPost", newPost.Id)
	if oldPost.Props != nil {
		if _, ok := oldPost.Props["msteams_sync_"+p.userID].(bool); ok {
			return
		}
	}

	link := p.links.GetLinkByChannelID(newPost.ChannelId)
	if link == nil {
		return
	}

	user, _ := p.API.GetUser(newPost.UserId)

	p.Update(link, user, newPost, oldPost)
}

func (p *Plugin) OnDeactivate() error {
	p.stop()
	return nil
}
func (p *Plugin) checkEnabledTeamByTeamId(teamId string) bool {
	if p.configuration.EnabledTeams == "" {
		return true
	}
	team, appErr := p.API.GetTeam(teamId)
	if appErr != nil {
		return false
	}
	isTeamEnabled := false
	enabledTeams := strings.Split(p.configuration.EnabledTeams, ",")
	for _, enabledTeam := range enabledTeams {
		if team.Name == enabledTeam {
			isTeamEnabled = true
			break
		}
	}
	return isTeamEnabled
}

func (p *Plugin) Send(link *links.ChannelLink, user *model.User, post *model.Post) (string, error) {
	p.API.LogDebug("Sending message to MS Teams", "link", link, "post", post)

	// TODO: Replace this with a template
	text := user.Username + "@mattermost: " + post.Message

	parentID := []byte{}
	if post.RootId != "" {
		parentID, _ = p.API.KVGet(mattermostTeamsPostKey(post.RootId))
	}

	var attachments []*msteams.Attachment
	for _, fileId := range post.FileIds {
		fileInfo, appErr := p.API.GetFileInfo(fileId)
		if appErr != nil {
			p.API.LogError("Unable to get file attachment", "error", appErr)
			continue
		}
		p.API.LogError("adding attachment", "attachment", msteams.Attachment{
			ID:           fileInfo.Id,
			Name:         fileInfo.Name,
			ThumbnailURL: fileInfo.ThumbnailPath,
			ContentURL:   fileInfo.Path,
			ContentType:  fileInfo.MimeType,
		})

		fileData, appErr := p.API.GetFile(fileInfo.Id)
		if appErr != nil {
			p.API.LogError("error get file attachment from mattermost", "error", appErr)
			continue
		}

		attachment, err := p.msteamsBotClient.UploadFile(link.MSTeamsTeam, link.MSTeamsChannel, fileInfo.Id+"_"+fileInfo.Name, int(fileInfo.Size), fileInfo.MimeType, bytes.NewReader(fileData))
		if err != nil {
			p.API.LogError("error uploading attachment", "error", err)
			continue
		}
		attachments = append(attachments, attachment)
	}

	newMessageId, err := p.msteamsBotClient.SendMessageWithAttachments(link.MSTeamsTeam, link.MSTeamsChannel, string(parentID), text, attachments)
	if err != nil {
		p.API.LogError("Error creating post", "error", err)
		return "", err
	}

	if post.Id != "" && newMessageId != "" {
		p.API.KVSet(mattermostTeamsPostKey(post.Id), []byte(newMessageId))
		p.API.KVSet(teamsMattermostPostKey(newMessageId), []byte(post.Id))
	}
	return newMessageId, nil
}

func (p *Plugin) Delete(link links.ChannelLink, user *model.User, post *model.Post) error {
	p.API.LogDebug("Sending message to MS Teams", "link", link, "post", post)

	parentID := []byte{}
	if post.RootId != "" {
		parentID, _ = p.API.KVGet(mattermostTeamsPostKey(post.RootId))
	}

	msgID, _ := p.API.KVGet(mattermostTeamsPostKey(post.Id))

	err := p.msteamsBotClient.DeleteMessage(link.MSTeamsTeam, link.MSTeamsChannel, string(parentID), string(msgID))
	if err != nil {
		p.API.LogError("Error deleting post", "error", err)
		return err
	}
	return nil
}

func (p *Plugin) Update(link *links.ChannelLink, user *model.User, newPost, oldPost *model.Post) error {
	p.API.LogDebug("Sending message to MS Teams", "link", link, "oldPost", oldPost, "newPost", newPost)

	parentID := []byte{}
	if oldPost.RootId != "" {
		parentID, _ = p.API.KVGet(mattermostTeamsPostKey(newPost.RootId))
	}

	msgID, _ := p.API.KVGet(mattermostTeamsPostKey(newPost.Id))

	// TODO: Replace this with a template
	text := user.Username + "@mattermost: " + newPost.Message

	p.API.LogDebug("update post data", "msgID", msgID, "parentID", parentID, "text", text)

	err := p.msteamsBotClient.UpdateMessage(link.MSTeamsTeam, link.MSTeamsChannel, string(parentID), string(msgID), text)
	if err != nil {
		p.API.LogError("Error updating the post", "error", err)
		return err
	}

	return nil
}
