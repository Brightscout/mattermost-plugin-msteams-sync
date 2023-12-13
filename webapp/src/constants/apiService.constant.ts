// Plugin api service (RTK query) configs
export const pluginApiServiceConfigs: Record<PluginApiServiceName, PluginApiService> = {
    connect: {
        path: '/connect',
        method: 'GET',
        apiServiceName: 'connect',
    },
    needsConnect: {
        path: '/needsConnect',
        method: 'GET',
        apiServiceName: 'needsConnect',
    },
    whitelistUser: {
        path: '/whitelist-user',
        method: 'GET',
        apiServiceName: 'whitelistUser',
    },
    getLinkedChannels: {
        path: '/linked-channels',
        method: 'GET',
        apiServiceName: 'getLinkedChannels',
    },
    disconnectUser: {
        path: '/disconnect',
        method: 'GET',
        apiServiceName: 'disconnectUser',
    },
    searchMSTeams: {
        path: '/msteams/teams',
        method: 'GET',
        apiServiceName: 'searchMSTeams',
    },
    searchMSChannels: {
        path: '/msteams/teams/{team_id}/channels',
        method: 'GET',
        apiServiceName: 'searchMSChannels',
    },
    linkChannels: {
        path: '/channels/link',
        method: 'POST',
        apiServiceName: 'linkChannels',
    },
    unlinkChannel: {
        path: '/channels/{channel_id}/unlink',
        method: 'DELETE',
        apiServiceName: 'unlinkChannel',
    },
};
