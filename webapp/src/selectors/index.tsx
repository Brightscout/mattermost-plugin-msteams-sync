import {ApiRequestCompletionState, ConnectedState, DialogState, ModalState, ReduxState, SnackbarState} from 'types/common/store.d';

const getPluginState = (state: ReduxState) => state['plugins-com.mattermost.msteams-sync'];

export const getApiRequestCompletionState = (state: ReduxState): ApiRequestCompletionState => getPluginState(state).apiRequestCompletionSlice;

export const getConnectedState = (state: ReduxState): ConnectedState => getPluginState(state).connectedStateSlice;

export const getSnackbarState = (state: ReduxState): SnackbarState => getPluginState(state).snackbarSlice;

export const getDialogState = (state: ReduxState): DialogState => getPluginState(state).dialogSlice;

export const getIsRhsLoading = (state: ReduxState): {isRhsLoading: boolean} => getPluginState(state).rhsLoadingSlice;

export const getCurrentTeam = (state: ReduxState): string => state.entities.teams.currentTeamId;

export const getLinkModalState = (state: ReduxState): ModalState => getPluginState(state).linkModalSlice;
