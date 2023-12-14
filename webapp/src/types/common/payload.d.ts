type PaginationQueryParams = {
    page: number;
    per_page: number;
}

type UnlinkChannelParams = {
    channelId: string;
}

type SearchParams = PaginationQueryParams & {
    search?: string;
}

type SearchMSChannelsParams = SearchParams & {
    teamId: string;
}
