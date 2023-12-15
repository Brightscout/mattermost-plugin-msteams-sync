import React from 'react';
import {Store, Action} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import reducer from 'reducers';

import EnforceConnectedAccountModal from 'components/enforceConnectedAccountModal';
import MSTeamsAppManifestSetting from 'components/appManifestSetting';
import ListConnectedUsers from 'components/getConnectedUsersSetting';

import {LinkChannelModal, RhsTitle} from 'components';

import {Rhs} from 'containers';

import {pluginTitle} from 'constants/common.constants';

import {iconUrl} from 'constants/illustrations.constants';

import {handleConnect, handleDisconnect, handleLink, handleModalLink, handleUnlinkChannels} from 'websocket';

import manifest from './manifest';

// eslint-disable-next-line import/no-unresolved
import {PluginRegistry} from './types/mattermost-webapp';
import App from './App';

export default class Plugin {
    enforceConnectedAccountId = '';
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        registry.registerReducer(reducer);
        registry.registerRootComponent(App);
        registry.registerRootComponent(LinkChannelModal);

        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        this.enforceConnectedAccountId = registry.registerRootComponent(EnforceConnectedAccountModal);

        registry.registerAdminConsoleCustomSetting('appManifestDownload', MSTeamsAppManifestSetting);
        registry.registerAdminConsoleCustomSetting('ConnectedUsersReportDownload', ListConnectedUsers);
        const {_, toggleRHSPlugin} = registry.registerRightHandSidebarComponent(Rhs, <RhsTitle/>);

        // TODO: update icons later
        registry.registerChannelHeaderButtonAction(
            <img
                width={24}
                height={24}
                src={iconUrl}
            />, () => store.dispatch(toggleRHSPlugin), null, pluginTitle);

        if (registry.registerAppBarComponent) {
            registry.registerAppBarComponent(iconUrl, () => store.dispatch(toggleRHSPlugin), pluginTitle);
        }

        registry.registerWebSocketEventHandler(`custom_${manifest.id}_connect`, handleConnect(store));
        registry.registerWebSocketEventHandler(`custom_${manifest.id}_disconnect`, handleDisconnect(store));
        registry.registerWebSocketEventHandler(`custom_${manifest.id}_unlink`, handleUnlinkChannels(store));
        registry.registerWebSocketEventHandler(`custom_${manifest.id}_link_channels`, handleModalLink(store));
        registry.registerWebSocketEventHandler(`custom_${manifest.id}_link`, handleLink(store));
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void;
        Components: any;
    }
}

window.registerPlugin(manifest.id, new Plugin());
