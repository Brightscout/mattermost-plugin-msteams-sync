import {combineReducers} from 'redux';

import {msTeamsPluginApi} from 'services';

import apiRequestCompletionSlice from 'reducers/apiRequest';
import connectedStateSlice from 'reducers/connectedState';
import snackbarSlice from 'reducers/snackbar';
import dialogSlice from 'reducers/dialog';
import rhsLoadingSlice from 'reducers/spinner';
import linkModalSlice from 'reducers/linkModal';

export default combineReducers({
    apiRequestCompletionSlice,
    connectedStateSlice,
    snackbarSlice,
    dialogSlice,
    rhsLoadingSlice,
    linkModalSlice,
    [msTeamsPluginApi.reducerPath]: msTeamsPluginApi.reducer,
});
