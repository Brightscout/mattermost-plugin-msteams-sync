import {Action} from 'redux';

import reducer, {showDialog, closeDialog} from 'reducers/dialog';

import {DialogState} from 'types/common/store.d';

const initialState:DialogState = {
    description: '',
    destructive: false,
    isLoading: false,
    primaryButtonText: '',
    show: false,
    title: '',
};

describe('Dialog State reducer', () => {
    it('should return the initial state', () => {
        expect(reducer({}, {} as Action)).toEqual({});
    });

    it('should handle `showDialog`', () => {
        const expectedState: DialogState = {show: true, description: 'description', primaryButtonText: 'Done', title: 'Title'};

        expect(reducer(initialState, showDialog({show: true, description: 'description', primaryButtonText: 'Done', title: 'Title'}))).toEqual(expectedState);
    });

    it('should handle `closeDialog`', () => {
        const expectedState: DialogState = {...initialState, show: false};

        expect(reducer({...initialState, show: true}, closeDialog())).toEqual(expectedState);
    });
});
