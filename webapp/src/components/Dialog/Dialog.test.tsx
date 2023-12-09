import React from 'react';

import {render, RenderResult} from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {DialogProps, LinearProgress} from '@brightscout/mattermost-ui-library';

import {Dialog} from './Dialog.component';

const onCloseHandler = jest.fn();
const onSubmitHandler = jest.fn();

const dialogProps = {
    onCloseHandler,
    onSubmitHandler,
};

let tree: RenderResult;

describe('Dialog', () => {
    beforeEach(() => {
        tree = render(<Dialog {...dialogProps}/>);
    });

    it('Should render correctly', () => {
        expect(tree).toMatchSnapshot();
    });

    it('Should close the dialog on clicking close button', () => {
        expect(tree.getAllByRole('button').length).toEqual(2);

        userEvent.click(tree.getAllByRole('button')[1]);
        expect(onCloseHandler).toHaveBeenCalledTimes(1);
    });
});
