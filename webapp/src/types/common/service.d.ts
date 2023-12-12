type HttpMethod = 'GET' | 'POST' | 'PATCH' | 'DELETE';

type PluginApiServiceName = 'needsConnect' | 'connect' | 'whitelistUser' | 'getLinkedChannels' | 'disconnectUser' | 'searchMSTeams' | 'searchMSChannels' | 'linkChannels';

type PluginApiService = {
    path: string,
    method: httpMethod,
    apiServiceName: PluginApiServiceName,
}

type APIError = {
    id: string,
    message: string,
}

type APIRequestPayload = LinkChannelsPayload | SearchLinkedChannelParams | void;
