import {ApiRequestCompletionState, ConnectedState, ReduxState, SnackbarState} from 'types/common/store.d';

export const getApiRequestCompletionState = (state: ReduxState['plugins-com.mattermost.msteams-sync']): ApiRequestCompletionState => state.apiRequestCompletionSlice;

export const getConnectedState = (state: ReduxState['plugins-com.mattermost.msteams-sync']): ConnectedState => state.connectedStateSlice;

export const getSnackbarState = (state: ReduxState['plugins-com.mattermost.msteams-sync']): SnackbarState => state.snackbarSlice;

export const getIsRhsLoading = (state: ReduxState['plugins-com.mattermost.msteams-sync']): {isRhsLoading: boolean} => state.rhsLoadingSlice;
