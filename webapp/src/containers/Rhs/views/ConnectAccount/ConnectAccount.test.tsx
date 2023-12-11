import React from 'react';

import userEvent from '@testing-library/user-event';
import {RenderResult, render} from '@testing-library/react';

import {mockDispatch} from 'tests/setup';

import {ConnectAccount} from './ConnectAccount.container';

let tree: RenderResult;

describe('Connect Account View', () => {
    beforeEach(() => {
        tree = render(<ConnectAccount/>);
    });

    it('should render correctly', () => {
        expect(tree).toMatchSnapshot();
    });

    it('should render connect account button', () => {
        const connectButton = tree.getByText('Connect Account');

        expect(connectButton).toBeVisible();
    });

    it('should dispatch an action when button is clicked', async () => {
        const connectButton = tree.getByText('Connect Account');

        await userEvent.click(connectButton);

        expect(mockDispatch).toHaveBeenCalledTimes(1);
    });
});
