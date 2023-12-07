import React, {useEffect} from 'react';
import {useDispatch} from 'react-redux';

import usePluginApi from 'hooks/usePluginApi';

//global styles
import {pluginApiServiceConfigs} from 'constants/apiService.constant';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

import {setConnected} from 'reducers/connectedState';
import {defaultPage, defaultPerPage} from 'constants/common.constants';
import {setIsRhsLoading} from 'reducers/spinner';

import 'styles/main.scss';

/**
 * This is the main App component for the plugin
 * @returns {JSX.Element}
 */
const App = (): JSX.Element => {
    const dispatch = useDispatch();
    const {makeApiRequestWithCompletionStatus, getApiState} = usePluginApi();

    useEffect(() => {
        const linkedChannelsParams: SearchLinkedChannelParams = {page: defaultPage, per_page: defaultPerPage};

        makeApiRequestWithCompletionStatus(pluginApiServiceConfigs.whitelistUser.apiServiceName);
        makeApiRequestWithCompletionStatus(pluginApiServiceConfigs.needsConnect.apiServiceName);
        makeApiRequestWithCompletionStatus(pluginApiServiceConfigs.getLinkedChannels.apiServiceName, linkedChannelsParams);
    }, []);

    const {data: needsConnectData, isLoading} = getApiState(pluginApiServiceConfigs.needsConnect.apiServiceName);

    useEffect(() => {
        dispatch(setIsRhsLoading(isLoading));
    }, [isLoading]);

    useApiRequestCompletionState({
        serviceName: pluginApiServiceConfigs.needsConnect.apiServiceName,
        handleSuccess: () => {
            const data = needsConnectData as NeedsConnectData;
            dispatch(setConnected({connected: data.connected, username: data.username, msteamsUserId: data.msteamsUserId, isAlreadyConnected: data.connected}));
        },
    });

    return <></>;
};

export default App;
