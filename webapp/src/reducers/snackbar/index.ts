import {PayloadAction, createSlice} from '@reduxjs/toolkit';

import {SnackbarActionPayload, SnackbarState} from 'types/common/store.d';

const initialState: SnackbarState = {
    message: '',
    severity: 'default',
    isOpen: false,
    icon: 'tick',
};

export const snackbarSlice = createSlice({
    name: 'snackbarState',
    initialState,
    reducers: {
        showAlert: (state, {payload}: PayloadAction<SnackbarActionPayload>) => {
            state.message = payload.message;
            state.severity = payload.severity;
            state.isOpen = true;
            state.icon = payload.icon;
        },
        closeAlert: (state) => {
            state.isOpen = false;
        },
    },
});

export const {showAlert, closeAlert} = snackbarSlice.actions;

export default snackbarSlice.reducer;
