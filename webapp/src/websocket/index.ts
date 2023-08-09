import {Action, Store} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import {setGlobalModalState} from 'src/reducers/globalModal';

import {setConnected} from 'src/reducers/connectedState';
import {ModalIds} from 'src/constants';

export function handleConnect(store: Store<GlobalState, Action<Record<string, unknown>>>) {
    return (_: WebsocketEventParams) => {
        store.dispatch(setConnected(true) as Action);
    };
}

export function handleDisconnect(store: Store<GlobalState, Action<Record<string, unknown>>>) {
    return (_: WebsocketEventParams) => {
        store.dispatch(setConnected(false) as Action);
    };
}

export function handleOpenLinkChannelsModal(store: Store<GlobalState, Action<Record<string, unknown>>>) {
    return (_: WebsocketEventParams) => {
        store.dispatch(setGlobalModalState({modalId: ModalIds.LINK_CHANNELS}) as Action);
    };
}
