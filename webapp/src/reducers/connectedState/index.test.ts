import {Action} from 'redux';

import reducer, {setConnected} from 'reducers/connectedState';

import {ConnectedState} from 'types/common/store.d';

const initialState: ConnectedState = {
    connected: false,
    isAlreadyConnected: false,
    username: '',
    msteamsUserId: '',
};

describe('Connected State reducer', () => {
    it('should return the initial state', () => {
        expect(reducer(initialState, {} as Action)).toEqual(initialState);
    });

    it('should handle `setConnected`', () => {
        const expectedState: ConnectedState = {...initialState, connected: true, username: 'john doe', msteamsUserId: '1234'};

        expect(reducer(initialState, setConnected({connected: true, username: 'john doe', msteamsUserId: '1234', isAlreadyConnected: false}))).toEqual(expectedState);
    });
});
