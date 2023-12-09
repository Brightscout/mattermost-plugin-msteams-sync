import {ReduxState} from 'types/common/store.d';

export const mockTestState: Pick<ReduxState, 'plugins-com.mattermost.msteams-sync'> = {
    'plugins-com.mattermost.msteams-sync': {

        apiRequestCompletionSlice: {
            requests: [
                'whitelistUser',
                'getLinkedChannels',
            ],
        },
        connectedStateSlice: {
            connected: true,
            isAlreadyConnected: true,
            username: 'Saurabh Sharma',
            msteamsUserId: '06e7db55-18e5-4f62-ac5f-4a21999c83e4',
        },
        snackbarSlice: {
            message: '',
            severity: 'default',
            isOpen: false,
        },
        dialogSlice: {
            description: '',
            destructive: false,
            show: false,
            primaryButtonText: '',
            isLoading: false,
            title: '',
        },
        rhsLoadingSlice: {
            isRhsLoading: false,
        },
        msTeamsPluginApi: {
            queries: {
                'whitelistUser(undefined)': {
                    status: 'fulfilled',
                    endpointName: 'whitelistUser',
                    requestId: '4S522qvjH78ABYk_r8nGc',
                    startedTimeStamp: 1702107980963,
                    data: {
                        presentInWhitelist: true,
                    },
                    fulfilledTimeStamp: 1702107981697,
                },
                'needsConnect(undefined)': {
                    status: 'fulfilled',
                    endpointName: 'needsConnect',
                    requestId: 'cS3T_IF0VW4kEJclgd-IV',
                    startedTimeStamp: 1702107980970,
                    data: {
                        canSkip: false,
                        connected: true,
                        msteamsUserId: '06e7db55-18e5-4f62-ac5f-4a21999c83e4',
                        needsConnect: false,
                        username: 'Saurabh Sharma',
                    },
                    fulfilledTimeStamp: 1702107985184,
                },
                'getLinkedChannels({"page":0,"per_page":20})': {
                    status: 'fulfilled',
                    endpointName: 'getLinkedChannels',
                    requestId: 'GXmcNCITzra1TGNGFtjD3',
                    originalArgs: {
                        page: 0,
                        per_page: 20,
                    },
                    startedTimeStamp: 1702107980968,
                    data: [],
                    fulfilledTimeStamp: 1702107981716,
                },
            },
            mutations: {},
            provided: {},
            subscriptions: {
                'whitelistUser(undefined)': {
                    '4S522qvjH78ABYk_r8nGc': {},
                },
                'needsConnect(undefined)': {
                    ghQe6Cf4HXIdPga1D47SD: {},
                    'cS3T_IF0VW4kEJclgd-IV': {},
                },
                'getLinkedChannels({"page":0,"per_page":20})': {
                    GXmcNCITzra1TGNGFtjD3: {},
                },
            },
            config: {
                online: true,
                focused: true,
                middlewareRegistered: false,
                refetchOnFocus: false,
                refetchOnReconnect: false,
                refetchOnMountOrArgChange: false,
                keepUnusedDataFor: 60,
                reducerPath: 'msTeamsPluginApi',
            },
        },
    },
};
