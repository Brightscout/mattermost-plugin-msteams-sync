import {Action} from 'redux';

import reducer, {setIsRhsLoading} from 'reducers/spinner';

const initialState:{isRhsLoading: boolean} = {
    isRhsLoading: false,
};

describe('Spinner state reducer', () => {
    it('should return the initial state', () => {
        expect(reducer(initialState, {} as Action)).toEqual(initialState);
    });

    it('should handle `setIsRhsLoading`', () => {
        const expectedState: {isRhsLoading: boolean} = {isRhsLoading: true};

        expect(reducer(initialState, setIsRhsLoading(true))).toEqual(expectedState);
    });
});