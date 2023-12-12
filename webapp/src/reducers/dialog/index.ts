import {PayloadAction, createSlice} from '@reduxjs/toolkit';

import {AppDialogs, DialogState} from 'types/common/store.d';

const initialState: Record<string, DialogState> = {};

export const dialogSlice = createSlice({
    name: 'dialogSlice',
    initialState,
    reducers: {
        showDialog: (state, {payload}: PayloadAction<AppDialogs>) => {
            state[payload.dialogId] = {
                show: true,
                ...payload.state,
            };
        },
        closeDialog: (state, {payload}: PayloadAction<string>) => {
            state[payload].show = false;
        },

    },
});

export const {showDialog, closeDialog} = dialogSlice.actions;

export default dialogSlice.reducer;
