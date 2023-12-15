import {TeamsState} from 'mattermost-redux/types/teams';

import {ApiRequestCompletionState, ConnectedState, ModalState, ReduxState, RefetchState, SnackbarState} from 'types/common/store.d';

const getPluginState = (state: ReduxState) => state['plugins-com.mattermost.msteams-sync'];

export const getApiRequestCompletionState = (state: ReduxState): ApiRequestCompletionState => getPluginState(state).apiRequestCompletionSlice;

export const getConnectedState = (state: ReduxState): ConnectedState => getPluginState(state).connectedStateSlice;

export const getSnackbarState = (state: ReduxState): SnackbarState => getPluginState(state).snackbarSlice;

export const getIsRhsLoading = (state: ReduxState): {isRhsLoading: boolean} => getPluginState(state).rhsLoadingSlice;

export const getCurrentTeam = (state: ReduxState): TeamsState => state.entities.teams;

export const getLinkModalState = (state: ReduxState): ModalState => getPluginState(state).linkModalSlice;

export const getRefetchState = (state: ReduxState): RefetchState => getPluginState(state).refetchSlice;
