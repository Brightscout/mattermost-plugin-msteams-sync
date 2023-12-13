import React, {useEffect} from 'react';
import {Action, Store} from 'redux';
import {useDispatch} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';

import {RhsTitle} from 'components';

import {pluginApiServiceConfigs} from 'constants/apiService.constant';
import {defaultPage, defaultPerPage, pluginTitle} from 'constants/common.constants';
import {iconUrl} from 'constants/illustrations.constants';

import {Rhs} from 'containers';

import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';
import usePluginApi from 'hooks/usePluginApi';

import {setConnected} from 'reducers/connectedState';
import {setIsRhsLoading} from 'reducers/spinner';

// global styles
import 'styles/main.scss';

/**
 * This is the main App component for the plugin
 * @returns {JSX.Element}
 */
const App = ({registry, store}:{registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>}): JSX.Element => {
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

    const {data: whitelistUserData} = getApiState(pluginApiServiceConfigs.whitelistUser.apiServiceName);

    useApiRequestCompletionState({
        serviceName: pluginApiServiceConfigs.needsConnect.apiServiceName,
        handleSuccess: () => {
            const data = needsConnectData as NeedsConnectData;
            dispatch(setConnected({connected: data.connected, username: data.username, msteamsUserId: data.msteamsUserId, isAlreadyConnected: data.connected}));
        },
    });

    useApiRequestCompletionState({
        serviceName: pluginApiServiceConfigs.whitelistUser.apiServiceName,
        handleSuccess: () => {
            const {presentInWhitelist} = whitelistUserData as WhitelistUserResponse;

            // Register the channel header button and app bar if the user is a whitelist user
            if (presentInWhitelist) {
                const {_, toggleRHSPlugin} = registry.registerRightHandSidebarComponent(Rhs, <RhsTitle/>);
                registry.registerChannelHeaderButtonAction(
                    <img
                        width={24}
                        height={24}
                        src={iconUrl}
                        style={{filter: 'grayscale(1)'}}
                    />, () => store.dispatch(toggleRHSPlugin), null, pluginTitle);
                if (registry.registerAppBarComponent) {
                    registry.registerAppBarComponent(iconUrl, () => store.dispatch(toggleRHSPlugin), pluginTitle);
                }
            }
        },
    });

    return <></>;
};

export default App;
