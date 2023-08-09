import {combineReducers} from 'redux';

import {msTeamsPluginApi} from 'src/services';

import apiRequestCompletionSlice from 'src/reducers/apiRequest';
import connectedReducer from 'src/reducers/connectedState';
import globalModalSlice from 'src/reducers/globalModal';

export default combineReducers({
    apiRequestCompletionSlice,
    connectedReducer,
    globalModalSlice,
    [msTeamsPluginApi.reducerPath]: msTeamsPluginApi.reducer,
});
