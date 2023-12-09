import React from 'react';

import {render, RenderResult} from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {WarningCard} from './WarningCard.component';
import {WarningCardProps} from './WarningCard.types';

const onConnect = jest.fn();

const warningCardProps: WarningCardProps = {
    onConnect,
};

let tree: RenderResult;

describe('Warning Channel', () => {
    beforeEach(() => {
        tree = render(<WarningCard {...warningCardProps}/>);
    });

    it('Should render correctly', () => {
        expect(tree).toMatchSnapshot();
    });

    it('Should call connect function on click of button', () => {
        expect(tree.getAllByRole('button').length).toEqual(1);

        userEvent.click(tree.getAllByRole('button')[0]);
        expect(onConnect).toHaveBeenCalledTimes(1);
    });
});
