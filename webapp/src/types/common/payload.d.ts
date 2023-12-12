type PaginationQueryParams = {
    page: number;
    per_page: number;
}

type UnlinkChannelParams = {
    channelId: string;
}

type SearchLinkedChannelParams = PaginationQueryParams & {
    search?: string;
}

type SearchMSChannelsParams = SearchLinkedChannelParams & {
    teamId: string;
}

type LinkChannelsPayload = {
    mattermostTeamID: string,
    mattermostChannelID: string,
    msTeamsTeamID: string,
    msTeamsChannelID: string,
}

