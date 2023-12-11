import {PayloadAction, createSlice} from '@reduxjs/toolkit';

import {DialogState, ModalState} from 'types/common/store.d';

const initialState: ModalState = {
    show: false,
    isLoading: false,
};

export const linkModalSlice = createSlice({
    name: 'linkModalSlice',
    initialState,
    reducers: {
        showLinkModal: (state) => {
            state.show = true;
        },
        hideLinkModal: (state) => {
            state.show = false;
        },
        setLinkModalLoading: (state, {payload}: PayloadAction<boolean>) => {
            state.isLoading = payload;
        },
    },
});

export const {showLinkModal, hideLinkModal, setLinkModalLoading} = linkModalSlice.actions;

export default linkModalSlice.reducer;
