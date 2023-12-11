import React from 'react';

import {RenderResult, render} from '@testing-library/react';

import {Rhs} from './Rhs.container';

let tree: RenderResult;

describe('RHS view', () => {
    beforeEach(() => {
        tree = render(<Rhs/>);
    });

    it('should render correctly', () => {
        expect(tree).toMatchSnapshot();
    });

    it('should render disconnect account button', () => {
        const disconnectButton = tree.getByText('Disconnect');
        expect(disconnectButton).toBeVisible();
    });
});
