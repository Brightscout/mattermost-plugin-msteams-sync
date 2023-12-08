import React, {useEffect, useRef} from 'react';
import {useDispatch} from 'react-redux';

import {Button} from '@brightscout/mattermost-ui-library';

import {Icon} from 'components/Icon';

import {getSnackbarState} from 'selectors';

import usePluginApi from 'hooks/usePluginApi';
import {closeAlert} from 'reducers/snackbar';
import {alertTimeout} from 'constants/common.constants';

import {SnackbarColor} from './Snackbar.types';

export const Snackbar = () => {
    const dispatch = useDispatch();
    const {state} = usePluginApi();
    const timeId = useRef(0);
    const {isOpen, message, severity, icon} = getSnackbarState(state);

    const handleClose = () => dispatch(closeAlert());

    useEffect(() => {
        if (isOpen) {
            timeId.current = window.setTimeout(() => {
                // Hide the snackbar after 4 seconds
                handleClose();
            }, alertTimeout);
        }

        return () => {
            clearTimeout(timeId.current);
        };
    }, [isOpen]);

    const snackbarColorMap: Record<SnackbarColor, string> = {
        error: 'bg-error',
        default: 'bg-default',
        success: 'bg-success',
    };

    return (
        <div
            className={`fixed bottom-20 right-20 left-20 py-8 px-12 rounded-4 d-flex gap-8 align-items-center justify-between elevation-2 msteams-sync-rhs__snackbar ${snackbarColorMap[severity]}`}
        >
            <div className='d-flex align-items-center gap-8'>
                <Icon
                    iconName={icon}
                    className='icon-white icon-16'
                />
                <h5 className='my-0 lh-24 wt-600 text-white'>{message}</h5>
            </div>
            <Button
                variant='text'
                className='snackbar__close'
                onClick={handleClose}
            >
                <Icon
                    iconName='close'
                    className='icon-white icon-16'
                />
            </Button>
        </div>
    );
};
